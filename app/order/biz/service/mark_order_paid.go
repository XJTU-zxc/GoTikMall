package service

import (
	"context"

	"github.com/XJTU-zxc/GoTikMall/app/order/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/order/biz/model"
	order "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
)

type MarkOrderPaidService struct {
	ctx context.Context
} // NewMarkOrderPaidService new MarkOrderPaidService
func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

// Run create note info
func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	// Finish your business logic.
	if req.UserId == 0 || req.OrderId == "" {
		return nil, kerrors.NewGRPCBizStatusError(2004006, "user_id or order_id is empty")
	}
	_, err = model.GetOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf(err.Error())
		return nil, kerrors.NewGRPCBizStatusError(2005006, err.Error())
	}
	err = model.UpdateOrderStatus(mysql.DB, s.ctx, req.UserId, req.OrderId, "prepaid")
	if err != nil {
		klog.Errorf(err.Error())
		return nil, kerrors.NewGRPCBizStatusError(2005016, err.Error())
	}
	return &order.MarkOrderPaidResp{}, nil
}
