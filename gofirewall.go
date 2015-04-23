package main

import (
	"fmt"
	"encoding/json"
)

type Firewall interface {
	Check(Request, *Response)
	Support(Request) bool
}

type Request struct {
	Cmd string `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

type Response struct {
	Code int `json:"code"`
	Reason string `json:"reason"`
}



func main () {
	var jsonBlob = []byte(`{"cmd":"UserProject","body":{"user_id":7,"project":"p6"}}`)

	var response Response
	checker(jsonBlob, &response)

	responseString, _ := json.Marshal(response)
	fmt.Println(string(responseString))
}

func checker(jsonBlob []byte, response *Response) {
	var request Request
	err := json.Unmarshal(jsonBlob, &request)
	if (err != nil) {
		panic(err)
	}


	var firewall UserProject

	support := firewall.Support(request)

	if support {
		firewall.Check(request, response)
	}
}


func get_firewalls() []Firewall {
	var firewalls []Firewall

	var userProject UserProject
	var email Email
	firewalls[0] = userProject
	firewalls[1] = email

	return firewalls
}

type UserProject struct {
	UserId int `json:"user_id"`
	Project string `json:"project"`
}

func (up UserProject) Check(request Request, response *Response) {
	var self UserProject

	rawBody,_ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		panic(err)
	}

	response.Code = self.UserId
	response.Reason = self.Project
}

func (up UserProject) Support(request Request) bool {
	return request.Cmd == "UserProject"
}

type Email struct {
	Email string `json:"email"`
}

func (up Email) Check(request Request, response *Response) {
	var self Email

	rawBody,_ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		panic(err)
	}

	response.Code = 2
	response.Reason = "strange"
}

func (up Email) Support(request Request) bool {
	return request.Cmd == "Email"
}

