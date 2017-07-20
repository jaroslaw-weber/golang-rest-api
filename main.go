//todo json rest api

package main

import (
	// "fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	// "strings"
	"github.com/gorilla/mux"

	"fmt"

	"log"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type User struct {
	Id     int64
	Name   string
	Emails []string
}

type Message struct {
	Id      int64
	Name    string
	Content string
}

type Room struct {
	Id   int64
	Name string
}

func main() {

	initializeDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/", indexRoute).Methods("GET")
	r.HandleFunc("/room/{name}", showRoomRoute).Methods("GET")
	r.HandleFunc("/room/{name}/add", dummyRoute).Methods("POST")
	r.HandleFunc("/room/{name}/remove", dummyRoute).Methods("POST")
	r.HandleFunc("/user/{name}", userProfileRoute).Methods("GET")
	r.HandleFunc("/user/{name}/add", userProfileRoute).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{Addr: ":8080"}
	err := srv.ListenAndServe()
	log.Fatal(err)

	//enable shutting down server
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		srv.Close()
		os.Exit(1)
	}()
	//enable shutting down server

}

func dummyRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dummy route!"))
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Index page!"))
}

func showRoomRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("room "))
	parameters := mux.Vars(r)
	roomName := parameters["name"]
	w.Write([]byte(roomName))
}

func userProfileRoute(w http.ResponseWriter, r *http.Request) {

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Addr:     "localhost:5432",
		Database: "postgres",
	})

	var foundUser User
	err := db.Model(&foundUser).First()
	if err != nil {
		panic(err)
	}
	userAsString := fmt.Sprintf("%#v", foundUser)
	w.Write([]byte(userAsString))
}

func initializeDatabase() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Addr:     "localhost:5432",
		Database: "postgres",
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	err = loadDatabaseData(db)
	if err != nil {
		panic(err)
	}

}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&User{}, &Message{}, &Room{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func loadDatabaseData(db *pg.DB) error {
	//dummy database data load
	user1 := &User{
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	err := db.Insert(user1)
	if err != nil {
		return err
	}
	log.Print("database loaded!")

	//test if inserted!!
	var foundUser User
	err = db.Model(&foundUser).First()
	if err != nil {
		panic(err)
	}
	userAsString := fmt.Sprintf("%#v", foundUser)
	log.Print(userAsString)

	return nil
}
