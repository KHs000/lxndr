package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/pkg/mongodb"
)

// StartDatabase ..
func StartDatabase() {
	configF, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configF.Close()

	bytes, err := ioutil.ReadAll(configF)
	if err != nil {
		log.Fatal(err)
	}

	type configs struct {
		ConnStr string `json:"connectionString"`
	}

	var config configs
	json.Unmarshal(bytes, &config)
	mongodb.Connect(config.ConnStr)
	client := mongodb.NewClient(config.ConnStr)
	mongodb.Client = client.(domain.MongoClient)
}

func logAccess(r *http.Request) { log.Printf("Request at %q", r.URL.Path) }

func recovery(message string) {
	if r := recover(); r != nil {
		log.Println(message)
	}
}

func validateMethod(w http.ResponseWriter, r *http.Request, verb string) {
	if r.Method != verb {
		resp := domain.Response{Message: "Method not allowed."}
		writeResponse(w, http.StatusBadRequest, resp)
		panic("Method not allowed.")
	}
}

func processRequestBody(r *http.Request, b interface{}) (map[string]interface{},
	*domain.Error) {
	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &domain.Error{
			Code:  http.StatusInternalServerError,
			Error: "Oh-uh, something's not quite right."}
	}

	err = json.Unmarshal(rb, &b)
	if err != nil {
		return nil, &domain.Error{
			Code:  http.StatusInternalServerError,
			Error: "Oh-uh, something's not quite right."}
	}

	return b.(map[string]interface{}), nil
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(data)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Internal server error.")
		return
	}

	w.Write(res)
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	defer recovery("Method not allowed.")
	validateMethod(w, r, "GET")
	resp := domain.Response{Message: "It works..."}
	writeResponse(w, http.StatusOK, resp)
}

// ExportRoutes ..
func ExportRoutes() {
	http.HandleFunc("/", defaultRoute)
	http.HandleFunc("/createUser", createUserHandler)
	http.HandleFunc("/editUser", editUserHandler)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/listUsers", listUsers)
	http.HandleFunc("/test", testHandler)
}
