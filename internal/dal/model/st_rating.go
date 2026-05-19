package model

import "time"

const TableNameStRating = "st_rating"

type StRating struct {
	ID        int        `gorm:"column:id;type:INTEGER;primaryKey" json:"id"`
	SiteID    int        `gorm:"column:site_id;type:int(11);not null;index" json:"site_id"`
	UserID    int        `gorm:"column:user_id;type:int(11);not null;default:0;index" json:"user_id"`
	VisitorID string     `gorm:"column:visitor_id;type:varchar(100);default:'';uniqueIndex:idx_site_visitor" json:"visitor_id"`
	Score     int        `gorm:"column:score;type:int(11);not null" json:"score"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP not null" json:"updated_at"`
}

func (*StRating) TableName() string {
	return TableNameStRating
}
