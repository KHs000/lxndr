package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
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

func processRequestBody(r *http.Request, b interface{}) (map[string]string, *Error) {
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

	return b.(map[string]string), nil
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	writeResponse(w, http.StatusOK, "It works...")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)

	type request struct {
		Email string `json:"email"`
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"]
	isNew := identifier.ValidateNewUser(email)
	if !isNew {
		writeResponse(w, http.StatusConflict, "This email has already been registred.")
		return
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		message := fmt.Sprintf("User Created. ID: %v.", primitive.ObjectID.String(id))
		writeResponse(w, http.StatusCreated, message)
		return
	}

	message := "There might be an error, please retry your operation."
	writeResponse(w, http.StatusInternalServerError, message)
	return
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	logAccess(r)

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	log.Println(b["email"])
}

// ExportRoutes ..
func ExportRoutes() {
	http.HandleFunc("/", defaultRoute)
	http.HandleFunc("/newUser", createUser)
	http.HandleFunc("/test", testFunc)
}

// Handler ..
// func Handler(path string, matchFields []string, f func(map[string]string) Res) {
// 	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
// 		logAccess(r)

// 		rb, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Oh-uh, something's not quite right."))
// 			return
// 		}

// 		b, e := validBody(rb, matchFields)
// 		if e != nil {
// 			w.WriteHeader(e.Code)
// 			w.Write([]byte(e.Error))
// 			return
// 		}

// 		resp := f(b.Body)
// 		if resp.E.Code != 0 {
// 			w.WriteHeader(resp.E.Code)
// 			w.Write([]byte(resp.E.Error))
// 			return
// 		}

// 		w.WriteHeader(resp.S.Code)
// 		w.Write([]byte(resp.S.Message))
// 		log.Println(resp.S.Message)
// 	})
// }
