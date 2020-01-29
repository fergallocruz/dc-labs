package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	wordDict := make(map[string]int)

	for _, s := range words {
		//fmt.Println(s)
		wordDict[s]++
	}
	return wordDict
}

func main() {
	wc.Test(WordCount)
}
