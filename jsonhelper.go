package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

//for getting, updating, adding member
func createResponseMemberJSON(success bool, errorMessage string, member Member) (string, error) {

	successString := "true"
	if !success {
		successString = "false"
	}

	responseString := fmt.Sprintf(`{"success": %v`, successString)
	if errorMessage != "" {
		responseString += fmt.Sprintf(`, "error": "%v"`, errorMessage)
	}

	responseString += fmt.Sprintf(`, "member": 
		{
		"ID" : %v,
		"Name" : "%v",
		"Email" : "%v"
		}`, member.ID, member.Name, member.Email)

	responseString += "}"

	jsonParsed, err := gabs.ParseJSON([]byte(responseString))

	return jsonParsed.String(), err
}

//for deleting member
func createResponseResultJSON(success bool, errorMessage string) string {

	successString := "true"
	if !success {
		successString = "false"
	}

	responseString := fmt.Sprintf(`{"success": %v`, successString)
	if errorMessage != "" {
		responseString += fmt.Sprintf(`, "error": "%v"`, errorMessage)
	}

	responseString += "}"

	jsonParsed, err := gabs.ParseJSON([]byte(responseString))
	if err != nil {
		panic(err) //todo
	}
	return jsonParsed.String()
}
