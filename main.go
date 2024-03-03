package gor

import (
	"context"
	"fmt"
	"net/http"
)

type GorHandler[T any, U any] func(ctx context.Context, req T) (resp U, err error)
type gorHandler func(ctx context.Context, req any) (resp any, err error)

func NewGorHandler[T any, U any](handler func(context.Context, T) (U, error)) gorHandler {
	return func(ctx context.Context, req any) (resp any, err error) {
		resp, err = handler(ctx, req.(T))
		return resp.(U), err
	}
}

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got path\n")
	})

	mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling task with id=%v\n", id)
	})
	return &Router{
		mux: mux,
	}
}

func (r *Router) Get(path string, handler any) {
	h, ok := handler.(gorHandler)
	if !ok {
		panic("handler is bad")
	}
	hand := NewGorHandler(h)
	r.mux.HandleFunc(fmt.Sprintf("GET %s", path), func(w http.ResponseWriter, r *http.Request) {
		resp, err := hand(context.Background(), r.PathValue("id"))
		fmt.Fprintf(w, "handling task with id=%v\n", resp, err)
	})
}

func (r *Router) Post(path string, handler gorHandler) {
	r.mux.HandleFunc(fmt.Sprintf("POST %s", path), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "got path\n")
	})
}

func (r *Router) Serve(addr string) error {
	return http.ListenAndServe("localhost:8090", r.mux)
}
