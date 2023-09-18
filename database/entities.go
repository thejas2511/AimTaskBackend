package database

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"primaryKey"; json: "email"`
	Password string `json:"password"`
	Role 	 string  `json:role`
	Role1     *Role    `gorm:"foreignKey:Role;references:Name" json:"variable"`
}

type Project struct{
	gorm.Model
	ID	uuid.UUID `gorm:"primarykey;type:uuid"`
	Name string  	`json:"proj_name"`
	Creator string     `json:"creator"`
	Role 		pq.StringArray			`gorm:"type:text[]";json:"role"`
	Group       string                   `json:"group"`
}

type Role struct{
	gorm.Model
	ID uuid.UUID `type:uuid"`
	Name string  `gorm:"primaryKey";json:"role_name"`
}

type Group struct{
	gorm.Model
	ID uuid.UUID `gorm:"primarykey;type:uuid"`
	Name string `json:"group_name"`
}

type GroupProjectMapping struct {
	gorm.Model
	ID uuid.UUID `gorm:"primarykey;type:uuid"`
	GroupId uuid.UUID `json:"group_id"`
	ProjectId uuid.UUID `json:"project_id"`
}

type GroupUserMapping struct {
	gorm.Model
	ID uuid.UUID `gorm:"primarykey;type:uuid"`
	GroupId uuid.UUID `json:"group_id"`
	Email string `json:"email"`
}
