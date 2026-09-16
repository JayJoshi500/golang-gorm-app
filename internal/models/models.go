package models

import (
	"time"
)

type RestoTables struct {
	ID        int64 `gorm:"type:int;primaryKey" json:"id"`
	Capacity  int64 `gorm:"type:int;not null" json:"capacity"`
	BookedVal int64 `gorm:"type:int;not null" json:"booked_val"`
}

type Booking struct {
	ID          int64     `gorm:"type:int;primaryKey" json:"id"`
	RestoTables int64     `gorm:"foreignKey:CompanyRefer;references:Code"`
	StartTime   time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime     time.Time `gorm:"type:timestamp;not null" json:"end_time"`
}

type Customer struct {
	ID        int64  `gorm:"type:int;primaryKey" json:"id"`
	FullName  string `gorm:"type:varchar(150);not null" json:"full_name"`
	ContactNo string `gorm:"type:varchar(150);not null" json:"contact_no"`
}

type RestoTimingSlots struct {
	ID        int64     `gorm:"type:int;primaryKey" json:"id"`
	DayOfWeek string    `gorm:"type:varchar(150);not null" json:"day_of_week"`
	StartTime time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime   time.Time `gorm:"type:timestamp;not null" json:"end_time"`
}
