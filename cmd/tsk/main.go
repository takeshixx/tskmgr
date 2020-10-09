package main

import (
	"log"

	"github.com/takeshixx/tskmgr/internal/ui"
)

func main() {
	ui, err := ui.NewUI()
	if err != nil {
		log.Fatal(err)
	}
	if err := ui.RunUI(); err != nil {
		log.Fatal(err)
	}
}
