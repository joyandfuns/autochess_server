package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type weapon struct {
	Id          int
	Name        string
	Level       int
	Description string
	Url         string
	Parent      []int
	Child       []int
}

var db *sql.DB
var err error

const (
	pg_host     = "127.0.0.1"
	pg_port     = 5432
	pg_user     = "joyandfun"
	pg_password = "123456"
	pg_dbname   = "autochess"
)

func SelectAllWeapon() ([]weapon, error) {
	rows, err := db.Query("SELECT id, name, level, description, url, parent, child FROM weapon")
	if err != nil {
		return nil, err
	}

	var weapons []weapon
	for rows.Next() {
		var w weapon
		if err := rows.Scan(&w.Id, &w.Name, &w.Level, &w.Description, &w.Url, pq.Array(&w.Parent), pq.Array(&w.Child)); err != nil {
			return nil, err
		}
		weapons = append(weapons, w)
	}

	rows.Close()
	db.Close()
	return weapons, nil
}

func query() ([]weapon, error) {
	log.Println("connecting pg...")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_password, pg_dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("PG Successful Connected!")

	weapons, err := SelectAllWeapon()
	if err != nil {
		return nil, err
	}
	return weapons, nil
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		jsonData, err := query()
		if err != nil {
			log.Fatal("query data error: ", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": jsonData,
		})
	})
	r.Run(":9000")
}
