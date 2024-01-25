package main

import (
	"flag"
	"github.com/ivanmashin/vanya/internal/configs"
	"log"
	"os"
)

func main() {
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("invalid number of arguments (exactly one positional argument is required: package path containing config.go source file)")
	}

	rootDir := os.Args[1]
	err := configs.Generate(rootDir)
	if err != nil {
		log.Fatalln(err)
	}
}
