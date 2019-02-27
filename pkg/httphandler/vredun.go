package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
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

func logAccess(r *http.Request) { log.Printf("Request at %q", r.URL.Path) }

func writeResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, fmt.Sprintf(`{"message": "%v"}`, message))
}

func processRequestBody(r *http.Request, b interface{}) (map[string]interface{}, *Error) {
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &Error{
			Code:  http.StatusInternalServerError,
			Error: "Oh-uh, something's not quite right."}
	}

	err = json.Unmarshal(rb, &b)
	if err != nil {
		return nil, &Error{
			Code:  http.StatusInternalServerError,
			Error: "Oh-uh, something's not quite right."}
	}

	return b.(map[string]interface{}), nil
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	writeResponse(w, http.StatusOK, "It works...")
}

// ExportRoutes ..
func ExportRoutes() {
	http.HandleFunc("/", defaultRoute)
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/editUser", editUser)
	http.HandleFunc("/deleteUser", deleteUser)
}
