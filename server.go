package main

import "fmt"

type Server struct {
	nicks   map[string]bool
	clients map[*Client]bool
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
		nicks:   make(map[string]bool),
	}
}

func (s *Server) IsNickAvailable(nick string) bool {
	if _, ok := s.nicks[nick]; !ok {
		return true
	} else {
		return false
	}
}

func (s *Server) RegisterClient(nick string, cli *Client) error {
	if !s.IsNickAvailable(nick) {
		return fmt.Errorf("Nick not available")
	}

	s.nicks[nick] = true
	s.clients[cli] = true

	return nil
}

func (s *Server) UnRegisterClient(cli *Client) {
	delete(s.nicks, cli.Nick)
	delete(s.clients, cli)
}

func (s *Server) HandleClientMessage(cli *Client, msg string) {
	for client := range s.clients {
		if client == cli {
			continue
		}
		client.SendMessage(cli, msg)
	}
}
