package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dto"
	"github.com/joexu01/logistics-app/middleware"
	"time"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	admin := &AdminLoginController{}
	group.POST("/login", admin.AdminLogin)
	group.GET("/info", admin.AdminInfo)
}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /account/login
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /account/login [post]
func (a *AdminLoginController) AdminLogin(c *gin.Context) {
	out := dto.AdminLoginOutput{Token: "admin"}
	middleware.ResponseSuccess(c, out)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息接口
// @Tags 管理员接口
// @ID /admin/info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /account/info [get]
func (a *AdminLoginController) AdminInfo(c *gin.Context) {
	out := &dto.AdminInfoOutput{
		Id:           1,
		UserName:     "admin",
		LoginTime:    time.Now(),
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}
