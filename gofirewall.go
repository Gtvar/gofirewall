package main

import (
	"fmt"
	"encoding/json"
)

/*
 * Interface for firewall
 */
type Firewall interface {
	Check(Request) Response
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
 * Constructor
 */
func (r *Response) NewResponse(code int, reason string) *Response {
	r.Code = code
	r.Reason = reason

	return r
}

/*
 * Firewall main 
 * Get json data and check it for firewalls
 */
func main () {
	var jsonBlob = in()

	defer func() {
		if err := recover(); err != nil {
			var ex = fmt.Errorf("%v", err)
			var response Response
			response = *response.NewResponse(1, ex.Error())

			out(response)
		}
	}()

	response := checker(jsonBlob)

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
func checker(jsonBlob []byte) Response {
	var request Request
	err := json.Unmarshal(jsonBlob, &request)
	if (err != nil) {
		panic(err)
	}

	var firewalls []Firewall
	var response Response
	firewalls = get_firewalls()

	for _, firewall := range firewalls {
		if firewall.Support(request) {
			return firewall.Check(request)
		}
	}
	
	return response
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

func (up UserProject) Check(request Request) Response {
	var self UserProject

	rawBody,_ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		panic(err)
	}

	var response Response
	response = *response.NewResponse(1, "test")

	return response
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

func (up Email) Check(request Request) Response {
	var self Email

	rawBody,_ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		panic(err)
	}

	var response Response
	response = *response.NewResponse(2, "strange")

	return response
}

func (up Email) Support(request Request) bool {
	return request.Cmd == "Email"
}

