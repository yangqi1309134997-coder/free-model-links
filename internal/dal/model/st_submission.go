package model

import "time"

const TableNameStSubmission = "st_submission"

type StSubmission struct {
	ID          int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	UserID      int        `gorm:"column:user_id;type:int(11);default:0" json:"user_id"`
	URL         string     `gorm:"column:url;type:varchar(255);not null" json:"url"`
	Title       string     `gorm:"column:title;type:varchar(100)" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	CategoryID  int        `gorm:"column:category_id;type:int(11)" json:"category_id"`
	Contact     string     `gorm:"column:contact;type:varchar(255)" json:"contact"`
	Status      string     `gorm:"column:status;type:varchar(20);default:'pending'" json:"status"`
	AdminNote   string     `gorm:"column:admin_note;type:text" json:"admin_note"`
	CreatedAt   *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

func (*StSubmission) TableName() string {
	return TableNameStSubmission
}
