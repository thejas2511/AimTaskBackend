package loader

import (
	"backend/database"
	"fmt"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() {
	var err error
	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err!=nil{
		panic("Error connecting to DB")
	}
}
func SyncDb(){
	DB.AutoMigrate(&database.User{})
	DB.AutoMigrate(&database.Project{})
	DB.AutoMigrate(&database.Role{})
	DB.AutoMigrate(&database.Group{})
	DB.AutoMigrate(&database.GroupProjectMapping{})
	DB.AutoMigrate(&database.GroupUserMapping{})

	//Mock Data
	roles := []database.Role{
        {ID: uuid.New(), Name: "Manager"},
        {ID: uuid.New(), Name: "CEO"},
        {ID: uuid.New(), Name: "Designer"},
    }
	fmt.Println(roles)
	DB.Create(&roles)
}

func GetConn(db *gorm.DB){
	db=DB
}