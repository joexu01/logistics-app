package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/dao"
	"github.com/joexu01/logistics-app/dto"
	"github.com/joexu01/logistics-app/fabsdk"
	"github.com/joexu01/logistics-app/middleware"
	"github.com/joexu01/logistics-app/public"
	"github.com/skip2/go-qrcode"
	"os"
	"strings"
)

type QRCodeController struct{}

func QRCodeRegister(group *gin.RouterGroup) {
	c := QRCodeController{}

	group.Static("/file", public.QRCodeDir)
	group.GET("/order/:order_id", c.GenerateRelevantQRCodeByOrderID)
	group.GET("/default", c.GetDefaultImage)
}

// GetDefaultImage godoc
// @Summary 获取默认图像
// @Description 获取默认图像
// @Tags 二维码
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.QRCodeImagesURL} "success"
// @Router /qrcode/default [get]
func (q *QRCodeController) GetDefaultImage(c *gin.Context) {
	url := "http://" + public.LANIPAddr + ":8880/qrcode/file/" + "default.png"
	imagesURL := &dto.QRCodeImagesURL{
		QueryProdInfoURL:         url,
		QueryLogisticsRecordURL:  url,
		UpdateLogisticsStatusUrl: url,
	}

	middleware.ResponseSuccess(c, imagesURL)
}

// GenerateRelevantQRCodeByOrderID godoc
// @Summary 获取二维码
// @Description 根据订单号获取全部二维码
// @Tags 二维码
// @Accept  json
// @Produce  json
// @Param id path string true "订单号"
// @Param collection query string false "Collection名称"
// @Success 200 {object} middleware.Response{data=dto.QRCodeImagesURL} "success"
// @Router /qrcode/order/{id} [get]
func (q *QRCodeController) GenerateRelevantQRCodeByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		middleware.ResponseError(c, 2001, errors.New("orderId order_id is empty"))
		return
	}

	collectionName := c.DefaultQuery("collection", fabsdk.CollectionTransaction1)

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

	if _, err := os.Open(public.QRCodeDir); err != nil {
		_ = os.Mkdir(`static`, os.ModePerm)
		_ = os.Mkdir(public.QRCodeDir, os.ModePerm)
	}

	queryProdInfo := `http://` + public.LANIPAddr + public.PubQueryProductInfo + orderInfo.BatchNumber
	queryTracking := `http://` + public.LANIPAddr + public.PubQueryLogisticsRecord + orderInfo.TrackingNumber
	updateStatus := `http://` + public.LANIPAddr + public.LogisticsUpdateStatus + orderInfo.TrackingNumber

	queryProdInfoFilename := orderID + `_prod.png`
	queryTrackingFilename := orderID + `_tracking.png`
	updateStatusFilename := orderID + `_logistics.png`

	_ = qrcode.WriteFile(
		queryProdInfo, qrcode.Medium, 256, public.QRCodeDir+`/`+queryProdInfoFilename)
	_ = qrcode.WriteFile(
		queryTracking, qrcode.Medium, 256, public.QRCodeDir+`/`+queryTrackingFilename)
	_ = qrcode.WriteFile(
		updateStatus, qrcode.Medium, 256, public.QRCodeDir+`/`+updateStatusFilename)

	imagesURL := &dto.QRCodeImagesURL{
		QueryProdInfoURL:         "http://" + public.LANIPAddr + ":8880/qrcode/file/" + queryProdInfoFilename,
		QueryLogisticsRecordURL:  "http://" + public.LANIPAddr + ":8880/qrcode/file/" + queryTrackingFilename,
		UpdateLogisticsStatusUrl: "http://" + public.LANIPAddr + ":8880/qrcode/file/" + updateStatusFilename,
	}

	middleware.ResponseSuccess(c, imagesURL)
}
