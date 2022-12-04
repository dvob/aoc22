package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var langNames = map[string]string{
	"sh": "Shell",
}

var text = `# Advent of Code 2022

Advent of code solved in various languages:
`

func main() {
	langs := map[string][]string{}
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, day := range entries {
		if !strings.HasPrefix(day.Name(), "day") {
			continue
		}

		langsOfDay, err := os.ReadDir(day.Name())
		if err != nil {
			log.Fatal(err)
		}
		for _, lang := range langsOfDay {
			if strings.HasSuffix(lang.Name(), ".txt") {
				continue
			}
			langs[lang.Name()] = append(langs[lang.Name()], day.Name())
		}
	}

	langOrder := []string{}
	for lang := range langs {
		langOrder = append(langOrder, lang)
	}

	sort.Strings(langOrder)

	buf := bytes.NewBuffer(nil)
	for _, lang := range langOrder {
		langName, ok := langNames[lang]
		if !ok {
			langName = strings.Title(lang)
		}
		buf.WriteString("* " + langName + "\n")
		sort.Strings(langs[lang])
		for _, day := range langs[lang] {
			buf.WriteString("  * [" + day + "](" + day + "/" + lang + ")\n")
		}
	}

	fmt.Println(text)
	fmt.Println(buf.String())
}
