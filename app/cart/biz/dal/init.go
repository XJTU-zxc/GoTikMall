package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
