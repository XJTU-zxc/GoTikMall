package service

import (
	"context"
	"strconv"

	"github.com/XJTU-zxc/GoTikMall/app/checkout/infra/mq"
	"github.com/XJTU-zxc/GoTikMall/app/checkout/infra/rpc"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/cart"
	checkout "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/checkout"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/email"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/order"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/payment"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

// Run create note info
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// Finish your business logic.
	/*
		Refer to the business process of the settlement service in the official textbook.
		1. get cart
		2. calculate cart
		3. create order
		4. empty cart
		5. pay
		6. change order result
		7. finish
	*/
	// get cart
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error(err)
		return nil, kerrors.NewGRPCBizStatusError(2004005, err.Error())
	}
	if cartResult == nil || cartResult.Cart == nil || len(cartResult.Cart.Items) == 0 {
		return nil, kerrors.NewGRPCBizStatusError(2005005, "cart is empty")
	}
	var (
		oi    []*order.OrderItem
		total float32
	)
	for _, cartItem := range cartResult.Cart.Items {
		productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: cartItem.ProductId})

		if resultErr != nil {
			klog.Error(resultErr)
			return nil, resultErr
		}
		if productResp.Product == nil {
			continue
		}

		p := productResp.Product

		cost := p.Price * float32(cartItem.Quantity)
		total += cost

		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{ProductId: cartItem.ProductId, Quantity: cartItem.Quantity},
			Cost: cost,
		})
	}

	// create order
	orderReq := &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: "USD",
		OrderItems:   oi,
		Email:        req.Email,
	}
	if req.Address != nil {
		addr := req.Address
		zipCodeInt, _ := strconv.Atoi(addr.ZipCode)
		orderReq.Address = &order.Address{
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			Country:       addr.Country,
			State:         addr.State,
			ZipCode:       int32(zipCodeInt),
		}
	}
	orderResult, err := rpc.OrderClient.PlaceOrder(s.ctx, orderReq)
	if err != nil {
		klog.Error(err.Error())
		return
	}
	klog.Info(orderResult)

	// empty cart
	emptyResult, err := rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error(err.Error())
		return
	}
	klog.Info(emptyResult)

	// charge
	var orderId string
	if orderResult != nil && orderResult.Order != nil {
		orderId = orderResult.Order.OrderId
	}
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
		},
	}
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		klog.Error(err.Error())
		return nil, err
	}

	// send email message
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You have just created an order in GoTikMall",
		Content:     "You have just created an order in GoTikMall",
	})
	msg := &nats.Msg{Subject: "email", Data: data, Header: make(nats.Header)}
	// open telemetry tracing
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))
	_ = mq.Nc.PublishMsg(msg)
	klog.Info(paymentResult)

	// change order state
	_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{UserId: req.UserId, OrderId: orderId})
	if err != nil {
		klog.Error(err)
		return
	}

	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return
}
