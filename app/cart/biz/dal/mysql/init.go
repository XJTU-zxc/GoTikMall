package mysql

import (
	"fmt"
	"os"

	"github.com/XJTU-zxc/GoTikMall/app/cart/biz/model"
	"github.com/XJTU-zxc/GoTikMall/app/cart/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	DB, err = gorm.Open(mysql.Open(fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	// if err := DB.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithTracerProvider(mtl.TracerProvider))); err != nil {
	// 	panic(err)
	// }
	fmt.Println("mysql init success")
	if os.Getenv("GO_ENV") != "online" {
		DB.AutoMigrate(
			&model.Cart{},
		)
	}
}
