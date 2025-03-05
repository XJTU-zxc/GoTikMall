package jwtutil

import (
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	jwt "github.com/golang-jwt/jwt/v5"
)

// 全局 CasBin 执行器
var Enforcer *casbin.Enforcer

// JWT 签名密钥
const JwtSecret = "test_key"

// 自定义 KeyMatch 函数以符合 govaluate.ExpressionFunction 类型
func keyMatchFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, nil
	}
	key1, ok1 := args[0].(string)
	key2, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return false, nil
	}
	return util.KeyMatch(key1, key2), nil
}

// InitAuth 初始化认证服务
func InitAuth() {
	var err error
	// 加载 RBAC 模型和策略文件
	Enforcer, err = casbin.NewEnforcer("rbac_model.conf", "policy.csv")
	if err != nil {
		panic(err)
	}
	// 启用权限匹配函数
	Enforcer.AddFunction("keyMatch", keyMatchFunc)
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID int32) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour) // 72 小时过期
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken 验证 JWT 令牌
// VerifyToken 验证 JWT 令牌
func VerifyToken(tokenStr string) (int32, bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false, fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false, fmt.Errorf("invalid user_id claim")
	}

	return int32(userID), true, nil
}
