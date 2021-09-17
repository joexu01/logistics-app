package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dao"
	"github.com/joexu01/logistics-app/dto"
	"github.com/joexu01/logistics-app/fabsdk"
	"github.com/joexu01/logistics-app/middleware"
	"strings"
)

type LogisticsController struct{}

func LogisticsRegister(router *gin.RouterGroup) {
	l := &LogisticsController{}
	router.POST("/update/:id", l.UpdateLogisticsRecord)

	router.GET("/product/:id", l.ReadProductInfo)
	router.GET("/tracking/:id", l.ReadTrackingResult)
	router.GET("/search/product/:keyword", l.SearchProductInfoByName)
}

// UpdateLogisticsRecord godoc
// @Summary 物流记录上链
// @Description 物流记录上链
// @Tags 物流企业-Logistics
// @Accept  json
// @Produce  json
// @Param body body dto.UpdateLogisticRecordInput true "body"
// @Param id path string true "物流追踪ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /logistics/update/{id} [post]
func (l *LogisticsController) UpdateLogisticsRecord(c *gin.Context) {
	input := &dto.UpdateLogisticRecordInput{}
	if err := input.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	trackingID:= c.Param("id")
	if trackingID == "" {
		middleware.ResponseError(c, 2001, errors.New("empty tracking ID"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMLogistics)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	args := strings.Join([]string{
		`"` + trackingID + `"`,
		`"` + `物流-Logistics Company: ` + input.Status + `"`,
	}, ",")

	transient := base64.StdEncoding.EncodeToString([]byte(input.Operator))

	invoke, err := sdkCtx.Invoke(
		fabsdk.FuncUpdateLogisticRecord, args, fabsdk.TransientKeyLogisticOperatorInput, transient)
	if err != nil {
		middleware.ResponseError(c, 2002, errors.New(string(invoke)))
		return
	}

	middleware.ResponseSuccess(c, string(invoke))
}

// ReadProductInfo godoc
// @Summary 读取货物批次详情
// @Description 读取货物批次详情
// @Tags 物流企业-Logistics
// @Accept  json
// @Produce  json
// @Param id path string true "批次编号"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /logistics/product/{id} [get]
func (l *LogisticsController) ReadProductInfo(c *gin.Context) {
	batchNum := c.Param("id")
	if batchNum == "" {
		middleware.ResponseError(c, 2001, errors.New("please specify batch number"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMLogistics)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + batchNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadProductInfo, arg, true)
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
// @Tags 物流企业-Logistics
// @Accept  json
// @Produce  json
// @Param id path string true "物流号"
// @Success 200 {object} middleware.Response{data=dao.LogisticsCombinedRecord} "success"
// @Router /logistics/tracking/{id} [get]
func (l *LogisticsController) ReadTrackingResult(c *gin.Context) {
	trackingNum := c.Param("id")
	if trackingNum == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMLogistics)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + trackingNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadLogisticsPriRecord, arg, false)
	if err != nil {
		middleware.ResponseError(c, 2003, errors.New(string(query)))
		return
	}

	resp := strings.ReplaceAll(string(query), `\`, ``)

	combineResult := &dao.LogisticsCombinedRecord{}
	_ = json.Unmarshal([]byte(resp), combineResult)

	middleware.ResponseSuccess(c, combineResult)
}

// SearchProductInfoByName godoc
// @Summary 根据货物名称搜索批次详情
// @Description 根据货物名称搜索批次详情
// @Tags 物流企业-Logistics
// @Accept  json
// @Produce  json
// @Param keyword path string true "关键词-商品名称"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /logistics/search/product/{keyword} [get]
func (l *LogisticsController) SearchProductInfoByName(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		middleware.ResponseError(c, 2001, errors.New("keyword is empty"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMLogistics)
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
