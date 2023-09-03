package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct { // how json response mapped to weather struct
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TemperatureC float64 `json:"temp_c"`
		Condition    struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch    int64   `json:"time_epoch"`
				TemperatureC float64 `json:"temp_c"`
				Condition    struct {
					Text string `json:"text"`
				} `json:"condition"`
				Chanceofrain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Rostov-on-Don"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=268f512efda9425785d162316230309&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available!")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TemperatureC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TemperatureC,
			hour.Chanceofrain,
			hour.Condition.Text,
		)

		if hour.Chanceofrain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
