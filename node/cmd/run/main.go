package main

import "crawlquery/node/router"

func main() {
	r := router.NewRouter()

	r.Run(":9090")
}
