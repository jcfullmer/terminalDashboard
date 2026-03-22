package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/distatus/battery"
	"github.com/jcfullmer/terminalDashboard/config"
)

type WeatherRes struct {
	CurrentUnits struct {
		Temperature2M string `json:"temperature_2m"`
		Precipitation string `json:"precipitation"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Interval      int     `json:"interval"`
		Temperature2M float64 `json:"temperature_2m"`
		Precipitation float64 `json:"precipitation"`
	} `json:"current"`
}

func GetWeather(conf *config.Config) (WeatherRes, error) {
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

func GetServices(conf *config.Config) (string, error) {
	var Status string
	for i, service := range conf.Services {
		cmd := exec.Command("systemctl", "is-active", service)
		err := cmd.Run()
		if err != nil {
			// Check the exit code for non-zero status
			if exitErr, ok := err.(*exec.ExitError); ok {
				Status += fmt.Sprintf("%s: %s", service, exitErr)
			} else {
				fmt.Printf("Error running command: %v\n", err)
			}
		} else {
			Status += fmt.Sprintf("%s is active", service)
		}
		if i+1 != len(conf.Services) {
			Status += "\n"
		}
	}
	return Status, nil
}

func GetBatteryPercent(conf *config.Config) (string, error) {
	batteries, err := battery.GetAll()
	if len(batteries) == 0 {
		if err != nil {
			return "", fmt.Errorf("access denied: %v", err)
		}
		return "", fmt.Errorf("no battery found")
	}

	bat := batteries[0]
	if bat.Full <= 0 {
		return "", fmt.Errorf("battery reports 0 capacity")
	}
	percentage := (float64(bat.Current) / float64(bat.Full)) * 100
	status := fmt.Sprintf("Charge: %.0f%%", percentage)
	if bat.ChargeRate > 0 {
		status += fmt.Sprintf("\nRate: %.0f mW", bat.ChargeRate)
	}
	return status, nil
}
