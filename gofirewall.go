package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net"
)

var addrFlag string

const (
	successCode           = 0
	forbiddenCode         = 11
	errorCommonCode       = 21
	errorDecodejson       = 22
	errorMissingParameter = 23
	errorMissingFirewall  = 24
	errorDBConnect        = 25
	errorDBCommon         = 26

	errorTextJsonDecode = "Error decode json"

	mongoDsn = "localhost:27017"
	mondoDB  = "pg_firewall"
)

/**
 * Init flags
 */
func init() {
	const (
		defaultAddr = ":8085"
	)
	flag.StringVar(&addrFlag, "addr", defaultAddr, "Address for bind")
}

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

/*
 * Response constructor
 */
func makeResponse(code int, reason string) *Response {
	return &Response{code, reason}
}

/**
 * Error struct
 */
type FirewallError struct {
	Error   error
	Message string
	Code    int
}

func (fe FirewallError) GetMessage() string {
	if fe.Message != "" {
		return fe.Message
	}

	return fe.Error.Error()
}

/*
 * Firewall main
 * Get json data and check it for firewalls
 */
func main() {
	flag.Parse()

	ln, _ := net.Listen("tcp", addrFlag)

	conn, _ := ln.Accept()

	defer func() {
		if err := recover(); err != nil {
			var ex = fmt.Errorf("%v", err)
			response := *makeResponse(errorCommonCode, ex.Error())

			out(conn, response)
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
			response = *makeResponse(err.Code, err.GetMessage())
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
		return Response{}, &FirewallError{err, errorTextJsonDecode, errorDecodejson}
	}

	var firewalls []Firewall
	firewalls = get_firewalls()

	for _, firewall := range firewalls {
		if firewall.Support(request) {
			return firewall.Check(request)
		}
	}

	return *makeResponse(errorMissingFirewall, "Unknown Firewall"), nil
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
	Code    int
	Reason  string
	UserId  int    `json:"user_id"`
	Project string `json:"project"`
}

/**
 * UserProject firewall check
 * @param  {[type]} up UserProject)  Check(request Request) (Response, *FirewallError [description]
 * @return {[type]}    [description]
 */
func (up UserProject) Check(request Request) (Response, *FirewallError) {
	rawBody, _ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &up)
	if err != nil {
		return Response{}, &FirewallError{err, errorTextJsonDecode, errorDecodejson}
	}

	FirewallError := up.Load()
	if FirewallError.Error != nil {
		return Response{}, FirewallError
	}

	return *makeResponse(up.Code, up.Reason), nil
}

func (up UserProject) Support(request Request) bool {
	return request.Cmd == "UserProject"
}

func (up *UserProject) Load() *FirewallError {
	session, err := mgo.Dial(mongoDsn)
	if err != nil {
		return &FirewallError{err, "", errorDBConnect}
	}
	defer session.Close()

	c := session.DB(mondoDB).C("user_project")
	err = c.Find(bson.M{"user_id": up.UserId, "project": up.Project}).One(&up)
	if err != nil {
		return &FirewallError{err, "", errorDBCommon}
	}

	return &FirewallError{}
}

/*
 * Email firewall
 *
 * {"cmd":"Email","body":{"email":"test@test.com"}}
 */
type Email struct {
	Code   int
	Reason string
	Email  string `json:"email"`
}

/**
 * Email firewall check
 * @param  {[type]} up Email)        Check(request Request) (Response, *FirewallError [description]
 * @return {[type]}    [description]
 */
func (up Email) Check(request Request) (Response, *FirewallError) {
	rawBody, _ := request.Body.MarshalJSON()

	err := json.Unmarshal(rawBody, &up)
	if err != nil {
		return Response{}, &FirewallError{err, errorTextJsonDecode, errorDecodejson}
	}

	return *makeResponse(forbiddenCode, "strange"), nil
}

func (up Email) Support(request Request) bool {
	return request.Cmd == "Email"
}
