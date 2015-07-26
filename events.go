package main

import "encoding/json"

type ErrorCode int

const (
	NoError ErrorCode = 1
	ErrorParsing

	// login errors
	NickNotAvail
	AlreadyLoggedIn
	NotLoggedIn
)

type (
	Event struct {
		Action string          `json:"action"`
		Ok     bool            `json:"ok", omitempty`
		Data   json.RawMessage `json:"data"`
	}

	MessageData struct {
		Sender  string `json: sender`
		Message string `json: message`
	}

	MessageRequestData struct {
		Id      int    `json: id`
		Message string `json: message`
	}

	MessageResponseData struct {
		Id        int       `json: int`
		ErrorCode ErrorCode `json: error_code`
	}

	LoginRequestData struct {
		Nick string `json: nick`
	}

	GenericResponseData struct {
		Id        int       `json:"id"`
		ErrorCode ErrorCode `json:"error_code", omitempty`
		Message   string    `json:"message", omitempty`
	}
)
