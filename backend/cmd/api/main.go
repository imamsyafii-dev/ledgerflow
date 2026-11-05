package main

import (
	"github.com/imamsyafii-dev/ledger/backend/internal/http/router"
)

func main() {
	engine := router.New()

	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
