package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jcfullmer/terminalDashboard/config"
	"github.com/jcfullmer/terminalDashboard/internal/utilities"

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
		pterm.Error.Printf("Error making spinner: %v", err)
	}
	weather, err := utilities.GetWeather(conf)
	if err != nil {
		pterm.Error.Printf("Error getting weather: %v", err)
	}

	spinner.UpdateText("Getting system Info")
	serviceStatus, err := utilities.GetServices(conf)
	if err != nil {
		pterm.Error.Printf("Error getting serviceStatus: %v", err)
	}

	// get current time
	area, _ := pterm.DefaultArea.WithCenter().Start()
	clock, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString(time.Now().Format("03:04PM"))).Srender()
	area.Update(clock)

	paddedBox := pterm.DefaultBox.WithLeftPadding(2).WithRightPadding(2).WithTopPadding(0).WithBottomPadding(0)
	weatherTitle := pterm.LightRed("Weather")

	var Panels []pterm.Panel
	box1 := paddedBox.WithTitle(weatherTitle).Sprintf("%s: %v %s", conf.LocationName, weather.Current.Temperature2M, weather.CurrentUnits.Temperature2M)
	Panels = append(Panels, pterm.Panel{box1})
	box2 := paddedBox.WithTitle(pterm.Green("Service Status")).Sprint(serviceStatus)
	Panels = append(Panels, pterm.Panel{box2})
	if conf.Battery {
		batt, err := utilities.GetBatteryPercent(conf)
		if err != nil {
			fmt.Printf("Error getting battery information: %v", err)
		} else {
			box3 := paddedBox.WithTitle(pterm.LightBlue("Battery")).Sprint(batt)
			Panels = append(Panels, pterm.Panel{box3})
		}
	}

	pterm.DefaultPanel.WithPanels([][]pterm.Panel{Panels}).Render()
}
