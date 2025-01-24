package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// print colourful text
func PrintRainbowGlowLargeText(text string) {
	colors := []string{
		"\033[31m", // Red
		"\033[33m", // Yellow
		"\033[32m", // Green
		"\033[36m", // Cyan
		"\033[34m", // Blue
		"\033[35m", // Magenta
	}

	cmd := exec.Command("figlet", text)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running figlet:", err)
		os.Exit(1)
	}

	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Loop through the generated ASCII art and print each line with rainbow color
	for i, line := range strings.Split(string(output), "\n") {
		color := colors[i%len(colors)] // Cycle through colors for each line
		fmt.Printf("%s%s\n", color, line)
	}
    
	fmt.Printf("\033[0m")
}

// Date and time
func Date() string {
    loc, err := time.LoadLocation("Asia/Kolkata")
    if err != nil {
        return ("Error loading location:")
    }
    return time.Now().In(loc).Format("Monday, 02 January 2006, 03:04:05 PM IST")
}



type weatherData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		SurfacePressure        float64  `json:"surface_pressure"`
		PressureMSL        float64  `json:"pressure_msl"`
		Cloudcover        float64  `json:"cloud_cover"`
	} `json:"current"`
	Daily struct {
		Sunrise []int64 `json:"sunrise"`
		Sunset  []int64 `json:"sunset"`
		TemperatureMax  []float32 `json:"temperature_2m_max"`
		TemperatureMin  []float32 `json:"temperature_2m_min"`
		DaylightDuration  []float32 `json:"daylight_duration"`
		UvIndexMax  []float32 `json:"uv_index_max"`
		RainSum  []float32 `json:"rain_sum"`
	} `json:"daily"`
}

type aqiData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		PM10 float64 `json:"pm10"`
		PM25 float64 `json:"pm2_5"`
		CarbonMonoxide float64 `json:"carbon_monoxide"`
		NitrogenDioxide float64 `json:"nitrogen_dioxide"`
		SulphurDioxide float64 `json:"sulphur_dioxide"`
		Ozone float64 `json:"ozone"`
		AerosolOpticalDepth float64 `json:"aerosol_optical_depth"`
		Dust float64 `json:"dust"`
		UVIndex float64 `json:"uv_index"`
		UVIndexClearSky float64 `json:"uv_index_clear_sky"`
		Ammonia float64 `json:"ammonia"`
		AlderPollen float64 `json:"alder_pollen"`
		BirchPollen float64 `json:"birch_pollen"`
		GrassPollen float64 `json:"grass_pollen"`
		MugwortPollen float64 `json:"mugwort_pollen"`
		OlivePollen float64 `json:"olive_pollen"`
		RagweedPollen float64 `json:"ragweed_pollen"`
	} `json:"current"`
}

// Get weather data
func GetWeatherData(lat, lon float64) {
	start := time.Now()

	weatherChan := make(chan weatherData, 1)
	aqiChan := make(chan aqiData, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(){
		defer wg.Done()
		getWeather(lat, lon, weatherChan)
	}()

	go func(){
		defer wg.Done()
		getAQI(lat, lon, aqiChan)
	}()

	wg.Wait()

	weather := <-weatherChan
	aqi := <-aqiChan
	printCombinedTable(weather, aqi)
	fmt.Printf("Total time taken: %.3fms\n\n", time.Since(start).Seconds())
}

// Get weather data
func getWeather(lat, lon float64, weatherChan chan weatherData) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=cloud_cover,pressure_msl,surface_pressure&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset,daylight_duration,sunshine_duration,uv_index_max,rain_sum&timeformat=unixtime&forecast_days=1", lat, lon)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error fetching weather data: %v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response: %v", err)
		return
	}

	var weather weatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		fmt.Printf("error parsing JSON: %v", err)
		return
	}

	// printWeatherTable(weather)
	weatherChan <- weather
}

// Get AQI
func getAQI(lat, lng float64, aqiChan chan aqiData) {
	url := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%.2f&longitude=%.2f&current=pm10,pm2_5,carbon_monoxide,nitrogen_dioxide,sulphur_dioxide,ozone,aerosol_optical_depth,dust,uv_index,uv_index_clear_sky,ammonia,alder_pollen,birch_pollen,grass_pollen,mugwort_pollen,olive_pollen,ragweed_pollen&timeformat=unixtime&forecast_days=1&domains=cams_global", lat, lng)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error fetching air quality data: %v", err)
	}
	defer res.Body.Close()

	// Read and parse the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response: %v", err)
	}

	// Unmarshal the JSON response
	var aqi aqiData
	if err := json.Unmarshal(body, &aqi); err != nil {
		fmt.Printf("error parsing JSON: %v", err)
	}

	// printAQITable(aqi)
	aqiChan <- aqi
}

func printCombinedTable(weather weatherData, aqi aqiData) {
	headerColor := "\033[1;36m"
	paramColor := "\033[1;33m"
	valueColor := "\033[1;32m"
	lineColor := "\033[1;34m"
	reset := "\033[0m"

	fmt.Printf("\n%-35s Weather and Air Quality Data%s\n", headerColor, reset)
	fmt.Printf("%s+----------------------+------------------+----------------------+------------------+%s\n", lineColor, reset)
	fmt.Printf("| %sWeather Parameter%s    | %sValue%s            | %sAQI Parameter%s        | %sValue%s            |\n",
		paramColor, reset, paramColor, reset, paramColor, reset, paramColor, reset)
	fmt.Printf("%s+----------------------+------------------+----------------------+------------------+%s\n", lineColor, reset)

	rows := [][]string{
		{
			fmt.Sprintf("%-20s", "Max Temperature"), 
			fmt.Sprintf("%s%.1f°C%s", valueColor, weather.Daily.TemperatureMax[0], reset),
			fmt.Sprintf("%-20s", "PM10"), 
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.PM10, 150, 50), aqi.Current.PM10, reset),
		},
		{
			fmt.Sprintf("%-20s", "Min Temperature"),
			fmt.Sprintf("%s%.1f°C%s", valueColor, weather.Daily.TemperatureMin[0], reset),
			fmt.Sprintf("%-20s", "PM2.5"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.PM25, 75, 25), aqi.Current.PM25, reset),
		},
		{
			fmt.Sprintf("%-20s", "Sunrise"),
			fmt.Sprintf("%s%s%s", valueColor, time.Unix(weather.Daily.Sunrise[0], 0).Format("15:04"), reset),
			fmt.Sprintf("%-20s", "Carbon Monoxide"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.CarbonMonoxide, 10, 5), aqi.Current.CarbonMonoxide, reset),
		},
		{
			fmt.Sprintf("%-20s", "Sunset"),
			fmt.Sprintf("%s%s%s", valueColor, time.Unix(weather.Daily.Sunset[0], 0).Format("15:04"), reset),
			fmt.Sprintf("%-20s", "Nitrogen Dioxide"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.NitrogenDioxide, 200, 100), aqi.Current.NitrogenDioxide, reset),
		},
		{
			fmt.Sprintf("%-20s", "Daylight Duration"),
			fmt.Sprintf("%s%.1f hrs%s", valueColor, weather.Daily.DaylightDuration[0], reset),
			fmt.Sprintf("%-20s", "Sulphur Dioxide"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.SulphurDioxide, 200, 100), aqi.Current.SulphurDioxide, reset),
		},
		{
			fmt.Sprintf("%-20s", "UV Index Max"),
			fmt.Sprintf("%s%.1f%s", valueColor, weather.Daily.UvIndexMax[0], reset),
			fmt.Sprintf("%-20s", "Ozone"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.Ozone, 180, 100), aqi.Current.Ozone, reset),
		},
		{
			fmt.Sprintf("%-20s", "Rain Sum"),
			fmt.Sprintf("%s%.1f mm%s", valueColor, weather.Daily.RainSum[0], reset),
			fmt.Sprintf("%-20s", "UV Index"),
			fmt.Sprintf("%s%.1f%s", getColor(aqi.Current.UVIndex, 11, 6), aqi.Current.UVIndex, reset),
		},
		{
			fmt.Sprintf("%-20s", "Surface Pressure"),
			fmt.Sprintf("%s%.1f hPa%s", valueColor, weather.Current.SurfacePressure, reset),
			fmt.Sprintf("%-20s", "Dust"),
			fmt.Sprintf("%s%.1f µg/m³%s", getColor(aqi.Current.Dust, 150, 50), aqi.Current.Dust, reset),
		},
		{
			fmt.Sprintf("%-20s", "Pressure MSL"),
			fmt.Sprintf("%s%.1f hPa%s", valueColor, weather.Current.PressureMSL, reset),
			fmt.Sprintf("%-20s", ""),
			"",
		},
		{
			fmt.Sprintf("%-20s", "Cloud Cover"),
			fmt.Sprintf("%s%.1f%%%s", valueColor, weather.Current.Cloudcover, reset),
			fmt.Sprintf("%-20s", ""),
			"",
		},
	}

	for _, row := range rows {


		// adjust padding if it is null or empty
		if row[1] == "" {
			fmt.Printf("| %s | %-16s | %-20s | %-16s |\n", row[0], row[1], row[2], row[3])
		}

		if row[3] == "" {
			fmt.Printf("| %s | %-27s | %-20s | %-16s |\n", row[0], row[1], row[2], row[3])
		} else {
			fmt.Printf("| %s | %-27s | %s | %-27s |\n", row[0], row[1], row[2], row[3])
		}
	}
	fmt.Printf("%s+----------------------+------------------+----------------------+------------------+%s\n", lineColor, reset)
}

// Color the text based on the value
func getColor(value float64, severeThreshold, moderateThreshold float64) string {
	if value >= severeThreshold {
		return "\033[1;31m"
	} else if value >= moderateThreshold {
		return "\033[1;33m"
	}
	return "\033[1;32m"
}


// // Print weather data in a table
// func printWeatherTable(weather weatherData) {
// 	headerColor := "\033[1;36m"
// 	paramColor := "\033[1;33m"
// 	valueColor := "\033[1;32m"
// 	lineColor := "\033[1;34m"
// 	reset := "\033[0m"

// 	fmt.Printf("\n%s Today Weather Data:%s\n", headerColor, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// 	fmt.Printf("| %sParameter%s             | %sValue%s              |\n", paramColor, reset, paramColor, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// 	fmt.Printf("| Max Temperature       | %s%6.2f°C%s           |\n", valueColor, weather.Daily.TemperatureMax[0], reset)
// 	fmt.Printf("| Min Temperature       | %s%6.2f°C%s           |\n", valueColor, weather.Daily.TemperatureMin[0], reset)
// 	fmt.Printf("| Sunrise               | %s%8s%s           |\n", valueColor, time.Unix(weather.Daily.Sunrise[0], 0).Format("15:04"), reset)
// 	fmt.Printf("| Sunset                | %s%8s%s           |\n", valueColor, time.Unix(weather.Daily.Sunset[0], 0).Format("15:04"), reset)
// 	fmt.Printf("| Daylight Duration     | %s%6.2f hrs%s       |\n", valueColor, weather.Daily.DaylightDuration[0], reset)
// 	fmt.Printf("| UV Index Max          | %s%6.2f%s             |\n", valueColor, weather.Daily.UvIndexMax[0], reset)
// 	fmt.Printf("| Rain Sum              | %s%6.2f mm%s          |\n", valueColor, weather.Daily.RainSum[0], reset)
// 	fmt.Printf("| Surface Pressure      | %s%6.2f hPa%s         |\n", valueColor, weather.Current.SurfacePressure, reset)
// 	fmt.Printf("| Pressure MSL          | %s%6.2f hPa%s        |\n", valueColor, weather.Current.PressureMSL, reset)
// 	fmt.Printf("| Cloud Cover           | %s%6.2f%%%s            |\n", valueColor, weather.Current.Cloudcover, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// }
// // Print AQI data in a table
// func printAQITable(aqi aqiData) {
// 	headerColor := "\033[1;36m"
// 	paramColor := "\033[1;33m"
// 	lineColor := "\033[1;34m"
// 	reset := "\033[0m"

// 	// Function to determine color based on value
// 	getColor := func(value float64, severeThreshold, moderateThreshold float64) string {
// 		if value >= severeThreshold {
// 			return "\033[1;31m" // Red for severe
// 		} else if value >= moderateThreshold {
// 			return "\033[1;33m" // Orange for moderate
// 		}
// 		return "\033[1;32m" // Green for ok
// 	}

// 	fmt.Printf("\n%sAir Quality Data:%s\n", headerColor, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// 	fmt.Printf("| %sParameter%s             | %sValue%s              |\n", paramColor, reset, paramColor, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// 	fmt.Printf("| PM10                  | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.PM10, 150, 50), aqi.Current.PM10, reset)
// 	fmt.Printf("| PM2.5                 | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.PM25, 75, 25), aqi.Current.PM25, reset)
// 	fmt.Printf("| Carbon Monoxide       | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.CarbonMonoxide, 10, 5), aqi.Current.CarbonMonoxide, reset)
// 	fmt.Printf("| Nitrogen Dioxide      | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.NitrogenDioxide, 200, 100), aqi.Current.NitrogenDioxide, reset)
// 	fmt.Printf("| Sulphur Dioxide       | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.SulphurDioxide, 200, 100), aqi.Current.SulphurDioxide, reset)
// 	fmt.Printf("| Ozone                 | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.Ozone, 180, 100), aqi.Current.Ozone, reset)
// 	fmt.Printf("| UV Index              | %s%6.2f%s             |\n", getColor(aqi.Current.UVIndex, 11, 6), aqi.Current.UVIndex, reset)
// 	fmt.Printf("| Dust                  | %s%6.2f µg/m³%s       |\n", getColor(aqi.Current.Dust, 150, 50), aqi.Current.Dust, reset)
// 	fmt.Printf("%s+-----------------------+--------------------+%s\n", lineColor, reset)
// }
