package runner

import "os"

type Runner struct {
	BaseURL string
}

var instance *Runner

func GetRunner() *Runner {
	if instance == nil {
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "https://qa-internship.avito.com"
		}
		instance = &Runner{BaseURL: baseURL}
	}
	return instance
}
