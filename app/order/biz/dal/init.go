package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/order/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
