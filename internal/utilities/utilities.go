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

type StoryInfo struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

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

func GetTop3HackerNewsStories() ([]string, error) {
	bestStoryIDs := "https://hacker-news.firebaseio.com/v0/beststories.json"
	storyURL := "https://hacker-news.firebaseio.com/v0/item/"
	var result []string
	res, err := http.Get(bestStoryIDs)
	if err != nil {
		return result, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return result, err
	}

	first3 := ids[:3]
	for _, story := range first3 {
		getStoryJSON := fmt.Sprintf("%s%d.json", storyURL, story)
		storyRes, err := http.Get(getStoryJSON)
		if err != nil {
			return result, err
		}
		body, err := io.ReadAll(storyRes.Body)
		if err != nil {
			return result, err
		}
		var storyInf StoryInfo
		if err := json.Unmarshal(body, &storyInf); err != nil {
			return result, err
		}
		titleLink := fmt.Sprintf("%s\n    %s", storyInf.Title, storyInf.URL)
		result = append(result, titleLink)

	}
	return result, nil
}
