package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/bcrypt"
)

var middleware = stats.New()

// HandleRoot / (GET)
func HandleRoot(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	obj := map[string]interface{}{} //My Dynamic JSON Object

	t := time.Now() //Current Request DateTime

	obj["ip_address"] = r.RemoteAddr               //Client remote address
	obj["date_hash"] = MD5(t.Format(time.RFC3339)) //MD5 HASH Current DateTime

	response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusOK, DATA: obj} //Populate Response Object

	respondWithJSON(w, r, http.StatusOK, response) //Send response

}

//HandleStats /info (POST)
func HandleStats(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	obj := map[string]interface{}{} //My Dynamic JSON Object

	t := time.Now()      //Current Request DateTime
	err := r.ParseForm() //Get Headers Form

	if err != nil {
		obj["error"] = err
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusBadRequest, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusBadRequest, response)
		return //Send response
	}

	v := r.Form //Local form

	var key = v.Get("key")

	if key != "" {
		auth, localKey := CheckAuth(key)
		if auth == 0 {
			obj["error"] = "Wrong Key!"
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusUnauthorized, DATA: obj} //Populate Response Object
			respondWithJSON(w, r, http.StatusUnauthorized, response)
			return
		} else if auth == 2 {
			obj["error"] = "Key out of Date!"
			obj["key"] = localKey.AUTHKEY
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusUnauthorized, DATA: obj} //Populate Response Object
			respondWithJSON(w, r, http.StatusUnauthorized, response)
			return
		} else if auth == 1 && localKey.TYPE == 1 {
			stats := middleware.Data()
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusOK, DATA: stats} //Populate Response Object
			respondWithJSON(w, r, http.StatusOK, response)                                                                                   //Send response
		} else {
			obj["error"] = "Not enough permissions!"
			obj["key"] = localKey.AUTHKEY
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusForbidden, DATA: obj} //Populate Response Object
			respondWithJSON(w, r, http.StatusForbidden, response)
			return
		}

	} else {
		obj["error"] = "Key format not valid!"
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusUnauthorized, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusUnauthorized, response)
		return
	}

}

//HandleRegister /register (POST)
func HandleRegister(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	obj := map[string]interface{}{} //My Dynamic JSON Object

	t := time.Now()      //Current Request DateTime
	err := r.ParseForm() //Get Headers Form

	if err != nil {
		obj["error"] = err
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusBadRequest, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusBadRequest, response)
		return //Send response
	}

	v := r.Form //Local form

	var localUser User
	localUser.NAME = v.Get("name")
	localUser.SURNAME = v.Get("surname")
	localUser.EMAIL = v.Get("email")
	localUser.PASSWORD = v.Get("password")

	if localUser.NAME != "" && localUser.SURNAME != "" && localUser.EMAIL != "" && localUser.PASSWORD != "" {

		hashedpsw, err3 := bcrypt.GenerateFromPassword([]byte(localUser.PASSWORD), bcrypt.MinCost)
		localUser.PASSWORD = string(hashedpsw)
		if err3 != nil {
			fmt.Println(err3)
		}
		status, inserted := InsertNewUser(localUser)
		if inserted == false {
			obj["error"] = "Error while registering account " + localUser.EMAIL + "!"
			obj["info"] = status
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusUnprocessableEntity, DATA: obj} //Populate Response Object
			respondWithJSON(w, r, http.StatusUnprocessableEntity, response)
			return
		} else {
			obj["info"] = "User " + localUser.EMAIL + " successfully registered!"
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusCreated, DATA: obj} //Populate Response Object
			respondWithJSON(w, r, http.StatusCreated, response)
			return
		}

	} else {
		obj["error"] = "User values format not valid!"
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusConflict, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusConflict, response)
		return
	}

}

//HandleStats /login (POST)
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	obj := map[string]interface{}{} //My Dynamic JSON Object

	t := time.Now()      //Current Request DateTime
	err := r.ParseForm() //Get Headers Form

	if err != nil {
		obj["error"] = err
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusBadRequest, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusBadRequest, response)
		return //Send response
	}

	v := r.Form //Local form

	EMAIL := v.Get("email")
	PASSWORD := v.Get("password")

	if EMAIL != "" && PASSWORD != "" {

		respo, logged, _ := CheckLoginUser(EMAIL, PASSWORD)
		if logged == false {
			respo["error"] = "Error while authorizing account " + EMAIL + "!"
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusUnauthorized, DATA: respo} //Populate Response Object
			respondWithJSON(w, r, http.StatusUnauthorized, response)
			return
		} else {
			respo["info"] = "User " + EMAIL + " successfully authorized!"
			response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusCreated, DATA: respo} //Populate Response Object
			respondWithJSON(w, r, http.StatusCreated, response)
			return
		}

	} else {
		obj["error"] = "User values format not valid!"
		response := JSONResponse{METHOD: r.Method, DATE: t.Format(time.RFC3339), PATH: r.RequestURI, STATUS: http.StatusConflict, DATA: obj} //Populate Response Object
		respondWithJSON(w, r, http.StatusConflict, response)
		return
	}

}

func main() {

	//Loading .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Getting vars from .env file
	port := os.Getenv("PORT")

	//Setting up MUX Router
	router := mux.NewRouter()

	//Routes
	router.HandleFunc("/", HandleRoot) // Home route
	router.HandleFunc("/info", HandleStats).Methods("POST")
	router.HandleFunc("/register", HandleRegister).Methods("POST")
	router.HandleFunc("/login", HandleLogin).Methods("POST")

	//Setting up middleware-focused library
	n := negroni.Classic()
	//Stats middleware
	n.Use(middleware)
	//Router
	n.UseHandler(router)

	log.Println("Running API Service on port " + port)
	//Serve
	n.Run(":" + port)

}
