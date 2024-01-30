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
	switch len(args) {
	case 0:
		rootDir, err = os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
	case 1:
		rootDir = os.Args[1]
	default:
		log.Fatalln("invalid number of arguments (exactly one positional argument is required: package path containing config.go source file)")
	}

	err = configs.Generate(rootDir)
	if err != nil {
		log.Fatalln(err)
	}
}
