package service

import (
	"context"
	"strconv"

	"github.com/XJTU-zxc/GoTikMall/app/auth/biz/jwtutil"
	auth "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	// Finish your business logic.
	userID, valid, err := jwtutil.VerifyToken(req.Token)
	if err != nil || !valid {
		return &auth.VerifyResp{Res: false}, nil
	}

	// 使用 Casbin 进行权限校验
	ok, err := jwtutil.Enforcer.Enforce(
		strconv.Itoa(int(userID)),
		"auth_service",
		"access",
	)
	if err != nil {
		return nil, err
	}
	return &auth.VerifyResp{Res: ok}, nil
}
