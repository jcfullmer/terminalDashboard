package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	LocationName string   `json:"name"`
	Latitude     string   `json:"latitude"`
	Longitude    string   `json:"longitude"`
	Services     []string `json:"services"`
	Battery      bool     `json:"battery"`
}

type geocodingAPI struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

func GetConfig() (*Config, error) {
	confPath, err := os.UserConfigDir()
	if err != nil {
		return &Config{}, err
	}
	confPath += "/.terminalDashboardConfig"
	err = checkConfig(confPath)
	if err != nil {
		return &Config{}, err
	}
	file, err := os.ReadFile(confPath)
	if err != nil {
		return &Config{}, err
	}
	conf := &Config{}
	err = json.Unmarshal(file, conf)
	if err != nil {
		return &Config{}, err
	}
	return conf, nil
}

func checkConfig(confPath string) error {
	_, err := os.Stat(confPath)
	if os.IsNotExist(err) {
		fmt.Println("Config not detected! Starting setup...")
		err = makeConfig(confPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeConfig(path string) error {
	reader := bufio.NewReader(os.Stdin)
	var userZip string
	var userCountryCode string
	var confServices []string
	var userBattery string
	fmt.Println("> What is your 2 character ISO Country code? US for the United States of America.")
	fmt.Scanln(&userCountryCode)
	fmt.Println("> What is your current zipcode?")
	fmt.Scanln(&userZip)
	fmt.Println("> What services do you want to keep track of? Separate entries with a comma.")
	line, _ := reader.ReadString('\n')
	if line != "" {
		line = strings.TrimSpace(line)
		confServices = strings.Split(line, ",")
		for i := range confServices {
			confServices[i] = strings.TrimSpace(confServices[i])
		}
	}
	requestURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json&countryCode=%s", userZip, strings.ToUpper(userCountryCode))
	res, err := http.Get(requestURL)
	if err != nil || res.StatusCode != 200 {
		log.Println("error getting latitude and longitude")
		return err
	}
	geoCodingRes := geocodingAPI{}
	body, _ := io.ReadAll(res.Body)
	if err = json.Unmarshal(body, &geoCodingRes); err != nil {
		return err
	}
	fmt.Println("> Do you have a device with a battery? y/N")
	fmt.Scanln(&userBattery)
	userBattLower := strings.ToLower(userBattery)
	confBatt := userBattLower == "y" || userBattLower == "yes"

	conf := Config{
		LocationName: geoCodingRes.Results[0].Name,
		Latitude:     fmt.Sprintf("%v", geoCodingRes.Results[0].Latitude),
		Longitude:    fmt.Sprintf("%v", geoCodingRes.Results[0].Longitude),
		Services:     confServices,
		Battery:      confBatt,
	}
	confData, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	os.WriteFile(path, confData, 0o644)
	return nil
}
