package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "-> ", log.Ldate|log.Ltime)

func Log(any any) {
	logger.Println(any)
}
