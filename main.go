package main

import (
	"fmt"
	"net/http"
	"os"

	"testLionParcell/handler"
	"testLionParcell/models"
	"testLionParcell/repository"
	"testLionParcell/usecase"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
	e  error
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}

	var (
		dbName = os.Getenv("DB_NAME")
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASS")
		dbPort = os.Getenv("DB_PORT")
	)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbUser,
		dbPass,
		dbName,
		dbPort)

	db, e = gorm.Open("postgres", dsn)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("Connection Established")
	}
	defer db.Close()
	db.AutoMigrate(
		&models.User{},
		&models.Payroll{},
		&models.PayrollLog{},
	)

	repo := repository.NewRepository(db)
	useCase := usecase.NewUsecase(repo)
	handlers := handler.NewHandler(useCase)

	router := mux.NewRouter()
	router.HandleFunc("/check", handlers.Check).Methods("GET")

	router.HandleFunc("/upload", handlers.Upload).Methods("POST")
	fmt.Println("run at port:8080")
	http.ListenAndServe(":8080", router)
}
