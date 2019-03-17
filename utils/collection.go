package utils

import (
	"log"
)

// Reverse - reverse array
func Reverse(reverseMe []byte) []byte {

	for i := 0; i < len(reverseMe)/2; i++ {
		helper := reverseMe[i]
		reverseMe[i] = reverseMe[len(reverseMe)-i-1]
		reverseMe[len(reverseMe)-i-1] = helper
	}

	log.Printf("Reverse res: %s\n", reverseMe)

	return reverseMe

}

// CheckDiff - return a difference between arr1 and arr2
func CheckDiff(arr1, arr2 []int) []int {

	var res []int
	for _, a := range arr1 {
		if index(arr2, a) == -1 {
			res = append(res, a)
		}
	}

	log.Printf("Difference res: %v\n", res)
	return res
}

func index(vs []int, t int) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
