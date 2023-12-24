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
	jsonErrorData, err := model.Serialize(model.ErrorDto{
		Message:   "profile empty",
		ErrorCode: -10,
	})

	if err != nil {
		log.Println("serialize errorDto failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonErrorData)
}

func UserProfileChangeUsername(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	person, err := requestEmployeeData(r, db)
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
	if employee, ok := db.Persons[userName]; ok {
		employee.Name = person.Name
		db.Persons[userName] = employee

		w.WriteHeader(http.StatusOK)
		return
	}
	response, err := model.Serialize(model.ErrorDto{
		Message:   "user has to be created first",
		ErrorCode: -11,
	})
	if err != nil {
		log.Println("serialize error dto failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("user not found please create user first")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func UserProfileChangeTeam(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	person, err := requestEmployeeData(r, db)
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
	if employee, ok := db.Persons[userName]; ok {
		employee.Team = person.Team
		db.Persons[userName] = employee

		w.WriteHeader(http.StatusOK)
		return
	}
	response, err := model.Serialize(model.ErrorDto{
		Message:   "user has to be created first",
		ErrorCode: -11,
	})
	if err != nil {
		log.Println("serialize error dto failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("user not found please create user first")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func UserProfileDeleteAbility(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	person, err := requestEmployeeData(r, db)
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
	if len(person.Abilities) != 1 {
		log.Println("invalid ability count:", len(person.Abilities))
		writeErrorDto(w, model.ErrorDto{
			Message:   "invalid ability count",
			ErrorCode: -12,
		})
		return
	}

	if employee, ok := db.Persons[userName]; ok {
		foundAbility, index := containesAbility(employee.Abilities, person.Abilities[0])
		if foundAbility {
			employee.Abilities = append(employee.Abilities[:index], employee.Abilities[index+1:]...)
			db.Persons[userName] = employee
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func UserProfileUpdateAbility(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	person, err := requestEmployeeData(r, db)
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
	if len(person.Abilities) != 1 {
		log.Println("invalid ability count:", len(person.Abilities))
		writeErrorDto(w, model.ErrorDto{
			Message:   "invalid ability count",
			ErrorCode: -12,
		})
		return
	}

	if employee, ok := db.Persons[userName]; ok {
		foundAbility, index := containesAbility(employee.Abilities, person.Abilities[0])
		if foundAbility {
			employee.Abilities[index] = person.Abilities[0]
		} else {
			employee.Abilities = append(employee.Abilities, person.Abilities[0])
			db.Persons[userName] = employee
		}

		w.WriteHeader(http.StatusOK)
		return
	}
	writeErrorDto(w, model.ErrorDto{
		Message:   "user has to be created first",
		ErrorCode: -11,
	})
}

func UserProfileCreate(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	person, err := requestEmployeeData(r, db)
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
	if _, ok := db.Persons[userName]; !ok {
		db.Persons[userName] = person
		w.WriteHeader(http.StatusOK)
		return
	}
	writeErrorDto(w, model.ErrorDto{
		Message: "already created",
	})
}

func containesAbility(abilities []model.Ability, ability model.Ability) (ok bool, index int) {
	for idx, abilityList := range abilities {
		if abilityList.Specific.Shortname == ability.Specific.Shortname {
			ok = true
			index = idx
			return
		}
	}
	return
}

func writeErrorDto(w http.ResponseWriter, dto model.ErrorDto) {
	response, err := model.Serialize(model.ErrorDto{
		Message: "already created",
	})
	if err != nil {
		log.Println("serialize error dto failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func requestEmployeeData(r *http.Request, db *persist.Db) (person model.Person, err error) {
	if r.Method != http.MethodPost {
		log.Println("expected POST method")
		return person, fmt.Errorf("expected POST method")
	}
	person, err = parsePerson(r)
	if err != nil {
		log.Println("could not parse username:", err)
		return person, fmt.Errorf("could not parse json: %s", err)
	}
	return
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
