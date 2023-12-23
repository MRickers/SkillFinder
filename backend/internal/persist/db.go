package persist

import (
	"skillFinder/internal/model"
)

type Db struct {
	// user -> person
	Persons map[string]model.Person
	Users   []model.User
	// sessionID -> user.Name
	Session map[string]string
}

func ContainsUser(d *Db, user model.User) bool {
	for _, dbUser := range d.Users {
		if user.Name == dbUser.Name {
			return true
		}
	}
	return false
}

func NewDb() (db Db) {
	db.Session = make(map[string]string)
	db.Persons = make(map[string]model.Person)
	return
}
