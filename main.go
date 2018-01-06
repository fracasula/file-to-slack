package main

import (
	"bufio"
	"log"
	"os"

	"./slack"
)

func main() {
	file, err := os.Open("./log.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	slackAPI := slack.API{
		Endpoint: slack.Endpoint{"/T8NNCR01G/B8NU7G8PM/9qSAhfSTSOvOjlVGdZCi1F2n"},
	}

	for scanner.Scan() {
		err := slackAPI.SendMessage(scanner.Text())

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
