package service

import (
	"context"

	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/model"
	cart "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type GetCartService struct {
	ctx context.Context
} // NewGetCartService new GetCartService
func NewGetCartService(ctx context.Context) *GetCartService {
	return &GetCartService{ctx: ctx}
}

// Run create note info
func (s *GetCartService) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	// Finish your business logic.
	if req.UserId == 0 {
		return nil, kerrors.NewGRPCBizStatusError(2004003, "user_id is required")
	}
	carts, err := model.GetCartByUserId(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var items []*cart.CartItem
	for _, v := range carts {
		items = append(items, &cart.CartItem{ProductId: v.ProductId, Quantity: v.Quantity})
	}

	return &cart.GetCartResp{Cart: &cart.Cart{UserId: req.GetUserId(), Items: items}}, nil
}
