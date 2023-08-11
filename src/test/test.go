package main

import (
	"fmt"
	"regexp"
)

func main() {
	input := "band. When the snow reflected, it was even more brilliant. --niji 5 - @lal603743923 (95%)"
	pattern := `\(\d+\%\)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(input)
	fmt.Println(len(match))
}
