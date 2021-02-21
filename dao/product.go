package dao

import "time"

type ProductInfo struct {
	Amount       int       `json:"amount"`
	Origin       string    `json:"origin"`
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
}

type OrderInfo struct {
	//OrderNumber    string  `json:"order_number"`
	BatchNumber    string  `json:"batch_number"`
	TrackingNumber string  `json:"tracking_number"`
	Sorter         string  `json:"sorter"` // 分拣员
	UnitPrice      float32 `json:"unit_price"`
	Quantity       int     `json:"quantity"`
	Client         string  `json:"client"`
}

type LogisticsRecord struct {
	Items []RecordSubItem `json:"items"`
}

type RecordSubItem struct {
	RecordTime time.Time `json:"record_time"`
	Status     string    `json:"status"`
}

type PrivateLogisticsRecord struct {
	Items []PrivateSubItem `json:"items"`
}

type PrivateSubItem struct {
	RecordTime time.Time `json:"record_time"`
	PeerID     string    `json:"peer_id"`
	Operator   string    `json:"operator"`
}

type LogisticsCombinedRecord struct {
	Public  *LogisticsRecord        `json:"public"`
	Private *PrivateLogisticsRecord `json:"private"`
}
