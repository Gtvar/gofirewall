package main

import (
	"fmt"
	"encoding/json"
)

/*
 * Interface for firewall
 */
type Firewall interface {
	Check(Request, *Response)
	Support(Request) bool
}

/*
 * Request struct
 */
type Request struct {
	Cmd string `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

/*
 * Response
 */
type Response struct {
	Code int `json:"code"`
	Reason string `json:"reason"`
}

/*
 * Firewall main
 * Get json data and check it for firewalls
 */
func main () {
	var jsonBlob = in()

	var response Response
	defer func() {
		if err := recover(); err != nil {
			var ex = fmt.Errorf("%v", err)
			response.Reason = ex.Error()

			out(response)
		}
	}()

	checker(jsonBlob, &response)

	out(response)
}

/*
 * Read data from source
 */
func in() []byte {
	var jsonBlob = []byte(`{"cmd":"UserProject","body":{"user_id":7,"project":"p6"}}`)
//	var jsonBlob = []byte(`{"cmd":"Email","body":{"email":"test@test.com"}}`)

	return jsonBlob
}

/*
 * Print response to out
 */
func out(response Response) {
	responseString, _ := json.Marshal(response)
	fmt.Println(string(responseString))
}

/*
 * Manage firewalls
 */
func checker(jsonBlob []byte, response *Response) {
	var request Request
	err := json.Unmarshal(jsonBlob, &request)
	if (err != nil) {
		panic(err)
	}

	var firewalls []Firewall
	firewalls = get_firewalls()

	for _, firewall := range firewalls {
		if firewall.Support(request) {
			firewall.Check(request, response)

			return
		}
	}
}

/*
 * Get list of firewalls
 */
func get_firewalls() []Firewall {
	var firewalls = make([]Firewall, 2)

	var userProject UserProject
	var email Email
	firewalls[0] = userProject
	firewalls[1] = email

	return firewalls
}

/*
 * UserProject firewall
 */
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

/*
 * Email firewall
 */
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

