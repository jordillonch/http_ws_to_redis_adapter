package test

import "github.com/jordillonch/http_ws_to_redis_adapter/cmd/di"

var common *di.CommonDi

func setUp() {
	if common == nil {
		common = di.InitWithEnvFile("../.env", "../.env.test")
	}
}
