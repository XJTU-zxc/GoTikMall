package service

import (
	"context"

	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/model"
	"github.com/XJTU-zxc/GoTikMall/app/cart/rpc"
	cart "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/cart"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type AddItemService struct {
	ctx context.Context
} // NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
}

// Run create note info
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// Finish your business logic.
	productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: req.Item.ProductId})
	if err != nil {
		return nil, err
	}
	if productResp.Product == nil || productResp.Product.Id == 0 {
		return nil, kerrors.NewGRPCBizStatusError(2004001, "product not found")
	}
	cartItem := &model.Cart{
		UserId:    req.UserId,
		ProductId: req.Item.ProductId,
		Quantity:  req.Item.Quantity,
	}
	if err := model.AddItem(mysql.DB, s.ctx, cartItem); err != nil {
		return nil, err
	}
	return &cart.AddItemResp{}, nil
}
