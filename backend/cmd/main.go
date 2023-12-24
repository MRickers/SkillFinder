package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"skillFinder/internal/controllers"
	"skillFinder/internal/model"
	"skillFinder/internal/persist"
)

type Config struct {
	Endpoint string `json:"endPoint"`
	LogFile  string `json:"logFile"`
}

func readConfig(configFile string) (config Config, err error) {
	file, err := os.ReadFile(configFile)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	return
}

func main() {
	const version = "v.0.1.0"
	fmt.Println("Skill finder", version)

	log.SetFlags(log.Lshortfile)

	var configFile string
	flag.StringVar(&configFile, "config", "skillFinder.cfg", "Configuration file for skill finder")
	flag.Parse()

	config, err := readConfig(configFile)
	if err != nil {
		fmt.Printf("%s, setting default config\n", err)
		config = Config{
			Endpoint: ":1234",
			LogFile:  "",
		}
	}
	db := persist.NewDb()
	db.Users = append(db.Users, model.User{Name: "rickers_m", Password: "admin"})
	db.Persons["test_user"] = model.Person{
		Name:      "Test User",
		Team:      "Future Project",
		Abilities: []model.Ability{},
	}

	// First we get a "copy" of the entry
	if employee, ok := db.Persons["test_user"]; ok {

		// Then we modify the copy
		employee.Abilities = append(employee.Abilities, model.Ability{
			Specific: model.Skill{
				Keyword:     "Safety Requirements",
				Shortname:   "Safety",
				Description: "Create Saftey requirements for system under development",
			},
			Level: 3,
		})

		// Then we reassign map entry
		db.Persons["test_user"] = employee
	}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.UserRegistration(w, r, &db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.UserLogin(w, r, &db)
	})

	http.HandleFunc("/api/v1/editProfile", controllers.LoggedIn(controllers.UserProfileCreate, &db))
	http.HandleFunc("/api/v1/editProfileName", controllers.LoggedIn(controllers.UserProfileChangeUsername, &db))
	http.HandleFunc("/api/v1/editProfileTeam", controllers.LoggedIn(controllers.UserProfileChangeTeam, &db))
	http.HandleFunc("/api/v1/editProfileAbility", controllers.LoggedIn(controllers.UserProfileUpdateAbility, &db))
	http.HandleFunc("/api/v1/editProfileDeleteAbility", controllers.LoggedIn(controllers.UserProfileDeleteAbility, &db))
	http.HandleFunc("/api/v1/getProfile", controllers.LoggedIn(controllers.GetUserProfile, &db))
	http.HandleFunc("/api/v1/getPersons", controllers.LoggedIn(controllers.GetAllEmployees, &db))

	fmt.Println("Listen for endpoint:", config.Endpoint)
	err = http.ListenAndServe(config.Endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}
}
