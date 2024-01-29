package main

import (
	"flag"
	"github.com/ivanmashin/vanya/internal/configs"
	"log"
	"os"
)

func main() {
	var (
		rootDir string
		err     error
	)

	args := flag.Args()
	if len(args) != 1 {
		rootDir, err = os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
	}

	rootDir = os.Args[1]
	err = configs.Generate(rootDir)
	if err != nil {
		log.Fatalln(err)
	}
}
