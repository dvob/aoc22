package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func getDayDirName(day int) string {
	return fmt.Sprintf("day%.02d", day)
}

func getInput(year, day int) ([]byte, error) {
	c := &http.Client{}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	session, found := os.LookupEnv("AOC_SESSION")
	if !found {
		return nil, fmt.Errorf("environment variable AOC_SESSION missing")
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: session,
	}
	req.AddCookie(cookie)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 399 {
		return nil, fmt.Errorf("http error: status=%d, body=%s", resp.StatusCode, input)
	}
	return input, nil
}

func run() error {
	var (
		year = 2022
	)
	flag.Parse()

	if flag.NArg() < 1 {
		return fmt.Errorf("missing argument day")
	}

	day, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return err
	}

	input, err := getInput(year, day)
	if err != nil {
		return err
	}

	dayDir := getDayDirName(day)
	err = os.MkdirAll(dayDir, 0750)
	if err != nil {
		return err
	}

	inputFile := filepath.Join(dayDir, "input.txt")
	f, err := os.Create(inputFile)
	if err != nil {
		return err
	}
	_, err = f.Write(input)
	if err != nil {
		return err
	}
	f.Close()

	sampleFile := filepath.Join(dayDir, "sample.txt")
	f, err = os.Create(sampleFile)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
