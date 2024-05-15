package main

import "github.com/khedhrije/podcaster-search-api/internal/bootstrap"

// @title           podcaster-search-api
// @version         1.0.0
// @description     This is the documentation for the podcaster-search-api.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   khedhri.je@gmail.com
// @host      		localhost:8083
//
// @securityDefinitions.apikey Bearer-JWT
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a valid JWT token.
//
// @securityDefinitions.apikey Bearer-APIKey
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a valid API key.

func main() {
	bootstrap.InitBootstrap().Run()
}
