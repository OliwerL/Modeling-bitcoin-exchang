package main

import (
	"fmt"
	"log"
	"time"
)

type TotalBitcoins struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

type MarketPrice struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

type BlockchainData struct {
	Metric1       string          `json:"metric1"`
	Metric2       string          `json:"metric2"`
	TotalBitcoins []TotalBitcoins `json:"total-bitcoins"`
	MarketPrice   []MarketPrice   `json:"market-price"`
	Type          string          `json:"type"`
	Average       string          `json:"average"`
	Timespan      string          `json:"timespan"`
}

type GoodData struct {
	coursedata   []float64
	bitcoinsdata []float64
	timespan     []time.Time
}

type ModelParameters struct {
	Intercept   float64   `json:"intercept"`
	SlopePrice  float64   `json:"slopePrice"`
	SlopeVolume float64   `json:"slopeVolume"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
}

type AllModelParameters struct {
	Parameters []ModelParameters `json:"parameters"`
}

func main() {
	var windowSize int
	fmt.Print("Enter window size: ")
	_, err := fmt.Scanln(&windowSize)
	if err != nil {
		log.Fatal("Invalid input for window size")
	}

	blockchainData := loadData("total-bitcoins.json")

	startYear := 2011
	flitredblockchainData := filterData(blockchainData, startYear, 2024)

	forecastedPrices, timestamps := rollingWindowForecast(flitredblockchainData, windowSize)

	actualPrices := make([]float64, len(flitredblockchainData.MarketPrice)-windowSize)
	for i := range actualPrices {
		actualPrices[i] = flitredblockchainData.MarketPrice[i+windowSize].Y
	}

	// Wykres
	generateForecastChart(actualPrices, forecastedPrices, timestamps)
	meanSquaredError := calculateMSE(actualPrices, forecastedPrices)
	meanAbsolutePercentageError := calculateMAPE(actualPrices, forecastedPrices)
	fmt.Printf("meanSquaredError:  %f\n", meanSquaredError)
	fmt.Printf("meanAbsolutePercentageError:  %f\n", meanAbsolutePercentageError)

}
