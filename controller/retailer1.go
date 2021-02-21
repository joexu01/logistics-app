package controller

import "github.com/gin-gonic/gin"

type Retailer1Controller struct{}

func Retailer1Register(router *gin.RouterGroup) {
	r := &Retailer1Controller{}
	router.POST("/sign/:id", r.SignForPackage)

	//router.GET("/product/:id", l.ReadProductInfo)
	//router.GET("/tracking/:id", l.ReadTrackingResult)
}

// SignForPackage godoc
// @Summary 签收货物
// @Description 签收货物
// @Tags 零售商1-Retailer1
// @Accept  json
// @Produce  json
// @Param body body dto.UpdateLogisticRecordInput true "body"
// @Param id path string true "物流追踪ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /logistics/update/{id} [post]
func (r *Retailer1Controller) SignForPackage(c *gin.Context) {

}
