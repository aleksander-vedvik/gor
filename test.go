package gor

import "context"

type User struct{}
type Server struct{}

func (s *Server) GetUser(ctx context.Context, userId int) (*User, error) {
	return nil, nil
}

func main() {
	srv := &Server{}
	r := NewRouter()

	r.Get("/user", srv.GetUser)
	r.Post("/user", NewGorHandler(srv.GetUser))
}
