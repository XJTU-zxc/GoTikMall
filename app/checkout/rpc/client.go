package rpc

import (
	"sync"

	"github.com/XJTU-zxc/GoTikMall/app/checkout/conf"
	checkoututils "github.com/XJTU-zxc/GoTikMall/app/checkout/utils"
	"github.com/XJTU-zxc/GoTikMall/common/clientsuite"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/client"
)

var (
	CartClient    cartservice.Client
	ProductClient productcatalogservice.Client
	PaymentClient paymentservice.Client
	OrderClient   orderservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
	commonSuite   client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		serviceName = conf.GetConf().Kitex.Service
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			CurrentServiceName: serviceName,
			RegistryAddr:       registryAddr,
		})
		initCartClient()
		initProductClient()
		initPaymentClient()
		initOrderClient()
	})
}

func initProductClient() {
	ProductClient, err = productcatalogservice.NewClient("product", commonSuite)
	checkoututils.MustHandleError(err)
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart", commonSuite)
	checkoututils.MustHandleError(err)
}

func initPaymentClient() {
	PaymentClient, err = paymentservice.NewClient("payment", commonSuite)
	checkoututils.MustHandleError(err)
}

func initOrderClient() {
	OrderClient, err = orderservice.NewClient("order", commonSuite)
	checkoututils.MustHandleError(err)
}
