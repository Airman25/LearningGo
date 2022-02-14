package main

import "fmt"

func faktorial(n int) int {
	f := 1
	for i := 1; i <= n; i++ {
		f *= i
	}
	return f
}

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
				fmt.Println(permutation)
				return
			}
		}
		jumps++
	}
}

func main() {
	var n, m int
	n = 9
	m = 430
	permutation := make([]int, n, n)
	for i := 0; i < n; i++ {
		permutation[i] = i + 1
	}
	fmt.Println(faktorial(n), m)
	fmt.Println(1, permutation)

	for i := 0; i < faktorial(n)-1; i++ { //
		fmt.Print(i+2, " ")
		nextPermutation(permutation)
	}
}
