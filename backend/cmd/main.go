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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.UserRegistration(w, r, &db)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.UserLogin(w, r, &db)
	})

	http.HandleFunc("/editProfile", controllers.LoggedIn(controllers.UpdateUserProfile, &db))
	http.HandleFunc("/getProfile", controllers.LoggedIn(controllers.GetUserProfile, &db))

	http.HandleFunc("/persons", controllers.LoggedIn(func(w http.ResponseWriter, r *http.Request, db *persist.Db) {
		log.Println("Accessing proteced area!")
	}, &db))

	fmt.Println("Listen for endpoint:", config.Endpoint)
	err = http.ListenAndServe(config.Endpoint, nil)
	if err != nil {
		fmt.Println(err)
	}
}
