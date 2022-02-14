package main

import "fmt"

func nextPermutation(permutation []int) {
	var jumps int
	for i := len(permutation) - 1; i > -1; i-- {
		for i2 := len(permutation) - 1; i2 > i; i2-- {
			if permutation[i] < permutation[i2] {
				permutation[i2], permutation[i] = permutation[i], permutation[i2]
				for i4 := 0; i4 < jumps; i4++ {
					for i3 := len(permutation) - 1; i3 > i+1; i3-- {
						if permutation[i3] < permutation[i3-1] {
							permutation[i3], permutation[i3-1] = permutation[i3-1], permutation[i3]
						}
					}
				}
				return
			}
		}
		jumps++
	}
}

func main() {
	var n, m int
	fmt.Scan(&n, &m)
	permutation := make([]int, n, n)
	for i := 0; i < n; i++ {
		permutation[i] = i + 1
	}
	for i := 0; i < m-1; i++ {
		nextPermutation(permutation)
	}
	fmt.Print(permutation)
}
