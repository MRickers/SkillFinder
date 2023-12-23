package model

import (
	"errors"
)

type Skill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Person struct {
	Name      string `json:"name"`
	Team      string `json:"team"`
	Abilities Skill  `json:"skill"`
}

type Observer interface {
	Update(message string)
}

type Observerable interface {
	AddObserver(observer Observer) error
	RemoveObserver(observer Observer) error
	NotifyObservers(message string)
}

type PersonObservers struct {
	observers []Observer
}

func ContainsPerson(persons []Person, person Person) bool {
	for _, personList := range persons {
		if personList == person {
			return true
		}
	}
	return false
}

func isObserverAvailable(observers []Observer, observer Observer) (int, bool) {
	for idx, obs := range observers {
		if obs == observer {
			return idx, true
		}
	}
	return 0, false
}

func (ob *PersonObservers) AddObserver(observer Observer) error {
	if _, found := isObserverAvailable(ob.observers, observer); found {
		return errors.New("observer already added")
	}
	ob.observers = append(ob.observers, observer)
	return nil
}

func (ob *PersonObservers) RemoveObserver(observer Observer) error {
	if index, found := isObserverAvailable(ob.observers, observer); found {
		ob.observers = append(ob.observers[:index], ob.observers[index+1:]...)
		return nil
	}
	return errors.New("observer not found")
}

func (ob *PersonObservers) NotifyObservers(message string) {
	for _, obs := range ob.observers {
		obs.Update(message)
	}
}
