package handler

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
)

// Response of Restful API
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

// msg of Response
const (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
	LOGOUT  = "LOGOUT"
)

// HomePage GET / -> /plat
func HomePage(c sola.Context) error {
	w := c.Response()
	w.Header().Add("Location", "/plat")
	w.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// Status ALL /status <Need Login>
func Status(c sola.Context) error {
	return c.JSON(http.StatusOK, Response{0, nil, SUCCESS})
}
