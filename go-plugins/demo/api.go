package demo

import (
	"net/http"

	"github.com/nht1206/go-study/go-plugins/core"
	"goji.io/v3/pat"
)

const (
	DefaultDemoPath = "/demo"
)

type Request struct {
	Data string `json:"data"`
}

type Response struct {
	Result     string `json:"result,omitempty"`
	IsContinue bool   `json:"is_continue,omitempty"`
	Err        string `json:"err,omitempty"`
}

type Driver interface {
	Handle(Request) Response
}

type Handler struct {
	core.Handler
	driver   Driver
	demoPath string
}

func NewHandler(driver Driver, opts ...HandleOption) *Handler {
	path := DefaultDemoPath
	h := &Handler{
		demoPath: path,
		driver:   driver,
		Handler:  core.NewHandler(),
	}

	for _, o := range opts {
		o(h)
	}

	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.HandleFunc(pat.Post(h.demoPath), func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := core.DecodeRequest(w, r, &req); err != nil {
			return
		}
		res := h.driver.Handle(req)
		core.EncodeResponse(w, res, res.Err != "")
	})
}
