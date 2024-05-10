package main

import "crawlquery/api/router"

func main() {
	r := router.NewRouter()

	r.Run(":8080")
}
