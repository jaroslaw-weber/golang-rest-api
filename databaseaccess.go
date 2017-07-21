package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//database connection settings
func getDatabaseConnection() *pg.DB {

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Addr:     "localhost:5432",
		Database: "postgres",
	})
	return db
}

//schema creation
//todo: move to migration
func createSchema() error {

	db := getDatabaseConnection()
	for _, model := range []interface{}{&Member{}, &BookCategory{}, &Book{}, &BookStatus{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
