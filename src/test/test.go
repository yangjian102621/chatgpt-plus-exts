package main

import (
	"fmt"
	"regexp"
)

func main() {
	input := "**Ma painting of a young girl in green sitting at a pond, in the style of Liu ye, traditional animation, cinematic lighting, book sculptures --ar 16:9 --s 300 --v 5.2** - \u003c@1129939476223893514\u003e (93%) (fast)"
	pattern := `\((\d+)\%\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	fmt.Println(matches[1])
}
