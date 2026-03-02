package config

import (
	"fmt"
	"net/http"
	"os"
	"os/user"

)

type Config struct {
	latitude  string
	longitude string
	service   string
}

func CheckConfig() error {
	confPath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	conf, err := os.Stat(confPath)
	if os.IsNotExist(err) {
		fmt.Println("Config not detected! Starting setup...")
		err = makeConfig(confPath)
	}
}

func makeConfig(path string) error {
	var userZip string
	var userCountryCode string
	fmt.Println("> What is your 2 character ISO Country code? US for the United States of America.)
	fmt.Scanln(&userCountryCode)
	fmt.Println("> What is your current zipcode?")
	fmt.Scanln(&userZip)
	requestURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json&countryCode=$s", userZip, userCountryCode)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error getting coordinates for weather, please run terminalDashboard -config make")
	}
}
