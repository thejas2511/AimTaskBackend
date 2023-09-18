package main

import (
	"backend/loader"
	"backend/routes"

	"gorm.io/gorm"
)
func init() {
	loader.LoadEnv()
	loader.GetDB()
	loader.SyncDb()

}
var db *gorm.DB
func main() {
	loader.GetConn(db)
	routes.InitializeRouter()
}