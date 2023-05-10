package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	Port     = 5432
	Username = "psnko"
	Password = "postgres"
	DBName   = "mypgdb"
)

type repo struct {
	DB *sqlx.DB
}

func NewRepo() *repo {
	psqlConn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", Username, Password, DBName)
	db, err := sqlx.Open("postgres", psqlConn)
	if err != nil {
		log.Fatal("Failed to connect to db" + err.Error())
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Failed to ping db")
	}
	return &repo{DB: db}
}

func (r *repo) SetLogin(userID int64, serviceName string, pwd string) error {
	if err := r.DB.QueryRow("select password from passwords where id=$1 and login=$2", userID, serviceName).Scan(&pwd); err == nil {
		return err
	}
	_, err := r.DB.Exec("INSERT INTO passwords (id,login,password,created_at) VALUES ($1,$2,$3,$4)", userID, serviceName, pwd, time.Now())
	if err == nil {
		err = errors.New("ok")
	}
	return err
}

func (r *repo) GetLogin(userID int64, serviceName string) (string, error) {
	var pwd string
	if err := r.DB.QueryRow("select password from passwords where id=$1 and login=$2", userID, serviceName).Scan(&pwd); err == nil {
		return pwd, err
	} else {
		return "", err
	}
}

func (r *repo) Clear() {
	//current := time.Now()
	_, err := r.DB.Exec("delete from urls where created_at>(CURRENT_TIMESTAMP - interval '1 day';")
	if err != nil {
		log.Println(fmt.Sprintf("Failed to clear memory. Error: %v", err))
	}
}

func (r *repo) Delete(userID int64, serviceName string) error {
	_, err := r.DB.Exec("delete from urls where id=$1 and login=$2", userID, serviceName)
	return err
}
