package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

/*
 * Interface for firewall
 */
type Firewall interface {
	Check(Request) (Response, *FirewallError)
	Support(Request) bool
}

/*
 * Request struct
 */
type Request struct {
	Cmd  string          `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

/*
 * Response
 */
type Response struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

/**
 * Error struct
 */
type FirewallError struct {
	Error   error
	Message string
	Code    int
}

/*
 * Constructor
 */
func makeResponse(code int, reason string) *Response {
	return &Response{code, reason}
}

/*
 * Firewall main
 * Get json data and check it for firewalls
 */
func main() {

	ln, _ := net.Listen("tcp", ":8085")

	conn, _ := ln.Accept()

	defer func() {
		if err := recover(); err != nil {
			var ex = fmt.Errorf("%v", err)
			response := *makeResponse(1, ex.Error())

			out(conn, response)

			server_loop(conn)
		}
	}()

	server_loop(conn)
}

/**
 * Loop connection for read-write
 * @param  {[type]} conn net.Conn      [description]
 * @return {[type]}      [description]
 */
func server_loop(conn net.Conn) {
	for {
		jsonBlob, _ := bufio.NewReader(conn).ReadString('\n')

		response, err := checker([]byte(jsonBlob))

		if err != nil {
			response = *makeResponse(err.Code, err.Message)
		}

		out(conn, response)
	}
}

/**
 * Print response to out
 * @param  {[type]} conn     net.Conn      [description]
 * @param  {[type]} response Response      [description]
 * @return {[type]}          [description]
 */
func out(conn net.Conn, response Response) {
	responseString, _ := json.Marshal(response)
	conn.Write(responseString)
	conn.Write([]byte("\n"))
}

/**
 * Manage firewalls
 * @param  {[type]} jsonBlob []byte)       (Response, *FirewallError [description]
 * @return {[type]}          [description]
 */
func checker(jsonBlob []byte) (Response, *FirewallError) {
	var request Request
	err := json.Unmarshal(jsonBlob, &request)
	if err != nil {
		return Response{}, &FirewallError{err, "Error decode json", 101}
	}

	var firewalls []Firewall
	var response Response
	firewalls = get_firewalls()

	for _, firewall := range firewalls {
		if firewall.Support(request) {
			return firewall.Check(request)
		}
	}

	return response, nil
}

/**
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

/**
 * UserProject firewall
 * {"cmd":"UserProject","body":{"user_id":7,"project":"p6"}}
 */
type UserProject struct {
	UserId  int    `json:"user_id"`
	Project string `json:"project"`
}

/**
 * UserProject firewall check
 * @param  {[type]} up UserProject)  Check(request Request) (Response, *FirewallError [description]
 * @return {[type]}    [description]
 */
func (up UserProject) Check(request Request) (Response, *FirewallError) {
	var self UserProject

	rawBody, _ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		return Response{}, &FirewallError{err, "Error decode json", 102}
	}

	return *makeResponse(1, "test"), nil
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

/**
 * Email firewall check
 * @param  {[type]} up Email)        Check(request Request) (Response, *FirewallError [description]
 * @return {[type]}    [description]
 */
func (up Email) Check(request Request) (Response, *FirewallError) {
	var self Email

	rawBody, _ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &self)
	if err != nil {
		return Response{}, &FirewallError{err, "Error decode json", 102}
	}

	return *makeResponse(2, "strange"), nil
}

func (up Email) Support(request Request) bool {
	return request.Cmd == "Email"
}
