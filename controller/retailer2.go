package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dao"
	"github.com/joexu01/logistics-app/dto"
	"github.com/joexu01/logistics-app/fabsdk"
	"github.com/joexu01/logistics-app/middleware"
	"strings"
)

type Retailer2Controller struct{}

func Retailer2Register(router *gin.RouterGroup) {
	r := &Retailer2Controller{}
	router.POST("/sign/:id", r.SignForPackage)

	router.GET("/product/:id", r.ReadProductInfo)
	router.GET("/tracking/:id", r.ReadTrackingResult)
	router.GET("/order/:id", r.ReadOrderInfo)
}

// SignForPackage godoc
// @Summary 签收货物
// @Description 签收货物
// @Tags 零售商2-Retailer2
// @Accept  json
// @Produce  json
// @Param body body dto.UpdateLogisticRecordInput true "body"
// @Param id path string true "物流追踪ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /retailer2/sign/{id} [post]
func (r *Retailer2Controller) SignForPackage(c *gin.Context) {
	trackingID := c.Param("id")
	if trackingID == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty"))
		return
	}

	input := &dto.UpdateLogisticRecordInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRetailer2)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	args := strings.Join([]string{
		`"` + trackingID + `"`,
		`"` + `零售商2-Retailer2: ` + input.Status + `"`,
	}, ",")

	invoke, err := sdkCtx.Invoke(fabsdk.FuncUpdateLogisticRecord, args, "", "")
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(invoke)))
		return
	}

	middleware.ResponseSuccess(c, string(invoke))
}

// ReadProductInfo godoc
// @Summary 读取货物批次详情
// @Description 读取货物批次详情
// @Tags 零售商2-Retailer2
// @Accept  json
// @Produce  json
// @Param id path string true "批次编号"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /retailer2/product/{id} [get]
func (r *Retailer2Controller) ReadProductInfo(c *gin.Context) {
	batchNum := c.Param("id")
	if batchNum == "" {
		middleware.ResponseError(c, 2001, errors.New("please specify batch number"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRetailer2)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + batchNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadProductInfo, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(query)))
		return
	}

	resp := strings.ReplaceAll(string(query), `\`, ``)

	prodInfo := &dao.ProductInfo{}

	_ = json.Unmarshal([]byte(resp), prodInfo)

	middleware.ResponseSuccess(c, prodInfo)
}

// ReadTrackingResult godoc
// @Summary 读取物流详情
// @Description 读取物流详情
// @Tags 零售商2-Retailer2
// @Accept  json
// @Produce  json
// @Param id path string true "物流号"
// @Success 200 {object} middleware.Response{data=dao.LogisticsRecord} "success"
// @Router /retailer2/tracking/{id} [get]
func (r *Retailer2Controller) ReadTrackingResult(c *gin.Context) {
	trackingNum := c.Param("id")
	if trackingNum == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRetailer2)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + trackingNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadLogisticsRecord, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(query)))
		return
	}

	record := &dao.LogisticsRecord{}
	err = json.Unmarshal(query, record)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	middleware.ResponseSuccess(c, record)
}

// ReadOrderInfo godoc
// @Summary 读取订单详情
// @Description 读取订单详情
// @Tags 零售商2-Retailer2
// @Accept  json
// @Produce  json
// @Param id path string true "订单号"
// @Success 200 {object} middleware.Response{data=dao.OrderInfo} "success"
// @Router /retailer2/order/{id} [get]
func (r *Retailer2Controller) ReadOrderInfo(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		middleware.ResponseError(c, 2001, errors.New("order ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRetailer2)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	args := strings.Join([]string{
		`"` + orderID + `"`,
		`"` + fabsdk.CollectionTransaction2 + `"`,
	}, ",")

	query, err := sdkCtx.Query(fabsdk.FuncReadOrderInfo, args, false)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(query)))
		return
	}

	resp := strings.ReplaceAll(string(query), `\`, ``)
	orderInfo := &dao.OrderInfo{}
	_ = json.Unmarshal([]byte(resp), orderInfo)

	middleware.ResponseSuccess(c, orderInfo)
}

