package main

import (
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/mat"
)

func calculateDifferences(data []float64) []float64 {
	diffs := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		diffs[i-1] = data[i] - data[i-1]
	}
	return diffs
}

func trainAndForecast(prices []float64, volumeDiffs []float64) (float64, float64, float64, float64) {
	priceDifferences := calculateDifferences(prices)
	volume := calculateDifferences(volumeDiffs)

	y := prices[1:]
	x1 := priceDifferences
	x2 := volume

	// współczynniki
	intercept, slopePrice, slopeVolume := linearRegression(y, x1, x2)

	lastPriceDiff := x1[len(x1)-1]
	lastVolumeDiff := x2[len(x2)-1]

	predictedPrice := intercept + slopePrice*lastPriceDiff + slopeVolume*lastVolumeDiff
	return predictedPrice, intercept, slopePrice, slopeVolume
}

func linearRegression(y, x1, x2 []float64) (float64, float64, float64) {

	X := mat.NewDense(len(y), 3, nil)
	Y := mat.NewVecDense(len(y), y)

	for i := 0; i < len(y); i++ {
		X.Set(i, 0, 1)
		X.Set(i, 1, x1[i])
		X.Set(i, 2, x2[i])
	}

	var coef mat.Dense
	qr := new(mat.QR)
	qr.Factorize(X)
	qr.SolveTo(&coef, false, Y)

	a := coef.At(0, 0)
	b1 := coef.At(1, 0)
	b2 := coef.At(2, 0)

	return a, b1, b2
}

func rollingWindowForecast(data BlockchainData, windowSize int) ([]float64, []time.Time) {
	var forecastedPrices []float64
	var timestamps []time.Time
	var sumIntercept, sumSlopePrice, sumSlopeVolume float64
	var allParams AllModelParameters

	numberOfModels := 0

	for i := windowSize; i < len(data.MarketPrice); i++ {

		windowPrices := make([]float64, windowSize)
		windowVolumes := make([]float64, windowSize)
		for j := 0; j < windowSize; j++ {
			windowPrices[j] = data.MarketPrice[i-windowSize+j].Y
			windowVolumes[j] = data.TotalBitcoins[i-windowSize+j].Y
		}

		predictedPrice, intercept, slopePrice, slopeVolume := trainAndForecast(windowPrices, windowVolumes)

		sumIntercept += intercept
		sumSlopePrice += slopePrice
		sumSlopeVolume += slopeVolume

		numberOfModels++

		forecastedPrices = append(forecastedPrices, predictedPrice)
		timestamps = append(timestamps, time.Unix(data.MarketPrice[i].X/1000, 0))
		params := ModelParameters{
			Intercept:   intercept,
			SlopePrice:  slopePrice,
			SlopeVolume: slopeVolume,
			StartDate:   time.Unix(data.MarketPrice[i-windowSize].X/1000, 0),
			EndDate:     time.Unix(data.MarketPrice[i].X/1000, 0),
		}

		allParams.Parameters = append(allParams.Parameters, params)

	}

	saveModelParameters(allParams, "model_parameters.json")

	modelsCount := float64(len(data.MarketPrice) - windowSize)
	avgSlopePrice := sumSlopePrice / modelsCount
	avgSlopeVolume := sumSlopeVolume / modelsCount

	fmt.Printf(" SlopePrice: %f, SlopeVolume: %f\n", avgSlopePrice, avgSlopeVolume)

	return forecastedPrices, timestamps
}

func calculateMSE(actual, predicted []float64) float64 {
	var sumError float64
	for i := range actual {
		predictionError := predicted[i] - actual[i]
		sumError += predictionError * predictionError
	}
	meanSquaredError := sumError / float64(len(actual))
	return meanSquaredError
}

func calculateMAPE(actual, predicted []float64) float64 {
	var sumError float64
	for i := range actual {
		if actual[i] != 0 {
			sumError += math.Abs((actual[i] - predicted[i]) / actual[i])
		}
	}
	meanAbsolutePercentageError := (sumError / float64(len(actual))) * 100
	return meanAbsolutePercentageError
}
