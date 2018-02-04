package main

import (
	"github.com/jeremyroberts0/itunes-to-spotify/api"
)

func main() {
	router := api.GetRouter()
	router.Run(":8081")
}
