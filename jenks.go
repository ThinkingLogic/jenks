package jenks

import (
	"math"
	"sort"
	"fmt"
	"strings"
	"strconv"
)

// Jenks natural breaks optimization (http://en.wikipedia.org/wiki/Jenks_natural_breaks_optimization)
// Based on the javascript implementation: https://gist.github.com/tmcw/4977508
// though that implementation has a bug - it has been fixed here.

// NaturalBreaks returns the best nClasses natural breaks in the data,
// using the Jenks natural breaks classification method (http://en.wikipedia.org/wiki/Jenks_natural_breaks_optimization).
// It tries to maximize the similarity of numbers in groups while maximizing the distance between the groups.
func NaturalBreaks(data []float64, nClasses int) []float64 {
	// sort data in numerical order, since this is expected by the matrices function
	data = sortData(data)

	// sanity check
	uniq := deduplicate(data)
	if nClasses >= len(uniq) {
		return uniq
	}

	// get our basic matrices (we only need lower class limits here)
	lowerClassLimits, _ := getMatrices(data, nClasses)

	// extract nClasses out of the computed matrices
	return breaks(data, lowerClassLimits, nClasses)
}

// AllNaturalBreaks finds all natural breaks in the data, for every set of breaks between 2 breaks and maxClasses.
// Uses the Jenks natural breaks classification method (http://en.wikipedia.org/wiki/Jenks_natural_breaks_optimization).
// It tries to maximize the similarity of numbers in groups while maximizing the distance between the groups.
func AllNaturalBreaks(data []float64, maxClasses int) [][]float64 {
	// sort data in numerical order, since this is expected by the matrices function
	data = sortData(data)

	// sanity check
	uniq := deduplicate(data)
	if maxClasses > len(uniq) {
		maxClasses = len(uniq)
	}

	// get our basic matrices (we only need lower class limits here)
	lowerClassLimits, _ := getMatrices(data, maxClasses)

	// extract nClasses out of the computed matrices
	allBreaks := [][]float64{}
	for i := 2; i <= maxClasses; i++ {
		nClasses := breaks(data, lowerClassLimits, i)
		if i == len(uniq) {
			nClasses = uniq
		}
		allBreaks = append(allBreaks, nClasses)
	}
	return allBreaks
}

// Round rounds the values of the given breaks as much as possible without changing the membership of each class.
// e.g. will attempt to round 111.11 to 111.1, then 111, then 110, then 100, then 0
// - ensuring that using the rounded break value doesn't change the membership of any class.
func Round(breaks []float64, data []float64) []float64 {
	data = sortData(data)
	rounded := make([]float64, len(breaks))
	for breakIdx := range breaks {
		// floor is the value that this break must remain above
		dataIdx := sort.SearchFloat64s(data, breaks[breakIdx])
		var floor float64
		if dataIdx == 0 { // make sure we can't go below breaks[i] - (breaks[i+1]-breaks[i])
			floor = data[0] - (breaks[breakIdx+1]-breaks[breakIdx])
		} else {
			floor = data[dataIdx-1]
		}
		rounded[breakIdx] = roundValue(breaks[breakIdx], floor)
	}
	return rounded
}

// roundValue works by replacing each digit (from right to left) with 0 until the value is no longer above the floor value.
func roundValue(initialValue float64, floor float64) float64 {
	b := []byte(strings.Trim(fmt.Sprintf("%f", initialValue), "0"))
	value := initialValue
	for i := len(b)-1; i >= 0; i-- {
		if b[i] != '.' {
			b[i] = '0'
			round, e := strconv.ParseFloat(string(b), 64)
			if e != nil || round <= floor {
				return value
			}
			value = round
		}
	}
	return value
}

// sortData checks to see if the data is sorted, returning it unchanged if so. Otherwise, it creates and sorts a copy.
func sortData(data []float64) []float64 {
	if !sort.Float64sAreSorted(data) {
		data2 := make([]float64, len(data))
		copy(data2, data)
		sort.Float64s(data2)
		data = data2
	}
	return data
}

// deduplicate returns a de-duplicated copy of the given slice, retaining the original order.
func deduplicate(data []float64) []float64 {
	keys := make(map[float64]bool)
	uniq := []float64{}
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniq = append(uniq, entry)
		}
	}
	return uniq
}

// getMatrices Computes the matrices required for Jenks breaks.
// These matrices can be used for any classing of data with 'classes <= n_classes'
func getMatrices(data []float64, nClasses int) ([][]int, [][]float64) {

	// in the original implementation, these matrices are referred to
	// as 'LC' and 'OP'
	//
	// * lowerClassLimits (LC): optimal lower class limits
	// * variance_combinations (OP): optimal variance combinations for all classes
	lowerClassLimits := make([][]int, len(data)+1)
	varianceCombinations := make([][]float64, len(data)+1)

	// the variance, as computed at each step in the calculation
	variance := 0.0

	// Initialize and fill each matrix with zeroes
	for i := 0; i < len(data)+1; i++ {
		lowerClassLimits[i] = make([]int, nClasses+1)
		varianceCombinations[i] = make([]float64, nClasses+1)
	}

	for i := 1; i < nClasses+1; i++ {
		lowerClassLimits[1][i] = 1
		varianceCombinations[1][i] = 0
		// in the original implementation, 'Infinity' is used but
		// math.MaxFloat64 will do.
		for j := 2; j < len(data)+1; j++ {
			varianceCombinations[j][i] = math.MaxFloat64
		}
	}

	for l := 2; l < len(data)+1; l++ {

		// sum was 'SZ' originally.
		// this is the sum of the values seen thus far when calculating variance.
		sum := 0.0
		// 'ZSQ' originally. the sum of squares of values seen thus far
		sumSquares := 0.0
		// 'WT' originally. 'w' is the number of data points considered so far.
		// it's used as the divisor in floating-point math, so using float rather than int
		w := 0.0

		for m := 1; m < l+1; m++ {

			// 'III' originally
			lowerClassLimit := l - m + 1
			currentIndex := lowerClassLimit - 1
			val := data[currentIndex]

			// here we're estimating variance for each potential classing
			// of the data, for each potential number of classes.
			w++

			// increase the current sum and sum-of-squares
			sum += val
			sumSquares += val * val

			// the variance at this point in the sequence is the difference
			// between the sum of squares and the total x 2, over the number
			// of samples.
			variance = sumSquares - (sum*sum)/w
			if currentIndex != 0 {
				for j := 2; j < nClasses+1; j++ {
					// if adding this element to an existing class
					// will increase its variance beyond the limit, break
					// the class at this point, setting the lower_class_limit
					// at this point.
					if varianceCombinations[l][j] >= (variance + varianceCombinations[currentIndex][j-1]) {
						lowerClassLimits[l][j] = lowerClassLimit
						varianceCombinations[l][j] = variance + varianceCombinations[currentIndex][j-1]
					}
				}
			}
		}

		lowerClassLimits[l][1] = 1;
		varianceCombinations[l][1] = variance;
	}
	// return the two matrices. for just providing breaks, only
	// 'lower_class_limits' is needed, but variances can be useful to
	// evaluate goodness of fit.
	return lowerClassLimits, varianceCombinations
}

// breaks is the second part of the jenks recipe:
// take the calculated matrices and derive an array of n breaks.
func breaks(data []float64, lowerClassLimits [][]int, nClasses int) []float64 {

	classBoundaries := make([]float64, nClasses)

	// the calculation of classes will never include the lower bound, so we need to explicitly set it
	// the upper bound is not included in the result - but it would be the maximum value in the data
	classBoundaries[0] = data[0]

	// the lowerClassLimits matrix is used as indexes into itself here:
	// the next value of `k` is obtained from .
	k := len(data) - 1
	for i := nClasses; i > 1; i -- {
		boundaryIndex := lowerClassLimits[k][i] - 1
		classBoundaries[i-1] = data[boundaryIndex]
		k = boundaryIndex
	}

	return classBoundaries
}
