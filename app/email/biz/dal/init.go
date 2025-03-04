package dal

import (
	"github.com/XJTU-zxc/GoTikMall/app/email/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/email/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
