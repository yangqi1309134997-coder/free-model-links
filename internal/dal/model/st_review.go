package model

import "time"

const TableNameStReview = "st_review"

type StReview struct {
	ID        int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	SiteID    int        `gorm:"column:site_id;type:int(11);not null;index" json:"site_id"`
	UserID    int        `gorm:"column:user_id;type:int(11);not null;default:0;index" json:"user_id"`
	Nickname  string     `gorm:"column:nickname;type:varchar(100);default:'Anonymous'" json:"nickname"`
	Content   string     `gorm:"column:content;type:text;not null" json:"content"`
	Status    string     `gorm:"column:status;type:varchar(20);default:'pending'" json:"status"`
	AdminNote string     `gorm:"column:admin_note;type:text" json:"admin_note"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StReview) TableName() string {
	return TableNameStReview
}
