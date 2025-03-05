package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/auth/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/auth/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
