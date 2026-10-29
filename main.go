package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ISBachvarov21/stackit-coding-challenge/models"
	"github.com/ISBachvarov21/stackit-coding-challenge/routes"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/julienschmidt/httprouter"
)

var db *pg.DB
var ownerUUID string

func initDatabase() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "stackit_coding_challenge",
	})

	err := createSchema(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createSchema(db *pg.DB) error {
	modelList := []interface{}{
		(*models.User)(nil),
		(*models.Event)(nil),
		(*models.Ticket)(nil),
	}

	for _, model := range modelList {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var err error
	ownerUUID = os.Getenv("OWNER_UUID")
	fmt.Println("OWNER_UUID:", ownerUUID)

	db, err = initDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := httprouter.New()
	routes.Register(router, db, ownerUUID)
	log.Fatal(http.ListenAndServe(":8080", router))
}
