package main

import (
	"flag"
	"log"
	"os"

	"./file"
	"./slack"
)

func main() {
	sync := flag.Bool("sync", false, "to send messages synchronously")

	flag.Parse()

	data, err := file.GetLinesFromFilename("./log.txt") // @todo make it come from os.Args

	if err != nil {
		log.Fatal(err)
	}

	slackAPI := slack.NewAPI("/T8NNCR01G/B8NU7G8PM/9qSAhfSTSOvOjlVGdZCi1F2n")

	if *sync {
		if err := slackAPI.SendDataSynchronously(data); err != nil {
			log.Fatal(err)
		}
	} else {
		errors := slackAPI.SendDataConcurrently(data)

		if len(errors) > 0 {
			for _, err := range errors {
				log.Print(err)
			}
			os.Exit(1)
		}
	}
}
