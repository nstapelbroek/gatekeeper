package handlers

import (
	"net/http"
)

func PostOpen(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello World!"))
}
