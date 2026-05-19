package model

import "time"

const TableNameStPaymentMethod = "st_payment_method"

type StPaymentMethod struct {
	ID        int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Icon      string     `gorm:"column:icon;type:text" json:"icon"`
	Sort      int        `gorm:"column:sort;type:int(11);default:0" json:"sort"`
	IsUsed    bool       `gorm:"column:is_used;type:bool;default:true" json:"is_used"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StPaymentMethod) TableName() string {
	return TableNameStPaymentMethod
}

const TableNameStSitePayment = "st_site_payment"

type StSitePayment struct {
	SiteID          int `gorm:"column:site_id;type:int(11);primaryKey" json:"site_id"`
	PaymentMethodID int `gorm:"column:payment_method_id;type:int(11);primaryKey" json:"payment_method_id"`
}

func (*StSitePayment) TableName() string {
	return TableNameStSitePayment
}
