package cmd

import (
	"sort"
)

func salientHelper(indexes map[int]bool, zs []float64, minDelta float64, begin, end int) {
	if end-begin < 2 {
		return
	}
	maxDelta := 0.0
	first, last := begin, end
	if zs[begin] <= zs[end] {
		highest := begin
		for i := begin + 1; i <= end; i++ {
			if zs[i] > zs[highest] {
				highest = i
			} else if delta := zs[highest] - zs[i]; delta > maxDelta {
				first, last = highest, i
				maxDelta = delta
			}
		}
	}
	if zs[begin] >= zs[end] {
		lowest := begin
		for i := begin + 1; i <= end; i++ {
			if zs[i] < zs[lowest] {
				lowest = i
			} else if delta := zs[i] - zs[lowest]; delta > maxDelta {
				first, last = lowest, i
				maxDelta = delta
			}
		}
	}
	if maxDelta == 0 || maxDelta < minDelta {
		return
	}
	indexes[first] = true
	indexes[last] = true
	salientHelper(indexes, zs, minDelta, begin, first)
	salientHelper(indexes, zs, minDelta, first, last)
	salientHelper(indexes, zs, minDelta, last, end)
}

func salient(zs []float64, minDelta float64) []int {
	switch len(zs) {
	case 0:
		return nil
	case 1:
		return []int{0}
	case 2:
		return []int{0, 1}
	default:
		indexes := map[int]bool{
			0:           true,
			len(zs) - 1: true,
		}
		salientHelper(indexes, zs, minDelta, 0, len(zs)-1)
		result := make([]int, 0, len(indexes))
		for i := range indexes {
			result = append(result, i)
		}
		sort.Ints(result)
		return result
	}
}
