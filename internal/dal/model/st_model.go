package model

import "time"

const TableNameStModel = "st_model"

type StModel struct {
	ID        int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Provider  string     `gorm:"column:provider;type:varchar(50)" json:"provider"`
	Icon      string     `gorm:"column:icon;type:text" json:"icon"`
	Sort      int        `gorm:"column:sort;type:int(11);default:0" json:"sort"`
	IsUsed    bool       `gorm:"column:is_used;type:bool;default:true" json:"is_used"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StModel) TableName() string {
	return TableNameStModel
}
