package utils

import "testing"

func TestMin_Int(t *testing.T) {
	tests := [][3]int{
		{0, 0, 0},
		{1, 1, 1},
		{-1, -1, -1},
		{-1, 1, -1},
		{1, -1, -1},
	}

	for _, test := range tests {
		result := Min(test[0], test[1])
		if result != test[2] {
			t.Errorf("Min(%d, %d) should be %d but is %d instead", test[0], test[1], test[2], result)
		}
	}
}

func TestMin_Float(t *testing.T) {
	tests := [][3]float64{
		{0, 0, 0},
		{1, 1, 1},
		{-1, -1, -1},
		{-1, 1, -1},
		{1, -1, -1},
		{0.5, 1.5, 0.5},
		{1.5, 0.5, 0.5},
		{0.5, -0.5, -0.5},
		{-0.5, 0.5, -0.5},
	}

	for _, test := range tests {
		result := Min(test[0], test[1])
		if result != test[2] {
			t.Errorf("Min(%f, %f) should be %f but is %f instead", test[0], test[1], test[2], result)
		}
	}
}
