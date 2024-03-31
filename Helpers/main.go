package Helpers

import (
	"log"
)

func Log(err error, msg string) {
	if err != nil {
		log.Println(err, " -> Error at : ", msg)
		//log.Panicln(err, " -> Error at : ", msg)
	}
}
