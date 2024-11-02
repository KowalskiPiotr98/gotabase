package tests

import (
	log "github.com/sirupsen/logrus"
)

func PanicOnErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
