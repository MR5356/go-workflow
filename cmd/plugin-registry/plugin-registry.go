package main

import (
	"github.com/MR5356/go-workflow/hub"
	"github.com/sirupsen/logrus"
)

func main() {
	server := hub.NewServer()
	if err := server.Run(80, "/tmp/plugin"); err != nil {
		logrus.Fatal(err)
	}
}
