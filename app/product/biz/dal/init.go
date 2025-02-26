package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/product/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
