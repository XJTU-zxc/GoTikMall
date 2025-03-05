package service

import (
	"context"

	"github.com/XJTU-zxc/GoTikMall/app/auth/biz/jwtutil"
	auth "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/auth"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService
func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	// Finish your business logic.
	token, err := jwtutil.GenerateToken(req.UserId)
	if err != nil {
		return nil, err
	}
	return &auth.DeliveryResp{Token: token}, nil
}
