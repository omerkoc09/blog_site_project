package main

import (
	"github.com/hayrat/go-template2/backend/app/router"
	"github.com/hayrat/go-template2/backend/pkg/app"
)

var Version = "v0.0.0"
var BuildTime = "Bilinmiyor"

func main() {

	r := router.NewAppRouter()
	a := app.New(r, Version, BuildTime)
	a.Static("/uploads", "./uploads")
	a.Start2()

}
