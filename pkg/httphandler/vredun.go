package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Error ..
type Error struct {
	Code  int
	Error string
}

// Success ..
type Success struct {
	Code    int
	Message string
}

// Res ..
type Res struct {
	E Error
	S Success
}

// Body ..
type Body struct {
	Body map[string]string
}

func logAccess(r *http.Request) { log.Printf("Request at %q", r.URL.Path) }

func validBody(rb []byte, ff []string) (*Body, *Error) {
	b := make(map[string]string)
	err := json.Unmarshal(rb, &b)
	if err != nil {
		log.Println("Could not read body.")
		log.Println(err)
		return nil, &Error{Code: 400, Error: "Could not read body."}
	}

	for _, f := range ff {
		if _, ok := b[f]; !ok {
			log.Printf("Missing parameter: %q in request body. /n", f)
			return nil, &Error{
				Code:  400,
				Error: fmt.Sprintf("Missing parameter: %q in request body.", f)}
		}
	}

	return &Body{Body: b}, nil
}

// Handler ..
func Handler(path string, matchFields []string, f func(map[string]string) Res) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		logAccess(r)

		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Oh-uh, something's not quite right."))
			return
		}

		b, e := validBody(rb, matchFields)
		if e != nil {
			w.WriteHeader(e.Code)
			w.Write([]byte(e.Error))
			return
		}

		resp := f(b.Body)
		if &resp.E != nil {
			log.Println(resp.E.Code)
			w.WriteHeader(resp.E.Code)
			w.Write([]byte(resp.E.Error))
			log.Println(resp.E.Error)
			return
		}

		w.WriteHeader(resp.S.Code)
		w.Write([]byte(resp.S.Message))
		log.Println(resp.S.Message)
	})
}
