package middleware

import (
	"context"
	"fmt"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/auth"
	"strconv"

	"github.com/XJTU-zxc/GoTikMall/app/auth/biz/jwtutil"
	"github.com/cloudwego/kitex/pkg/endpoint"
)

// AuthMiddleware 身份验证中间件
func AuthMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) (err error) {
		// 获取请求的方法名
		method, ok := ctx.Value("methodName").(string)
		if !ok {
			return next(ctx, req, resp)
		}
		// 检查是否在白名单中
		if jwtutil.IsWhitelisted(method) {
			return next(ctx, req, resp)
		}
		// 进行身份验证
		verifyReq, ok := req.(*auth.VerifyTokenReq)
		if !ok {
			return next(ctx, req, resp)
		}
		userID, valid, err := jwtutil.VerifyToken(verifyReq.Token)
		if err != nil || !valid {
			return err
		}
		// 进行权限校验
		ok, err = jwtutil.Enforcer.Enforce(
			strconv.Itoa(int(userID)),
			method,
			"access",
		)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("权限不足")
		}
		return next(ctx, req, resp)
	}
}
