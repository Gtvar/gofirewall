package main

import (
	"fmt"
	"encoding/json"
)

type Firewall interface {
	Check(Request) Response
	Support(Request) bool
}

type Request struct {
	Cmd string `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

type Response struct {
	Code int `json:"code"`
	Reason string `json:"reason"`
	Request Request `json:"request"`
}


func main () {


	/* check structure request:
	 cmd - Command
	 body - Body of command
	*/

	// upload module for command
	// set body to module

	// module validate body
	// module check body and return Response structure

	// return Response structure

	var jsonBlob = []byte(`{"cmd":"UserProject","body":{"user_id":7,"project":"p6"}}`)

	var request Request

	err := json.Unmarshal(jsonBlob, &request)
	if (err != nil) {
		fmt.Println("erorr", err)
	}

	fmt.Println(request)

	var firewall UserProject

	support := firewall.Support(request)

	if support {
		response := firewall.Check(request)

		fmt.Println(response)

		responseString, _ := json.Marshal(response)
		fmt.Println(string(responseString))
	}


}

type UserProject struct {
	UserId int `json:"user_id"`
	Project string `json:"project"`
}

func (up UserProject) Check(request Request) Response {

	var self UserProject

	rawBody,_ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println(self)

	var result Response

	result.Code = self.UserId
	result.Reason = self.Project
	result.Request = request

	return result
}

func (up UserProject) Support(request Request) bool {
	return request.Cmd == "UserProject"
}

