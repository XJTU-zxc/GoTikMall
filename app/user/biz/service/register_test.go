package service

import (
	"context"
	"testing"

	"github.com/XJTU-zxc/GoTikMall/app/user/biz/dal/mysql"
	user "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/user"
	"github.com/joho/godotenv"
)

func TestRegister_Run(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("godotenv.Load err: %v", err)
	}
	mysql.Init()
	ctx := context.Background()
	s := NewRegisterService(ctx)
	// init req and assert value

	req := &user.RegisterReq{
		Email: "test@test.com",
		Password: "test",
		ConfirmPassword: "test",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
