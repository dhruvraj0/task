package main

import (
	"fmt"
	"sort"
)

func main() {
	a := []int{1, 4, 3, 9, 2}
	ans := minmax(a)
	fmt.Println(ans)
}

func minmax(arr []int) [][]int {
	var sorted []int

	for i := 0; i < len(arr); i++ {
		sorted = append(sorted, arr[i])
	}
	sort.Ints(sorted)

	var final [][]int
	ma := make(map[int]int)

	for i := 0; i < len(arr); i++ {
		if _, ok := ma[sorted[i]]; !ok {
			ma[sorted[i]] = i
		}
	}

	for i := 0; i < len(arr); i++ {
		index := ma[arr[i]]
		j := index - 1
		k := index + 1

		for k < len(sorted) {
			if index == 0 {
				res := []int{arr[i] - 1, arr[i], sorted[k]}
				final = append(final, res)
				k++
			} else if j >= 0 && k < len(sorted) {
				res := []int{sorted[j], arr[i], sorted[k]}
				final = append(final, res)
				j--
				k++
			} else {
				break
			}
		}
	}
	return final
}
