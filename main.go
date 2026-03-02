package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const WeatherAPI = "https://api.open-meteo.com/v1/forecast?latitude=41.9573&longitude=-76.518&current=temperature_2m,precipitation&wind_speed_unit=mph&temperature_unit=fahrenheit&precipitation_unit=inch"

func main() {
	spinner, err := pterm.DefaultSpinner.Start("Gathering weather")
	if err != nil {
		fmt.Printf("Error making spinner: %v", err)
	}
	weather, err := getWeather()
	if err != nil {
		fmt.Printf("Error getting weather: %v", err)
	}
	time.Sleep(time.Second * 3)
	spinner.Success()
	spinner.Info()
	spinner.UpdateText("Getting system Info")

	// get current time
	area, _ := pterm.DefaultArea.WithCenter().Start()
	clock, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString(time.Now().Format("03:04PM"))).Srender()
	area.Update(clock)

	paddedBox := pterm.DefaultBox.WithLeftPadding(4).WithRightPadding(4).WithTopPadding(1).WithBottomPadding(1)
	title := pterm.LightRed("Weather")

	box1 := paddedBox.WithTitle(title).Sprint("Current Weather\n", weather.Current.Temperature2M, weather.CurrentUnits.Temperature2M)

	pterm.DefaultPanel.WithPanels([][]pterm.Panel{
		{{box1}},
	}).Render()
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

func getWeather() (WeatherRes, error) {
	resp, err := http.Get(WeatherAPI)
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
