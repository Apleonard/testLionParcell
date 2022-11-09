package models

import "github.com/jinzhu/gorm"

type User struct {
	ID     uint
	Name   string `gorm:"size:255;not null" json:"name"`
	Status string `gorm:"size:10;not null" json:"status"`
}

type Payroll struct {
	Batch         int64   `gorm:"not null" json:"batch"`
	AccountName   string  `gorm:"size:255;not null" json:"account_name"`
	AccountNumber string  `gorm:"size:255;unique;not null" json:"account_number"`
	UserID        int64   `gorm:"not null" json:"user_id"`
	Amount        float64 `gorm:"not null" json:"amount"`
	Status        string  `gorm:"size:255;not null" json:"status"`
}

type PayrollLog struct {
	FileName     string `gorm:"not null" json:"file_name"`
	Batch        int64  `gorm:"size:255;not null" json:"batch"`
	TotalSuccess int64  `gorm:"not null" json:"total_success"`
	TotalFailed  int64  `gorm:"not null" json:"total_failed"`
	gorm.Model
}
