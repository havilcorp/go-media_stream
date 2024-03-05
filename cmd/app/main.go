package main

import (
	"go-media-stream/internal/app"

	"github.com/sirupsen/logrus"
)

func main() {
	err := app.Run()
	if err != nil {
		logrus.Error(err)
	}
}
