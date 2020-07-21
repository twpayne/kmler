package cmd

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSalient(t *testing.T) {
	for i, tc := range []struct {
		zs       []float64
		minDelta float64
		expected []int
	}{
		{
			zs:       nil,
			minDelta: 0,
			expected: nil,
		},
		{
			zs:       []float64{0.0},
			minDelta: 0,
			expected: []int{0},
		},
		{
			zs:       []float64{0.0, 0.0},
			minDelta: 0,
			expected: []int{0, 1},
		},
		{
			zs:       []float64{0.0, 1.0, 2.0},
			minDelta: 0,
			expected: []int{0, 2},
		},
		{
			zs:       []float64{2.0, 1.0, 0.0},
			minDelta: 0,
			expected: []int{0, 2},
		},
		{
			zs:       []float64{0.0, 1.0, 0.0},
			minDelta: 0,
			expected: []int{0, 1, 2},
		},
		{
			zs:       []float64{1.0, 0.0, 1.0},
			minDelta: 0,
			expected: []int{0, 1, 2},
		},
		{
			zs:       []float64{0.0, 1.0, 2.0, 3.0},
			minDelta: 0,
			expected: []int{0, 3},
		},
		{
			zs:       []float64{0.0, 1.0, 2.0, 3.0, 2.0, 1.0, 0.0},
			minDelta: 0,
			expected: []int{0, 3, 6},
		},
		{
			zs:       []float64{1.0, 2.0, 0.0, 3.0},
			minDelta: 0,
			expected: []int{0, 1, 2, 3},
		},
		{
			zs:       []float64{3.0, 0.0, 2.0, 1.0},
			minDelta: 0,
			expected: []int{0, 1, 2, 3},
		},
	} {
		if i != 5 {
			continue
		}
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := salient(tc.zs, tc.minDelta)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
