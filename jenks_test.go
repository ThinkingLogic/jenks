package jenks_test

import (
	"reflect"
	"testing"
	"github.com/ThinkingLogic/jenks"
)

func TestNaturalBreaks(t *testing.T) {
	type args struct {
		data     []float64
		nClasses int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{name: "two breaks",
			args: args{nClasses: 2, data: []float64{1, 2, 3,   12, 13, 14,    21, 22, 23,    27, 28, 29}},
			want: []float64{1, 21}},
		{name: "three breaks",
			args: args{nClasses: 3, data: []float64{1, 2, 3,   12, 13, 14,    21, 22, 23,    27, 28, 29}},
			want: []float64{1, 12, 21}},
		{name: "four breaks",
			args: args{nClasses: 4, data: []float64{1, 2, 3,   12, 13, 14,    21, 22, 23,    27, 28, 29}},
			want: []float64{1, 12, 21, 27}},
		{name: "more breaks than unique values",
			args: args{nClasses: 4, data: []float64{1.1, 1.1, 1.1, 1.3, 1.3, 1.3, 1.2, 1.2, 1.2}},
			want: []float64{1.1, 1.2, 1.3}},
		{name: "http://www.real-statistics.com/multivariate-statistics/cluster-analysis/jenks-natural-breaks#example1",
			args: args{nClasses: 4, data: []float64{28.9, 33.5, 36.1, 38.6, 40.7, 42.7, 43.6, 45.8, 48.2, 48.6, 49.0, 51.0, 52.1, 52.2, 52.2, 52.4, 53.6, 54.2, 55.8, 55.8, 56.4, 56.8, 56.8, 57.7, 57.9, 58.2, 58.3, 58.4, 60.1, 60.1, 60.2, 61.1, 61.4, 61.9, 62.1, 62.5, 62.7, 63.1, 63.6, 64.2, 64.3, 64.4, 64.6, 64.7, 64.7, 64.8, 65.4, 65.8, 65.9, 66.2, 66.4, 66.6, 66.8, 67.0, 67.0, 67.1, 67.2, 67.2, 67.4, 68.2, 68.2, 68.3, 69.4, 69.5, 69.8, 70.2, 70.3, 70.5, 70.6, 71.2, 71.2, 71.2, 71.2, 71.8, 71.9, 72.0, 72.0, 72.0, 72.3, 72.5, 72.6, 73.0, 73.0, 73.0, 73.0, 73.2, 73.4, 73.4, 73.4, 74.0, 74.2, 74.4, 74.4, 74.9, 74.9, 75.4, 75.6, 76.0, 76.3, 76.3, 76.3, 76.4, 76.7, 77.2, 77.3, 77.6, 77.7, 78.3, 78.5, 78.5, 78.6, 78.7, 78.9, 79.2, 79.2, 79.2, 79.8, 79.8, 79.9, 80.7, 80.7, 81.2, 81.4, 81.5, 81.8, 82.0, 82.1, 82.2, 82.3, 82.4, 82.8, 83.0, 83.1, 83.3, 83.4, 83.6, 83.8, 83.8, 84.0, 84.2, 85.2, 85.4, 85.8, 86.1, 86.3, 87.1, 87.5, 87.7, 87.7, 87.8, 88.3, 88.9, 89.3, 90.3, 93.1, 94.2, 94.7, 95.7, 97.8, 99.2}},
			want: []float64{28.9, 55.8, 69.4, 80.7}},
		{name: "http://www.real-statistics.com/multivariate-statistics/cluster-analysis/jenks-natural-breaks#example2",
			args: args{nClasses: 3, data: []float64{5, 8, 9, 12, 15}},
			want: []float64{5, 8, 12}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jenks.NaturalBreaks(tt.args.data, tt.args.nClasses); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NaturalBreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllNaturalBreaks(t *testing.T) {
	type args struct {
		data       []float64
		maxClasses int
	}
	tests := []struct {
		name string
		args args
		want [][]float64
	}{
		{name: "two to four breaks",
			args: args{maxClasses: 4, data: []float64{1, 2, 3, 12, 13, 14, 21, 22, 23, 27, 28, 29}},
			want: [][]float64{
				{1, 21},
				{1, 12, 21},
				{1, 12, 21, 27},
			},
		},
		{name: "more breaks than unique values",
			args: args{maxClasses: 10, data: []float64{1, 2, 4, 1, 2, 4, 1, 2, 4}},
			want: [][]float64{
				{1, 4},
				{1, 2, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jenks.AllNaturalBreaks(tt.args.data, tt.args.maxClasses); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllNaturalBreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound(t *testing.T) {
	type args struct {
		data       []float64
		breaks     []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{name: "Round at small scale",
			args: args{breaks: []float64{1.1, 12.1, 21.1, 27.1}, data: []float64{1.1, 2.1, 3.1, 12.1, 13.1, 14.1, 21.1, 22.1, 23.1, 27.1, 28.1, 29.1}},
			want: []float64{0, 10, 20, 27},
		},
		{name: "Round at large scale",
			args: args{breaks: []float64{101, 1201, 2101, 2701}, data: []float64{101, 201, 1001, 1201, 1301, 1401, 2101, 2201, 2301, 2701, 2801, 2901}},
			want: []float64{0, 1200, 2000, 2700},
		},
		{name: "Don't round lowest bound below diff between lowest and next lowest",
			args: args{breaks: []float64{1.01, 2.01, 2.51}, data: []float64{1.01, 1.11, 1.12, 2.01, 2.11, 2.12, 2.51, 2.61, 2.72}},
			want: []float64{1, 2, 2.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jenks.Round(tt.args.breaks, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}
