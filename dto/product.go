package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/joexu01/logistics-app/public"
)

type ProductInfoInput struct {
	BatchNum string `json:"batch_num" form:"batch_num" comment:"产品批次"`
	Amount   int    `json:"amount" form:"amount" comment:"数量"`
	Origin   string `json:"origin" form:"origin" comment:"产地"`
	Name     string `json:"name" form:"name" comment:"姓名"`
}

func (params *ProductInfoInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type OrderInput struct {
	OrderNumber    string  `json:"order_number" form:"order_number"`
	BatchNumber    string  `json:"batch_number" form:"batch_number"`
	TrackingNumber string  `json:"tracking_number" form:"tracking_number"`
	Sorter         string  `json:"sorter" form:"sorter"` // 分拣员
	UnitPrice      float32 `json:"unit_price" form:"unit_price"`
	Quantity       int     `json:"quantity" form:"quantity"`
	Collection     string  `json:"collection" form:"collection"`
}

func (params *OrderInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type UpdateLogisticRecordInput struct {
	TrackingNum string `json:"tracking_num" form:"tracking_num"`
	Status      string `json:"status" form:"status"`
	Operator    string `json:"operator" form:"operator"`
}

func (params *UpdateLogisticRecordInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

type ReadOrderInfoInput struct {
	OrderNum       string `json:"order_num" form:"order_num"`
	CollectionName string `json:"collection_name" form:"collection_name"`
}

func (params *ReadOrderInfoInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
