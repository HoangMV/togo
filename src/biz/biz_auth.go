package biz

import (
	"errors"
	"math/rand"
	"time"

	"github.com/HoangMV/todo/lib/log"
	"github.com/HoangMV/todo/src/models/entity"
	"github.com/HoangMV/todo/src/models/request"
	"github.com/HoangMV/todo/src/models/response"
	"github.com/google/uuid"
)

func (biz *Biz) Register(req *request.LoginReq) error {
	// Add New user
	user := &entity.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := biz.dao.InsertUser(user); err != nil {
		log.Get().Errorf("Biz::Register InsertUser err: %v, data: %+v", err, user)
		return errors.New("add user error")
	}

	// Add User max todo require
	rand.Seed(time.Now().UnixNano())
	min := 10
	max := 30
	conf := &entity.UserTodoConfig{
		UserID:  user.ID,
		MaxTodo: rand.Intn(max-min+1) + min,
	}

	if err := biz.dao.InsertUserMaxTodo(conf); err != nil {
		log.Get().Errorf("Biz::Register InsertUserMaxTodo err: %v, data: %+v", err, conf)
	}

	return nil
}

func (biz *Biz) Login(req *request.LoginReq) (*response.LoginResp, error) {
	user, err := biz.dao.GetUserByUsername(req.Username)
	if err != nil {
		log.Get().Errorf("Biz::Login GetUserByUsername err: %v", err)
		return nil, errors.New("username incorrect")
	}

	if req.Password != user.Password {
		return nil, errors.New("password incorrect")
	}

	// create token
	token := uuid.NewString()
	biz.dao.SetTokenToCache(token, user.ID)

	resp := &response.LoginResp{
		UserID: user.ID,
		Token:  token,
	}

	return resp, nil
}

func (biz *Biz) CheckAuth(token string) int {
	return biz.dao.GetTokenInCache(token)
}
