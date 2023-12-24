package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"skillFinder/internal/model"
	"skillFinder/internal/persist"

	"github.com/google/uuid"
)

func LoggedIn(fn func(w http.ResponseWriter, r *http.Request, db *persist.Db), db *persist.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")

		if err != nil || cookie.Value == "" {
			log.Println("no session found")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized access"))
			return
		}

		sessionID := cookie.Value

		if user, ok := db.Session[sessionID]; !ok {
			log.Println("session expired")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			log.Println("User request:", user)
		}
		fn(w, r, db)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, err := parseUser(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := validateUserInput(db, user); err != nil {
		log.Println("invalid credentials:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !persist.ContainsUser(db, user) {
		log.Println("user not registered")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if alreadyLoggedIn(db, user) {
		log.Println("already logged in")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
		return
	}
	sessionID := uuid.NewString()
	db.Session[sessionID] = user.Name

	cookie := http.Cookie{Name: "session", Value: sessionID}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))

}

func UserRegistration(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Println("error decoding request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validateUserInput(db, user); err != nil {
		log.Println("invalid user input:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if persist.ContainsUser(db, user) {
		log.Println("user already registered")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Users = append(db.Users, user)
	db.Persons[user.Name] = model.Person{}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registration successful"))
}

func alreadyLoggedIn(db *persist.Db, user model.User) bool {
	if _, ok := db.Session[user.Name]; ok {
		return true
	}
	return false
}

func parseUser(r *http.Request) (user model.User, err error) {
	err = json.NewDecoder(r.Body).Decode(&user)
	return
}

func validateUserInput(db *persist.Db, user model.User) error {
	if user.Name == "" || user.Password == "" {
		return fmt.Errorf("missing required fields")
	}
	if len(user.Name) < 5 {
		return fmt.Errorf("username must be at least 5 characters long")
	}

	return nil
}
