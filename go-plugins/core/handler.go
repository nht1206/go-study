package core

import (
	"fmt"
	"net"
	"net/http"

	"goji.io/v3"
	"goji.io/v3/pat"
)

// Handler is the base to create plugin handlers.
// It initializes connections and sockets to listen to.
type Handler struct {
	mux *goji.Mux
}

// NewHandler creates a new Handler with an http mux.
func NewHandler() Handler {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get(PingPath), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", DefaultContentType)
		fmt.Fprintf(w, "pong")
	})

	mux.HandleFunc(pat.Get(Prestop), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", DefaultContentType)
		w.WriteHeader(http.StatusOK)
	})

	return Handler{
		mux: mux,
	}
}

func (h Handler) Serve(l net.Listener) error {
	server := http.Server{
		Addr:    l.Addr().String(),
		Handler: h.mux,
	}

	return server.Serve(l)
}

func (h *Handler) HandleFunc(path goji.Pattern, fn func(http.ResponseWriter, *http.Request)) {
	h.mux.HandleFunc(path, fn)
}
