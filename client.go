package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	server *Server

	Nick     string
	conn     *websocket.Conn
	loggedin bool

	writeQ chan *Event
}

func NewClient(s *Server, conn *websocket.Conn) *Client {
	return &Client{
		server:   s,
		Nick:     "",
		loggedin: false,
		conn:     conn,
		writeQ:   make(chan *Event),
	}
}

func (c *Client) Handle() {
	go c.readLoop()
	c.writeLoop()
}

func (c *Client) readLoop() {
	for {
		var msg Event
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		fmt.Printf("%v\n", msg)

		switch msg.Action {

		case "login":
			c.handleLogin(msg.Data)

		case "message":
			c.handleMessage(msg.Data)
		}
	}

}

func (c *Client) writeLoop() {
	for {
		select {
		case ev := <-c.writeQ:
			c.conn.WriteJSON(ev)
		}
	}
}

func (c *Client) sendEvent(action string, isok bool, data interface{}) {
	rmsg, _ := json.Marshal(data)
	ev := &Event{
		Action: action,
		Ok:     isok,
		Data:   json.RawMessage(rmsg),
	}
	c.writeQ <- ev

}

func (c *Client) handleLogin(msg json.RawMessage) {
	response := &GenericResponseData{
		ErrorCode: NoError,
		Message:   "",
	}

	var loginData LoginRequestData
	err := json.Unmarshal(msg, &loginData)
	if err != nil {
		response.ErrorCode = ErrorParsing
		response.Message = "Unable to parse message"
		c.sendEvent("error", false, response)
		return
	}

	if c.loggedin {
		response.ErrorCode = AlreadyLoggedIn
		response.Message = "Already Logged In"
		c.sendEvent("error", false, response)
		log.Println("Error parsing message event")
		return
	}
	if c.server.IsNickAvailable(loginData.Nick) {
		c.server.RegisterClient(loginData.Nick, c)
		c.Nick = loginData.Nick
		c.loggedin = true
		c.sendEvent("ack_login", true, response)
	} else {
		response.ErrorCode = NickNotAvail
		response.Message = "The Nick has already been taken"
		c.sendEvent("error", false, response)
		log.Println("Nick Collision")
	}

	return
}

func (c *Client) handleMessage(msg json.RawMessage) {
	response := &GenericResponseData{
		Id:        -1,
		ErrorCode: NoError,
		Message:   "",
	}

	var mData MessageRequestData
	err := json.Unmarshal(msg, &mData)
	if err != nil {
		response.ErrorCode = ErrorParsing
		response.Message = "Unable to parse message"
		c.sendEvent("error", false, response)
		log.Println("Error parsing message event")
		return
	}

	// if not logged in, inform the client
	// send the ID of the message
	if !c.loggedin {
		response.Id = mData.Id
		response.ErrorCode = NotLoggedIn
		response.Message = "You are not logged in"
		c.sendEvent("error", false, response)
		log.Println("User not logged in, sending messages")
	}

	fmt.Println(mData.Message)

	c.server.HandleClientMessage(c, mData.Message)

	c.sendEvent("ack_message", true, response)
}

func (c *Client) SendMessage(sender *Client, msg string) {
	msgEv := &MessageData{
		Sender:  sender.Nick,
		Message: msg,
	}

	c.sendEvent("message", true, msgEv)
}

func (c *Client) nextId() int {
	return 42
}
