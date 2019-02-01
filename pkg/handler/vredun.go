package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Error ..
type Error struct {
	Code  int
	Error string
}

// Succe ..
type Succe struct {
	Code    int
	Message string
}

// Resp ..
type Resp struct {
	Error Error
	Succe Succe
}

func logAccess(r *http.Request) { log.Printf("Request at %q", r.URL.Path) }

func parseBody(b interface{}, r *http.Request) (interface{}, error) {
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(rb), &b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Handler ..
func Handler(path string, bodyParser interface{}, f func(interface{}) Resp) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		logAccess(r)

		b, err := parseBody(bodyParser, r)
		if err != nil {
			return
		}

		resp := f(b)
		if &resp.Error != nil {
			w.WriteHeader(resp.Error.Code)
			w.Write([]byte(resp.Error.Error))
			return
		}

		w.WriteHeader(resp.Succe.Code)
		w.Write([]byte(resp.Succe.Message))
		log.Println(resp.Succe.Message)
	})
}
