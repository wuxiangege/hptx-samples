package dao

import "gitee.com/chunanyong/zorm"

type SoItemTbl struct {
	zorm.EntityStruct

	SysNo         int64   `column:"sysno"`
	SoSysNo       int64   `column:"so_sysno"`
	ProductSysNo  int64   `column:"product_sysno"`
	ProductName   string  `column:"product_name"`
	CostPrice     float64 `column:"cost_price"`
	OriginalPrice float64 `column:"original_price"`
	DealPrice     float64 `column:"deal_price"`
	Quantity      int32   `column:"quantity"`
}

func (entity *SoItemTbl) GetTableName() string {
	return "order.so_item"
}

func (entity *SoItemTbl) GetPKColumnName() string {
	return "sysno"
}
