package model

import "time"

const TableNameStAd = "st_ad"

type StAd struct {
	ID          int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Slot        string     `gorm:"column:slot;type:varchar(50);not null" json:"slot"`
	ContentType string     `gorm:"column:content_type;type:varchar(20);default:'html'" json:"content_type"`
	Content     string     `gorm:"column:content;type:text" json:"content"`
	ImageURL    string     `gorm:"column:image_url;type:varchar(500)" json:"image_url"`
	LinkURL     string     `gorm:"column:link_url;type:varchar(500)" json:"link_url"`
	StartAt     *time.Time `gorm:"column:start_at;type:datetime" json:"start_at"`
	EndAt       *time.Time `gorm:"column:end_at;type:datetime" json:"end_at"`
	Sort        int        `gorm:"column:sort;type:int(11);default:0" json:"sort"`
	IsUsed      bool       `gorm:"column:is_used;type:bool;default:true" json:"is_used"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StAd) TableName() string {
	return TableNameStAd
}
