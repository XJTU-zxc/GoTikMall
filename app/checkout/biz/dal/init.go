package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/checkout/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
