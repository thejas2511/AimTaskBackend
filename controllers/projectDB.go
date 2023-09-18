package controllers

import (
	"backend/database"
	"backend/loader"
	"fmt"

	"github.com/google/uuid"
)

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func AddProjectToDB(proj database.Project){
	err:=loader.DB.Create(&proj)
	if err!=nil{
		fmt.Println(err)
	}
}
// func AddProjectMapToDB(proj database.ProjectMapping){
// 	err:=loader.DB.Create(&proj)
// 	if err!=nil{
// 		fmt.Println(err)
// 	}
// }
func GetAllProjectsByEmail(email string) []database.Project{
	var Projects []database.Project;
	loader.DB.Where("creator = ?",email).Find(&Projects)
	return Projects
}

func GetAllProjectsByRole(role string) []database.Project{
	var Projects []database.Project;
	loader.DB.Find(&Projects)
	return Projects
}

// func getAllProjectIds(email string) []database.ProjectMapping{
// 	var ProjectMappings []database.ProjectMapping
// 	loader.DB.Where("email=?",email).Find(&ProjectMappings)
	
// 	return ProjectMappings
// }
func getProjectById(id uuid.UUID) database.Project{
	var Project database.Project;
	loader.DB.Where("id = ?",id).Find(&Project)
	return Project
}

func GetProjectsByGroupDB(id uuid.UUID)[]database.Project{
	var projects []database.GroupProjectMapping;
	loader.DB.Where("group_id=?",id).Find(&projects)
	var final []database.Project
	for _,obj := range projects{
		var temp database.Project
		loader.DB.Where("id=?",obj.ProjectId).Find(&temp)
		final=append(final,temp)
	}
	return final
}

func DeleteProjectDB(id uuid.UUID){
	res:=loader.DB.Where("id=?",id).Delete(&database.Project{})
	res2:=loader.DB.Where("project_id=?",id).Delete(&database.GroupProjectMapping{})
	if res!=nil || res2!=nil{
		fmt.Println(res,res2)
	}
}