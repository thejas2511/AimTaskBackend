package controllers

import (
	"backend/database"
	"backend/loader"
	"fmt"

	"github.com/google/uuid"
)

func AddGroupToDB(group database.Group){
	res:=loader.DB.Create(&group)
	
	if res.Error!=nil{
		panic(res.Error)
	}
}

func DeleteGroupUserByID(id uuid.UUID ,user string){
	res:=loader.DB.Where("group_id=? AND email!=?",id,user).Delete(&database.GroupUserMapping{})
	if res.Error!=nil{
		panic(res.Error)
	}
}

func DeleteGroupProjectByID(id uuid.UUID ){
	res:=loader.DB.Where("group_id=?",id).Delete(&database.GroupProjectMapping{})
	if res.Error!=nil{
		panic(res.Error)
	}
}

func AddGroupProjectMappingToDB(group []database.GroupProjectMapping){
	err:=loader.DB.Create(&group)
	if err!=nil{
		fmt.Println(err)
	}
}
func AddGroupUserMappingToDB(group []database.GroupUserMapping){
	err:=loader.DB.Create(&group)
	if err!=nil{
		fmt.Println(err)
	}
}
func GetGroupsByEmailDB(email string)[]database.Group{
	var groups []database.GroupUserMapping
	loader.DB.Where("email=?",email).Find(&groups)
	var final []database.Group
	for _,obj:=range groups{
		var temp database.Group
		
		loader.DB.Where("id=?",obj.GroupId).Find(&temp)
		
		final=append(final,temp)

	}
	return final
}

func GetEmailByGroupDB(id uuid.UUID)[]string{
	var users []database.GroupUserMapping
	loader.DB.Where("group_id=?",id).Find(&users)
	var res []string

	for _,obj:=range users{
		
		res=append(res,obj.Email)
	}
	fmt.Println(res,45)
	return res
}

