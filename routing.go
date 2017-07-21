//todo: refactor
//todo: move database access to databaseaccess.go

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

//setup server with REST api
func setRoutesAndStartServer() {

	//database set schema. use only once. should be in the migration part butfor testing purposes will leave it here
	//createSchema()

	//setup router
	r := mux.NewRouter()
	//implemented routes
	r.HandleFunc("/member/{id}", memberRoute).Methods("GET")
	r.HandleFunc("/member/{id}", deleteMemberRoute).Methods("DELETE")
	r.HandleFunc("/member/{id}", updateMemberRoute).Methods("PATCH")
	r.HandleFunc("/member", addMemberRoute).Methods("POST")

	//todo routes
	r.HandleFunc("/", indexRoute).Methods("GET")
	r.HandleFunc("/book/category/{id}", notImplementedRoute).Methods("GET", "PUT")
	r.HandleFunc("/book/{id}", notImplementedRoute).Methods("GET", "PUT")
	r.HandleFunc("/book/category", notImplementedRoute).Methods("POST")
	r.HandleFunc("/book", notImplementedRoute).Methods("POST")
	http.Handle("/", r)

	//start server
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

//route for not implemented api's
func notImplementedRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API not implemented yet!"))
}

//index route.
//todo add info about api
func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Library REST API example in go!"))
}

func writeResponse(w http.ResponseWriter, jsonResponse string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}

//GET for member
func memberRoute(w http.ResponseWriter, r *http.Request) {
	db := getDatabaseConnection()
	ps := mux.Vars(r)
	var foundMember Member
	id := ps["id"]
	err := db.Model(&foundMember).Where("member.id=?", id).Select()
	if err != nil {
		writeResponse(w, createResponseResultJSON(false, "no member with this id"))
	} else {
		responseJSON, err := createResponseMemberJSON(true, "", foundMember)
		if err != nil {
			panic(err) //todo
		} else {
			writeResponse(w, responseJSON)
		}
	}
}

//DELETE for member
func deleteMemberRoute(w http.ResponseWriter, r *http.Request) {

	db := getDatabaseConnection()
	ps := mux.Vars(r)
	var foundMember Member
	id := ps["id"]
	err := db.Model(&foundMember).Where("member.id=?", id).Select()
	if err != nil {
		writeResponse(w, createResponseResultJSON(false, "no member with this id"))
		log.Println(err)
	} else {
		err = db.Delete(&foundMember)
		if err != nil {
			writeResponse(w, createResponseResultJSON(false, "failed to delete member"))
			log.Println(err)
		} else {
			jsonResponse := createResponseResultJSON(true, "")
			writeResponse(w, jsonResponse)
		}
	}
}

//POST for member
func addMemberRoute(w http.ResponseWriter, r *http.Request) {

	db := getDatabaseConnection()
	memberName := r.FormValue("Name")
	memberEmail := r.FormValue("Email")
	newMember := Member{Name: memberName, Email: memberEmail}
	err := db.Insert(newMember)
	if err != nil {
		writeResponse(w, createResponseResultJSON(false, "failed to add member"))
		log.Println(err)
	}
	jsonResponse, err := createResponseMemberJSON(true, "", newMember)
	if err != nil {
		panic(err)
	}
	writeResponse(w, jsonResponse)
}

//PATCH form member
func updateMemberRoute(w http.ResponseWriter, r *http.Request) {

	db := getDatabaseConnection()
	ps := mux.Vars(r)
	id := ps["id"]
	memberName := r.FormValue("Name")
	memberEmail := r.FormValue("Email")
	var foundMember Member
	err := db.Model(&foundMember).Where("member.id=?", id).Select()
	if err != nil {
		writeResponse(w, createResponseResultJSON(false, "no member with this id"))
	}
	if memberName == "" {
		memberName = foundMember.Name
	}
	if memberEmail == "" {
		memberEmail = foundMember.Email
	}
	updatedMember := Member{ID: foundMember.ID, Name: memberName, Email: memberEmail}
	err = db.Update(updatedMember)
	if err != nil {
		writeResponse(w, createResponseResultJSON(false, "failed to update member"))
		log.Println(err)
	}
	jsonResponse, err := createResponseMemberJSON(true, "", updatedMember)
	if err != nil {
		panic(err)
	}
	writeResponse(w, jsonResponse)
}
