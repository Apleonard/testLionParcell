package repository

import (
	"errors"
	"testLionParcell/models"

	"github.com/jinzhu/gorm"
)

type Repositories interface {
	CreatePayroll(data *models.Payroll) error
	CheckUser(data *models.Payroll) error
	CreatePayrollLOg(data *models.PayrollLog) error
}

type repositories struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repositories {
	return &repositories{
		db: db,
	}
}

func (r *repositories) CreatePayroll(data *models.Payroll) error {
	r.db.Create(data)

	return nil
}

func (r *repositories) CheckUser(data *models.Payroll) error {
	result := &models.User{}
	user := &models.User{}
	user.ID = uint(data.UserID)
	user.Name = data.AccountName
	user.Status = data.Status

	err := r.db.Where("id = ?", user.ID).First(&result).Error
	if err != nil {
		return err
	}

	if user.ID != result.ID || user.Name != result.Name || user.Status != result.Status {
		return errors.New("user doesnt match")
	}

	return nil
}

func (r *repositories) CreatePayrollLOg(data *models.PayrollLog) error {
	r.db.Create(data)
	return nil
}
