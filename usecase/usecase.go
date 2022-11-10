package usecase

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"sync"
	"testLionParcell/models"
	"testLionParcell/repository"
	"time"

	"github.com/labstack/gommon/log"
)

type Usecases interface {
	Upload(file multipart.FileHeader) error
	processUploadFile(content [][]string, file multipart.FileHeader) *models.Payroll
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

	fmt.Println("masuk usecase")

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
	u.processUploadFile(content, file)

	return nil
}

func (u *usecases) processUploadFile(content [][]string, file multipart.FileHeader) *models.Payroll {

	payroll := &models.Payroll{}
	payrollLog := models.PayrollLog{}
	success := int64(0)
	fail := int64(0)

	type Fail struct {
		Name         string
		ErrorMessage string
	}

	chanCreateInternalProviderFail := make(chan *Fail)
	defer close(chanCreateInternalProviderFail)

	var wg sync.WaitGroup

	go func() {
		for i, contentData := range content[1:] {
			wg.Add(1)
			for i, v := range contentData {
				switch i {
				case 0:
					batchInt, err := strconv.Atoi(v)
					if err != nil {
						// return nil
					}
					payroll.Batch = int64(batchInt)
				case 1:
					payroll.AccountName = v
				case 2:
					payroll.AccountNumber = v
				case 3:
					userIDInt, err := strconv.Atoi(v)
					if err != nil {
						// return nil
					}
					payroll.UserID = int64(userIDInt)
				case 4:
					AmountFloat64, err := strconv.ParseFloat(v, 64)
					if err != nil {
						// return nil
					}
					payroll.Amount = AmountFloat64
				case 5:
					payroll.Status = v
				}
			}
			err := u.repo.CheckUser(payroll)
			if err != nil {
				success++
				fmt.Println("error check user excel row-", i, payroll)
			} else {
				fail++
				u.repo.CreatePayroll(payroll)
			}

			//set payroll log
			payrollLog.Batch = payroll.Batch
			payrollLog.TotalSuccess = success
			payrollLog.TotalFailed = fail
			payrollLog.CreatedAt = time.Now()
			payrollLog.UpdatedAt = time.Now()
			payrollLog.FileName = file.Filename
		}

		err := u.repo.CreatePayrollLOg(&payrollLog)
		if err != nil {
			// return nil
		}

	}()
	wg.Wait()

	return payroll
}
