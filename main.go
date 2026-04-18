package main

import "gin-gonic-gorm/bootstrap"

func main() {
	// migrate -path database/migrations -database "mysql://root:@tcp(localhost:3306)/go_gin_gonic" up
	bootstrap.BoostrapApp()
}
