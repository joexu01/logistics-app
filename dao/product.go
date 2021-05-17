package dao

type ProductInfo struct {
	Amount       int    `json:"amount"`
	Origin       string `json:"origin"`
	Name         string `json:"name"`
	LastModified string `json:"last_modified"`
	BatchNumber  string `json:"batch_number"`
}

type OrderInfo struct {
	OrderID string `json:"order_id"`
	//OrderNumber    string  `json:"order_number"`
	BatchNumber    string  `json:"batch_number"`
	TrackingNumber string  `json:"tracking_number"`
	Sorter         string  `json:"sorter"` // 分拣员
	UnitPrice      float32 `json:"unit_price"`
	Quantity       int     `json:"quantity"`
	Client         string  `json:"client"`
	Status         string  `json:"status"`
	ProductName    string  `json:"product_name"`
}

type LogisticsRecord struct {
	Items []RecordSubItem `json:"items"`
}

type RecordSubItem struct {
	RecordTime string `json:"record_time"`
	Status     string `json:"status"`
}

type PrivateLogisticsRecord struct {
	Items []PrivateSubItem `json:"items"`
}

type PrivateSubItem struct {
	RecordTime string `json:"record_time"`
	PeerID     string `json:"peer_id"`
	Operator   string `json:"operator"`
}

type LogisticsCombinedRecord struct {
	Public  *LogisticsRecord        `json:"public"`
	Private *PrivateLogisticsRecord `json:"private"`
}
