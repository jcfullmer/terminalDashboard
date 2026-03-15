package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/jcfullmer/terminalDashboard/config"
)

func getWeather(conf *config.Config) (WeatherRes, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,precipitation&wind_speed_unit=mph&temperature_unit=fahrenheit&precipitation_unit=inch",
		conf.Latitude, conf.Longitude)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherRes{}, fmt.Errorf("error getting weather: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherRes{}, fmt.Errorf("error getting weather: %v", err)
	}
	defer resp.Body.Close()

	weather := WeatherRes{}
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return WeatherRes{}, fmt.Errorf("error getting weather: %v", err)
	}

	return weather, nil
}

func getServices(conf *config.Config) (string, error) {
	var Status string
	for _, service := range conf.Services {
		cmd := exec.Command("systemctl", "is-active", service)
		err := cmd.Run()
		if err != nil {
			// Check the exit code for non-zero status
			if exitErr, ok := err.(*exec.ExitError); ok {
				Status += fmt.Sprintf("%s: %s\n", service, exitErr)
			} else {
				fmt.Printf("Error running command: %v\n", err)
			}
		} else {
			Status += fmt.Sprintf("%s is active\n", service)
		}
	}
	return Status, nil
}
