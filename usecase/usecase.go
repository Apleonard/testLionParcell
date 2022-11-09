package usecase

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"testLionParcell/models"
	"testLionParcell/repository"

	"github.com/labstack/gommon/log"
)

type Usecases interface {
	Upload(file multipart.FileHeader) error
	processUploadFile(content [][]string) *models.Payroll
}

type usecases struct {
	repo repository.Repositories
}

func NewUsecase(repo repository.Repositories) Usecases {
	return &usecases{
		repo: repo,
	}
}

func (u *usecases) Upload(file multipart.FileHeader) error {

	f, err := file.Open()
	if err != nil {
		log.Error("file open", err)
		return nil
	}
	defer f.Close()

	reader := csv.NewReader(f)
	content, err := reader.ReadAll()
	if err != nil {
		return err
	}

	test := u.processUploadFile(content)
	fmt.Println(test)

	return nil
}

func (u *usecases) processUploadFile(content [][]string) *models.Payroll {

	payroll := &models.Payroll{}

	for i, contentData := range content[1:] {
		for i, v := range contentData {
			switch i {
			case 0:
				batchInt, err := strconv.Atoi(v)
				if err != nil {
					return nil
				}
				payroll.Batch = int64(batchInt)
			case 1:
				payroll.AccountName = v
			case 2:
				payroll.AccountNumber = v
			case 3:
				userIDInt, err := strconv.Atoi(v)
				if err != nil {
					return nil
				}
				payroll.UserID = int64(userIDInt)
			case 4:
				AmountFloat64, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil
				}
				payroll.Amount = AmountFloat64
			case 5:
				payroll.Status = v
			}
		}
		fmt.Println(i, payroll)
		err := u.repo.CheckUser(payroll)
		if err != nil {
			fmt.Println("error check user excel row-", i, payroll)
		} else {
			u.repo.CreatePayroll(payroll)
		}
	}
	return payroll
}
