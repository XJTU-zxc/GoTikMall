package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/payment/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
