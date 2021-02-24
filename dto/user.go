package dto

import (
	"time"
)

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""` // 返回的Token
}

type AdminInfoOutput struct {
	Id           int       `json:"id"`
	UserName     string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}