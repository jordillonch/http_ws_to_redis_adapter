package main

import (
	"github.com/jordillonch/http_ws_to_redis_adapter/cmd/di"
)

func main() {
	adapterDi := di.Init()
	adapterDi.HttpWsServerServices.Start()
}
