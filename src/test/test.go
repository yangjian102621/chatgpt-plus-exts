package main

import (
	"strings"
)

func main() {
	input := "lal603743923_A_Chinese_girl_walks_barefoot_on_the_beach_wearing_04223eb8-e11d-4c41-8b12-e15639729299.png"
	index := strings.LastIndex(input, "_")
	if index != -1 {
		hash := input[index+1 : len(input)-4]
		println(hash)
	}
}
