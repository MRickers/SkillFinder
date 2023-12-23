package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"skillFinder/internal/persist"
)

func GetAllEmployees(w http.ResponseWriter, r *http.Request, db *persist.Db) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	list, err := serializePersons(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(list)
}

func serializePersons(db *persist.Db) ([]byte, error) {
	personData, err := json.Marshal(db.Persons)

	if err != nil {
		err = fmt.Errorf("serialize person list failed: %s", err)
		return []byte(""), err
	}
	return personData, nil

}
