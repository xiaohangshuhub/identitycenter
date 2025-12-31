package user

import "errors"

var (
	ErrInvalidPassword       = errors.New("invalid password")             // 密码错误
	ErrUserNotFound          = errors.New("user not found")               // 用户不存在
	ErrUserExist             = errors.New("user already exists")          // 用户已存在
	ErrInvalidUserOrPassword = errors.New("invalid username or password") // 用户名或密码错误
	ErrSessionlogIdIsnil     = errors.New("session log id is nil")        // session log id 为空
	ErrSessionlogUserIdIsnil = errors.New("session log user id is nil")   // session log user id 为空
	ErrAuthGrantTypeIsnil    = errors.New("auth grant type is nil")       // auth grant type is nil

)
