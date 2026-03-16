package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jcfullmer/terminalDashboard/config"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	WeatherText := fmt.Sprintf("Getting Weather for %s", conf.LocationName)

	spinner, err := pterm.DefaultSpinner.Start(WeatherText)
	if err != nil {
		fmt.Printf("Error making spinner: %v", err)
	}
	weather, err := getWeather(conf)
	if err != nil {
		fmt.Printf("Error getting weather: %v", err)
	}

	// spinner.Success()
	spinner.UpdateText("Getting system Info")
	serviceStatus, err := getServices(conf)
	if err != nil {
		fmt.Printf("Error getting serviceStatus: %v", err)
	}

	// get current time
	area, _ := pterm.DefaultArea.WithCenter().Start()
	clock, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString(time.Now().Format("03:04PM"))).Srender()
	area.Update(clock)

	paddedBox := pterm.DefaultBox.WithLeftPadding(2).WithRightPadding(2).WithTopPadding(1) //.WithBottomPadding(0)
	title := pterm.LightRed("Weather")

	box1 := paddedBox.WithTitle(title).Sprint("Current Weather\n", weather.Current.Temperature2M, weather.CurrentUnits.Temperature2M)
	box2 := paddedBox.WithTitle(pterm.Green("Service Status")).Sprint("Service Status\n", serviceStatus)

	pterm.DefaultPanel.WithPanels([][]pterm.Panel{
		{{box1}, {box2}},
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
