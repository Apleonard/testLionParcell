package repository

import (
	"fmt"
	"testLionParcell/models"

	"github.com/jinzhu/gorm"
)

type Repositories interface {
	CreatePayroll(data *models.Payroll) error
	CheckUser(data *models.Payroll) error
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
	user := &models.User{}
	user.ID = uint(data.UserID)
	user.Name = data.AccountName
	user.Status = data.Status

	err := r.db.First(&user).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
