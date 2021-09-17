package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dao"
	"github.com/joexu01/logistics-app/fabsdk"
	"github.com/joexu01/logistics-app/middleware"
	"strings"
)

type RegulatorController struct{}

func RegulatorRegister(router *gin.RouterGroup) {
	r := &RegulatorController{}

	router.GET("/product/:id", r.ReadProductInfo)
	router.GET("/tracking/:id", r.ReadTrackingResult)
	router.GET("/private/:id", r.ReadCombinedTrackingResult)
	router.GET("/search/product/:keyword", r.SearchProductInfoByName)
}

// ReadProductInfo godoc
// @Summary 读取货物批次详情
// @Description 读取货物批次详情
// @Tags 监管机构-Regulator
// @Accept  json
// @Produce  json
// @Param id path string true "批次编号"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /regulator/product/{id} [get]
func (r *RegulatorController) ReadProductInfo(c *gin.Context) {
	batchNum := c.Param("id")
	if batchNum == "" {
		middleware.ResponseError(c, 2001, errors.New("please specify batch number"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRegulator)
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
// @Tags 监管机构-Regulator
// @Accept  json
// @Produce  json
// @Param id path string true "物流号"
// @Success 200 {object} middleware.Response{data=dao.LogisticsRecord} "success"
// @Router /regulator/tracking/{id} [get]
func (r *RegulatorController) ReadTrackingResult(c *gin.Context) {
	trackingNum := c.Param("id")
	if trackingNum == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRegulator)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	arg := `"` + trackingNum + `"`

	query, err := sdkCtx.Query(fabsdk.FuncReadLogisticsRecord, arg, true)
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

// ReadCombinedTrackingResult godoc
// @Summary 读取物流详情
// @Description 读取物流详情
// @Tags 监管机构-Regulator
// @Accept  json
// @Produce  json
// @Param id path string true "物流号"
// @Success 200 {object} middleware.Response{data=dao.LogisticsCombinedRecord} "success"
// @Router /regulator/private/{id} [get]
func (r *RegulatorController) ReadCombinedTrackingResult(c *gin.Context) {
	trackingNum := c.Param("id")
	if trackingNum == "" {
		middleware.ResponseError(c, 2001, errors.New("tracking ID is empty string"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRegulator)
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
// @Tags 监管机构-Regulator
// @Accept  json
// @Produce  json
// @Param keyword path string true "关键词-商品名称"
// @Success 200 {object} middleware.Response{data=dao.ProductInfo} "success"
// @Router /regulator/search/product/{keyword} [get]
func (r *RegulatorController) SearchProductInfoByName(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		middleware.ResponseError(c, 2001, errors.New("keyword is empty"))
		return
	}

	sdkCtx, err := fabsdk.NewFabSDKCtx(fabsdk.NUMRegulator)
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
