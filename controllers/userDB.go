package controllers

import (
	"backend/database"
	"backend/loader"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckUserByEmail(email string)bool{
	var user database.User
	loader.DB.Where("email = ?",email).First(&user)
	
	
	if user.Email==""{
		return false
	}
	return true
}

func GetUserByEmail(email string)database.User{
	var user database.User
	loader.DB.Where("email=?",email).First(&user)
	return user
}

func DeleteUserbyEmail(user database.User){
	err := loader.DB.Delete(&user)
	if err != nil {
		fmt.Println(err)
	}
}
func DeleteGroupUserMappingByEmail(email string){
	res:=loader.DB.Where("email=?",email).Delete(&database.GroupUserMapping{})
	if res!=nil{
		fmt.Println(res)
	}
}

func AddRoleToDB(user database.Role){
	err:=loader.DB.Create(&user)
	if err != nil {
		fmt.Println(err)
	}
}

func AddUserToDB(user database.User)map[string]interface{}{
	hash,err:=bcrypt.GenerateFromPassword([]byte(user.Password),10)
	exist:=CheckUserByEmail(user.Email)
	msg:="User added"
	status:=2
	if exist {
		status=1
		msg="User exists"
	}else{
		user.Password=string(hash)
		fmt.Println("hi")
		loader.DB.Create(&user)
	}
	if err!=nil{
		panic("Error in Hashing")
	}
	res:=map[string]interface{}{"message":msg,"status":status}
	
	return res
}

func GetAllRolesDB()[]string{
	var roles []database.Role
	loader.DB.Find(&roles)
	var res []string
	for _, obj := range roles {
        res = append(res, obj.Name)
    }
	return res
}

func GetAllUsersDB()[]string{
	var roles []database.User
	loader.DB.Find(&roles)
	var res []string
	for _, obj := range roles {
        res = append(res, obj.Email)
    }
	return res
}