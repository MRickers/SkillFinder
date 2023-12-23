package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"skillFinder/internal/model"
	"skillFinder/internal/persist"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	userName, err := getUser(r, db)
	if err != nil {
		log.Println("could not request user:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if person, ok := db.Persons[userName]; ok {
		jsonData, err := json.Marshal(person)
		if err != nil {
			log.Println("could not serialize person:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(jsonData)
		return
	}
	errorDto := model.ErrorDto{
		Message:   "profile empty",
		ErrorCode: -10,
	}
	jsonErrorData, err := json.Marshal(errorDto)

	if err != nil {
		log.Println("serialize errorDto failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonErrorData)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	if r.Method != http.MethodPost {
		log.Println("expected POST method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := parsePerson(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userName, err := getUser(r, db)
	if err != nil {
		log.Println("could not request user:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.Persons[userName] = person
	w.WriteHeader(http.StatusOK)
}

func getUser(r *http.Request, db *persist.Db) (user string, err error) {
	cookie, err := r.Cookie("session")

	if err != nil {
		return
	}

	if cookie.Value == "" {
		return user, fmt.Errorf("empty session id")
	}

	sessionID := cookie.Value

	if user, ok := db.Session[sessionID]; ok {
		return user, nil
	}
	return
}

func parsePerson(r *http.Request) (person model.Person, err error) {
	if err = json.NewDecoder(r.Body).Decode(&person); err != nil {
		err = fmt.Errorf("deserialize person failed: %s", err)
	}

	return
}
