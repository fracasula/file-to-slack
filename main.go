package main

import "./slack"

func main() {
	slackAPI := slack.API{
		Endpoint: slack.Endpoint{"/T8NNCR01G/B8NU7G8PM/9qSAhfSTSOvOjlVGdZCi1F2n"},
	}

	slackAPI.SendMessage("Ciao!")
}
