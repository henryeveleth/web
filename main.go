package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Index(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING INDEX")

	w.Header().Set("Content-Type", "application/json")
	users, err := GetAllUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func Show(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING SHOW")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	if id, err := strconv.Atoi(vars["id"]); err == nil {
		user, err := GetUser(id)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			js, err := json.Marshal(err)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(js)
			return
		} else {
			js, err := json.Marshal(user)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(js)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING CREATE")

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response := ResponseError{Message: "Bad arguments."}

		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(js)
		return
	}

	err = user.persist()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING UPDATE")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)

	if id, err := strconv.Atoi(vars["id"]); err == nil {
		user, err := GetUser(id)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			js, err := json.Marshal(err)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(js)
			return
		} else {
			err := decoder.Decode(&user)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				response := ResponseError{Message: "Bad arguments."}

				js, err := json.Marshal(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(js)
				return
			}

			err = user.persist()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING DELETE")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	if id, err := strconv.Atoi(vars["id"]); err == nil {
		user, err := GetUser(id)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)

			js, err := json.Marshal(err)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(js)
			return
		} else {
			user.DeletedAt = mysql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}

			err = user.persist()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/v1/users/{id}", Show).Methods("GET")
	router.HandleFunc("/v1/users/{id}", Update).Methods("PUT", "PATCH")
	router.HandleFunc("/v1/users/{id}", Delete).Methods("DELETE")

	router.HandleFunc("/v1/users/", Index).Methods("GET")
	router.HandleFunc("/v1/users/", Create).Methods("POST")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
