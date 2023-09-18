package controllers

import (
	"backend/database"
	"backend/loader"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
func containsElement(arr []string, target string) bool {
    for _, element := range arr {
        if element == target {
            return true
        }
    }
    return false
}
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	var user database.User
	json.NewDecoder(r.Body).Decode(&user)
	
	res:=AddUserToDB(user)
	
	
	w.WriteHeader(http.StatusOK)
	
	
	json.NewEncoder(w).Encode(&res)
}

func Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	fmt.Println("bongu",r.URL.Path,w)
	type Req struct{
		Email string
		Password string
	}
	fmt.Println(r.Host,r.URL.Scheme,r.URL.RequestURI(),"get all")
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	var status int
	var msg string
	var tokenString string
	var err error
	dbUser:=GetUserByEmail(request.Email)
	if dbUser.Email==""{
		status=1
		msg="User does not exist"
	}else{
		err=bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(request.Password))
		if err!=nil{
			status=2
			msg="Password Mismatch"
		}else{
			status=3
			msg="Correct User"
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":request.Email ,
				"exp": time.Now().Add(time.Hour*24*30).Unix(),
			})
			
			// Sign and get the complete encoded token as a string using the secret
			tokenString, err= token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			cookie := http.Cookie{
				Name:     "Authorization",
				Value:    tokenString,
				Expires:  time.Now().Add(time.Hour*24),
				HttpOnly: true,
				
			}
			
			http.SetCookie(w, &cookie)
			
			if err!=nil{
				msg="Error while jwt creation"
			}
		}
	}
	enc,err:=loader.Encrypt(tokenString,os.Getenv("ENCODE_KEY"))
	fmt.Println(w)
	res:=map[string]interface{}{"status":status,"msg":msg,"access":enc,"role":dbUser.Role}
	json.NewEncoder(w).Encode(&res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type req struct{
		email string
	}
	var user database.User
	json.NewDecoder(r.Body).Decode(&user)
	dbUser:=GetUserByEmail(user.Email)
	DeleteUserbyEmail(dbUser)
	DeleteGroupUserMappingByEmail(user.Email)
	w.WriteHeader(http.StatusOK)
	res:=map[string]interface{}{"status":1,"msg":"Deleted User"}
	json.NewEncoder(w).Encode(&res)
}

func GetAllRoles(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	var roles[] database.Role
	result := loader.DB.Find(&roles)
	if result.Error!=nil{
		panic("Error while fetching all users")
	}
	var res []string
	for _, obj := range roles {
        res = append(res, obj.Name)
    }
	json.NewEncoder(w).Encode(&res)
}


func Validate(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	w.WriteHeader(http.StatusOK)
	res:=map[string]interface{}{"message":"Logged in"}
	json.NewEncoder(w).Encode(&res)
}

func AddNewRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Name string
	}
	var role database.Role
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	role.ID = uuid.New()
	role.Name=request.Name
	roles:=GetAllRolesDB()
	if containsElement(roles,role.Name){
		res:=map[string]interface{}{"status":1,"msg":"Role already exists"}
		json.NewEncoder(w).Encode(&res)
	}else{
	AddRoleToDB(role)
	res:=map[string]interface{}{"status":2,"msg":"New Role Added"}
	json.NewEncoder(w).Encode(&res)
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	var users[] database.User
	result := loader.DB.Find(&users)
	if result.Error!=nil{
		panic("Error while fetching all users")
	}
	var res []string
	for _, obj := range users {
        res = append(res, obj.Email)
    }
	json.NewEncoder(w).Encode(&res)

}

func AddUsersFromCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	res:=GetAllRolesDB()
	file, _, err := r.FormFile("csvFile")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	count:=0
	len:=0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, "Error reading CSV: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var user database.User
		user.Email=record[0]
		user.Password=record[1]
		user.Role=record[2]
		if record[0]!="email"{
			
			if  containsElement(res,record[2]){
				res1:=AddUserToDB(user)
				if res1["status"]==2{
					len+=1
				}else{
					fmt.Println(record)
					count+=1
				}
	
			}else{
				fmt.Println(record)
				count+=1
			}
		}
			
	}
	response:=map[string]interface{}{"Error":count,"Success":len}
	json.NewEncoder(w).Encode(&response)
	
}



