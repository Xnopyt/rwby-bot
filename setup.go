package main

import (
	"bufio"
	"os"
)

type loadedData struct {
	botToken    string
	FName       string
	LName       string
	Postcode    string
	CNum        string
	CMon        string
	CYea        string
	CCVV        string
	Anticaptcha string
}

func loadData() loadedData {
	tokens, err := os.Open("tokens.txt")
	if err != nil {
		panic("Couldn't open tokens.txt")
	}
	defer tokens.Close()
	scanner := bufio.NewScanner(tokens)
	i := 0
	var data loadedData
	for scanner.Scan() {
		switch i {
		case 0:
			data.botToken = scanner.Text()
		case 1:
			data.FName = scanner.Text()
		case 2:
			data.LName = scanner.Text()
		case 3:
			data.CNum = scanner.Text()
		case 4:
			data.CMon = scanner.Text()
		case 5:
			data.CYea = scanner.Text()
		case 6:
			data.CCVV = scanner.Text()
		case 7:
			data.Anticaptcha = scanner.Text()
		}
		i++
	}
	err = scanner.Err()
	if err != nil {
		panic("Failed to scan tokens.txt!")
	}
	return data
}
