package model

import "time"

const TableNameStTag = "st_tag"

type StTag struct {
	ID        int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Color     string     `gorm:"column:color;type:varchar(20);default:'#3b82f6'" json:"color"`
	Sort      int        `gorm:"column:sort;type:int(11);default:0" json:"sort"`
	IsUsed    bool       `gorm:"column:is_used;type:bool;default:true" json:"is_used"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StTag) TableName() string {
	return TableNameStTag
}

const TableNameStSiteTag = "st_site_tag"

type StSiteTag struct {
	SiteID int `gorm:"column:site_id;type:int(11);primaryKey" json:"site_id"`
	TagID  int `gorm:"column:tag_id;type:int(11);primaryKey" json:"tag_id"`
}

func (*StSiteTag) TableName() string {
	return TableNameStSiteTag
}
