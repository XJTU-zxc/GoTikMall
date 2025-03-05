package main

import (
	"net"
	"time"

	"github.com/XJTU-zxc/GoTikMall/app/checkout/biz/dal"
	"github.com/XJTU-zxc/GoTikMall/app/checkout/conf"
	"github.com/XJTU-zxc/GoTikMall/app/checkout/infra/mq"
	"github.com/XJTU-zxc/GoTikMall/app/checkout/infra/rpc"
	"github.com/XJTU-zxc/GoTikMall/common/mtl"
	"github.com/XJTU-zxc/GoTikMall/common/serversuite"
	"github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)
var (
	serviceName  = conf.GetConf().Kitex.Service
	RegisterAddr = conf.GetConf().Registry.RegistryAddress[0]
)

func main() {
	err := godotenv.Load()
	if err != nil {
		klog.Error(err.Error())
	}

	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetircsPort, RegisterAddr)

	dal.Init()

	mq.Init()
	rpc.InitClient()

	opts := kitexInit()

	svr := checkoutservice.NewServer(new(CheckoutServiceImpl), opts...)

	err = svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServerSuite{
		CurrentServerName: serviceName,
		RegistryAddr:      RegisterAddr,
	}))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}
