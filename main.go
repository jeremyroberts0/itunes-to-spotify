package main

import (
	"fmt"
	"time"

	"github.com/jeremyroberts0/itunes-to-spotify/api"
)

func main() {
	fmt.Printf("Server starting at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	router := api.GetRouter()
	router.Run(":8081")
}
