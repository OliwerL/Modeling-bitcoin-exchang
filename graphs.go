package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

func loadData(filePath string) BlockchainData {
	file, _ := os.ReadFile(filePath)
	var data BlockchainData
	json.Unmarshal(file, &data)

	return data
}

func filterData(data BlockchainData, startYear, endYear int) BlockchainData {
	filtered := BlockchainData{}
	for _, entry := range data.MarketPrice {
		year := time.Unix(entry.X/1000, 0).Year()
		if year >= startYear && year <= endYear {
			filtered.MarketPrice = append(filtered.MarketPrice, entry)
		}
	}
	for _, entry := range data.TotalBitcoins {
		year := time.Unix(entry.X/1000, 0).Year()
		if year >= startYear && year <= endYear {
			filtered.TotalBitcoins = append(filtered.TotalBitcoins, entry)
		}
	}
	return filtered
}

func saveModelParameters(allParams AllModelParameters, filename string) error {
	file, _ := os.Create(filename)

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(allParams)

	return nil
}

func generateForecastChart(actual, forecasted []float64, timestamps []time.Time) {

	actualSeries := chart.TimeSeries{
		Name:    "Actual Market Price",
		XValues: timestamps,
		YValues: actual,
	}

	forecastedSeries := chart.TimeSeries{
		Name:    "Forecasted Market Price",
		XValues: timestamps,
		YValues: forecasted,
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorYellow,
			StrokeDashArray: []float64{5, 5},
		},
	}

	graph := chart.Chart{
		Series: []chart.Series{actualSeries, forecastedSeries},
		XAxis: chart.XAxis{
			Name:  "Time",
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:  "Price",
			Style: chart.StyleShow(),
		},
	}

	fileName := "forecast_chart.png"
	f, _ := os.Create(fileName)

	defer f.Close()

	graph.Render(chart.PNG, f)
	fmt.Printf("Wykres z prognozami został zapisany do pliku: %s\n", fileName)
}
