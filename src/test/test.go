package main

import (
	"fmt"
	"regexp"
)

func main() {
	input := "**A chinese girl, long hair and shawl. looking at view, At the age of 15 or 16, her skin was better than snow and was beautiful. The appearance was extremely beautiful, the whole body was dressed in red, and the hair was tied with a gold band. When the snow reflected, it was even more brilliant. --niji 5** - <@549957079511597067> (Waiting to start)"
	pattern := `\*\*(.*?)\*\*`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		fmt.Printf("%+v", matches[1])
	}
	//for _, match := range matches {
	//	if len(match) > 1 {
	//		fmt.Println("Found:", match[1])
	//	}
	//}
}
