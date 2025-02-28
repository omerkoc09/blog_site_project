package main

import (
	"github.com/hayrat/go-template2/backend/idare/router"
	"github.com/hayrat/go-template2/backend/pkg/app"
)

var Version = "v0.0.0"
var BuildTime = "Bilinmiyor"

func main() {
	r := router.NewIdareRouter()
	a := app.New(r, Version, BuildTime)
	a.Start()
}
