package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/user/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
