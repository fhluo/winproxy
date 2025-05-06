package main

import (
	"github.com/fhluo/winproxy/go/cmd"
	"log"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
