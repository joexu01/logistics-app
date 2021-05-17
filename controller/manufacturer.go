package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dao"
	"github.com/joexu01/logistics-app/dto"
	"github.com/joexu01/logistics-app/fabsdk"
	"github.com/joexu01/logistics-app/middleware"
	"strconv"
	"strings"
)

type ManufacturerController struct{}

func ManufacturerRegister(router *gin.RouterGroup) {
	m := &ManufacturerController{}
	router.POST("/product", m.CreateProductInfo)
	router.POST("/send", m.SendProductOff)
	router.PUT("/order/accept/:id", m.AcceptOrder)

	router.GET("/product/:id", m.ReadProductInfo)
	router.GET("/order/read/:id", m.ReadOrderInfo)
	router.GET("/tracking/:id", m.ReadTrackingResult)
	router.GET("/order/unaccepted", m.ReadUnacceptedOrders)
	router.GET("/order/reject/:id", m.RejectOrder)
	router.GET("/search/product/:keyword", m.SearchProductInfoByName)
}

// CreateProductInfo godoc
// @Summary 商品信息上链
// @Description 商品信息上链
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param body body dto.ProductInfoInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /manufacturer/product [post]
func (m *ManufacturerController) CreateProductInfo(c *gin.Context) {
	info := &dto.ProductInfoInput{}
	if err := info.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	amount := strconv.Itoa(info.Amount)

	args := strings.Join([]string{
		`"` + info.BatchNum + `"`,
		`"` + amount + `"`,
		`"` + info.Name + `"`,
		`"` + info.Origin + `"`}, ",")

	bytes, err := sdkCtx.Invoke(fabsdk.FuncNewProductInfo, args, "", "")
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(bytes)))
		return
	}

	middleware.ResponseSuccess(c, string(bytes))
}

// SendProductOff godoc
// @Summary 发货记录上链
// @Description 发货记录上链
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param body body dto.UpdateLogisticRecordInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /manufacturer/send [post]
func (m *ManufacturerController) SendProductOff(c *gin.Context) {
	//logisticsID, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	middleware.ResponseError(c, 2001, err)
	//	return
	//}

	input := &dto.UpdateLogisticRecordInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	args := strings.Join([]string{
		`"` + input.TrackingNum + `"`,
		`"` + `制造商-Manufacturer: ` + input.Status + `"`,
	}, ",")

	invoke, err := sdkCtx.Invoke(fabsdk.FuncUpdateLogisticRecord, args, "", "")
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(invoke)))
		return
	}
	resp := strings.ReplaceAll(string(invoke), `\`, ``)

	middleware.ResponseSuccess(c, resp)
}

// ReadProductInfo godoc
// @Summary 读取货物批次详情
// @Description 读取货物批次详情
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param id path string true "批次编号"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /manufacturer/product/{id} [get]
func (m *ManufacturerController) ReadProductInfo(c *gin.Context) {
	batchNum := c.Param("id")
	if batchNum == "" {
		middleware.ResponseError(c, 2001, errors.New("please specify batch number"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + batchNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadProductInfo, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	resp := strings.ReplaceAll(string(query), `\`, ``)

	prodInfo := &dao.ProductInfo{}

	_ = json.Unmarshal([]byte(resp), prodInfo)

	middleware.ResponseSuccess(c, prodInfo)
}

// ReadOrderInfo godoc
// @Summary 读取订单详情
// @Description 读取订单详情
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param id path string true "订单号"
// @Param collection query string false "Collection名称"
// @Success 200 {object} middleware.Response{data=dao.OrderInfo} "success"
// @Router /manufacturer/order/read/{id} [get]
func (m *ManufacturerController) ReadOrderInfo(c *gin.Context) {
	collectionName := c.DefaultQuery("collection", fabsdk.CollectionTransaction1)
	orderID := c.Param("id")
	if orderID == "" {
		middleware.ResponseError(c, 2001, errors.New("order ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	args := strings.Join([]string{
		`"` + orderID + `"`,
		`"` + collectionName + `"`,
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

// ReadTrackingResult godoc
// @Summary 读取物流详情
// @Description 读取物流详情
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param id path string true "物流号"
// @Success 200 {object} middleware.Response{data=dao.LogisticsRecord} "success"
// @Router /manufacturer/tracking/{id} [get]
func (m *ManufacturerController) ReadTrackingResult(c *gin.Context) {
	trackingNum := c.Param("id")
	if trackingNum == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
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

// ReadUnacceptedOrders godoc
// @Summary 读取接受到的订单请求
// @Description 读取接受到的订单请求
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param collection query string false "Collection名称"
// @Success 200 {object} middleware.Response{data=dao.OrderInfo} "success"
// @Router /manufacturer/order/unaccepted [get]
func (m *ManufacturerController) ReadUnacceptedOrders(c *gin.Context) {
	collectionName := c.DefaultQuery("collection", fabsdk.CollectionTransaction1)

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	arg := `"` + collectionName + `"`

	query, err := sdkCtx.Query(fabsdk.FuncGetOrdersUnaccepted, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New(string(query)))
		return
	}

	result := strings.ReplaceAll(string(query), "\n", "")
	result = strings.ReplaceAll(result, `\`, ``)

	var orders []dao.OrderInfo
	err = json.Unmarshal([]byte(result), &orders)
	if err != nil {
		middleware.ResponseError(c, 2003, fmt.Errorf("error unmarshaling JSON: %v", err))
		return
	}

	middleware.ResponseSuccess(c, orders)
}

// AcceptOrder godoc
// @Summary 接受订单
// @Description 接受订单
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param id path string true "订单ID"
// @Param collection query string false "Collection名称"
// @Param body body dto.AcceptOrderInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /manufacturer/order/accept/{id} [put]
func (m *ManufacturerController) AcceptOrder(c *gin.Context) {
	input := &dto.AcceptOrderInput{}
	err := input.BindingValidParams(c)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	collection := c.DefaultQuery("collection", "transactionCollection1")
	orderID := c.Param("id")
	if orderID == "" {
		middleware.ResponseError(c, 2002, errors.New("orderID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	args := strings.Join([]string{
		`"` + collection + `"`,
		`"` + orderID + `"`,
	}, ",")

	bytes, err := json.Marshal(input)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	transient := base64.StdEncoding.EncodeToString(bytes)

	invoke, err := sdkCtx.Invoke(fabsdk.FuncAcceptOrder, args, fabsdk.TransientKeyAcceptOrderInput, transient)

	if err != nil {
		middleware.ResponseError(c, 2005,
			fmt.Errorf("error invoking function %s: %s", fabsdk.FuncAcceptOrder, invoke))
		return
	}

	middleware.ResponseSuccess(c, string(invoke))
}

// RejectOrder godoc
// @Summary 拒绝订单
// @Description 拒绝订单
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param id path string true "订单ID"
// @Param collection query string true "Collection名称"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /manufacturer/order/reject/{id} [get]
func (m *ManufacturerController) RejectOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		middleware.ResponseError(c, 2001, errors.New("order ID is empty"))
		return
	}
	collection := c.Query("collection")
	if collection == "" {
		middleware.ResponseError(c, 2002, errors.New("collection name is empty"))
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	args := strings.Join([]string{
		`"` + collection + `"`,
		`"` + orderID + `"`,
	}, ",")

	invoke, err := sdkCtx.Invoke(fabsdk.FuncRejectOrderRequest, args, "", "")
	if err != nil {
		middleware.ResponseError(c, 2004, errors.New(string(invoke)))
		return
	}

	middleware.ResponseSuccess(c, string(invoke))
}

// SearchProductInfoByName godoc
// @Summary 根据货物名称搜索批次详情
// @Description 根据货物名称搜索批次详情
// @Tags 制造商-Manufacturer
// @Accept  json
// @Produce  json
// @Param keyword path string true "关键词-商品名称"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /manufacturer/search/product/{keyword} [get]
func (m *ManufacturerController) SearchProductInfoByName(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		middleware.ResponseError(c, 2001, errors.New("keyword is empty"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMManufacturer)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + keyword + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadProductInfoByProductName, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(query)))
		return
	}

	result := strings.ReplaceAll(string(query), "\n", "")
	result = strings.ReplaceAll(result, `\`, ``)

	var products []dao.ProductInfo
	err = json.Unmarshal([]byte(result), &products)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	middleware.ResponseSuccess(c, products)
}
