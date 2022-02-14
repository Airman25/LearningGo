package main

import "fmt"

func faktorial(n int) int {
	f := 1
	for i := 1; i <= n; i++ {
		f *= i
	}
	return f
}

func numberright(permutation []int, place int) int {
	counter := 0
	for i := place + 1; i < len(permutation); i++ {
		if permutation[place] > permutation[i] {
			counter++
		}
	}
	return counter
}

func numberOfPermutation(permutation []int) int {
	placenumber := 0
	for i := 0; i < len(permutation); i++ {
		numbersright := numberright(permutation, i)
		placenumber += numbersright * faktorial(len(permutation)-i-1)
		//fmt.Println(numbersright, len(permutation)-i)
	}
	return placenumber + 1
}

func main() {
	var n, number int
	fmt.Scan(&n)
	permutation := make([]int, n, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&number)
		permutation[i] = number
	}
	fmt.Println(numberOfPermutation(permutation))
}
