package service

import (
	"context"
	"errors"

	"github.com/XJTU-zxc/GoTikMall/app/user/biz/dal/mysql"
	"github.com/XJTU-zxc/GoTikMall/app/user/biz/model"
	user "github.com/XJTU-zxc/GoTikMall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Finish your business logic.
	klog.Infof("LoginReq:%+v", req)
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email or password is empty")
	}
	userRow, err := model.GetByEmail(mysql.DB, s.ctx, req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userRow.PasswordHashed), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	return &user.LoginResp{UserId: int32(userRow.ID)}, nil
}
