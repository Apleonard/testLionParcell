package usecase

import (
	"encoding/csv"
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
	BatchUpload(file multipart.FileHeader) error
	processUploadFile(content [][]string, file multipart.FileHeader) *models.Payroll
	processBatchUploadFile(content [][]string, file multipart.FileHeader) *models.Payroll
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
		return err
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

	chanFail := make(chan *Fail)
	defer close(chanFail)

	var wg sync.WaitGroup

	go func() {
		defer wg.Done()
		for _, contentData := range content[1:] {
			wg.Add(1)
			// worker(i)
			for i, v := range contentData {
				switch i {
				case 0:
					batchInt, err := strconv.Atoi(v)
					if err != nil {
						chanFail <- &Fail{}
					}
					payroll.Batch = int64(batchInt)
				case 1:
					payroll.AccountName = v
				case 2:
					payroll.AccountNumber = v
				case 3:
					userIDInt, err := strconv.Atoi(v)
					if err != nil {
						chanFail <- &Fail{}
					}
					payroll.UserID = int64(userIDInt)
				case 4:
					AmountFloat64, err := strconv.ParseFloat(v, 64)
					if err != nil {
						chanFail <- &Fail{}
					}
					payroll.Amount = AmountFloat64
				case 5:
					payroll.Status = v
				}
			}
			err := u.repo.CheckUser(payroll)
			if err != nil {
				fail++
			} else {
				success++
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
			return
		}

	}()
	wg.Wait()

	return payroll
}

func (u *usecases) BatchUpload(file multipart.FileHeader) error {
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

	u.processBatchUploadFile(content, file)

	return nil
}

func (u *usecases) processBatchUploadFile(content [][]string, file multipart.FileHeader) *models.Payroll {

	var (
		payrollLog  models.PayrollLog
		datas       []interface{}
		FailedDatas []interface{}
		success     int64
		fail        int64
	)

	for i := 1; i < len(content); i++ {
		batchInt, _ := strconv.Atoi(content[i][0])
		userIDInt, _ := strconv.Atoi(content[i][3])
		amountFloat64, _ := strconv.ParseFloat(content[i][4], 64)

		data := &models.Payroll{
			Batch:         int64(batchInt),
			UserID:        int64(userIDInt),
			AccountName:   content[i][1],
			AccountNumber: content[i][2],
			Amount:        amountFloat64,
			Status:        content[i][5],
		}

		//check to users db
		err := u.repo.CheckUser(data)
		if err != nil {
			//append failed data if error
			fail++
			FailedData := &models.PayrollFail{
				Batch:         data.Batch,
				AccountName:   data.AccountName,
				AccountNumber: data.AccountNumber,
				UserID:        data.UserID,
				Amount:        data.Amount,
				Status:        data.Status,
			}
			FailedDatas = append(FailedDatas, FailedData)
		} else {
			//append datas if success
			success++
			datas = append(datas, data)
		}
		payrollLog.Batch = data.Batch
	}

	//insert to payrolls db
	err := u.repo.CreateBatchPayroll(datas)
	if err != nil {
		return nil
	}

	//insert to payroll failed log db
	err = u.repo.CreateBatchFailedPayroll(FailedDatas)
	if err != nil {
		return nil
	}

	//set payroll log
	payrollLog.TotalSuccess = success
	payrollLog.TotalFailed = fail
	payrollLog.CreatedAt = time.Now()
	payrollLog.UpdatedAt = time.Now()
	payrollLog.FileName = file.Filename

	//insert to payroll logs db
	err = u.repo.CreatePayrollLOg(&payrollLog)
	if err != nil {
		return nil
	}

	return nil
}
