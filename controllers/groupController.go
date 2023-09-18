package controllers

import (
	"backend/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func AddNewGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Name string
		Email string
	}
	var request Req
	var group database.Group
	json.NewDecoder(r.Body).Decode(&request)
	group.ID = uuid.New()
	group.Name=request.Name
	AddGroupToDB(group)
	var Map database.GroupUserMapping
	var groups []database.GroupUserMapping
	Map.ID=uuid.New()
	Map.GroupId=group.ID
	Map.Email=request.Email
	groups=append(groups,Map)
	AddGroupUserMappingToDB(groups)
	res1:=map[string]interface{}{"status":1,"msg":"New Group Added"}
	json.NewEncoder(w).Encode(&res1)
	
}

func AddNewGroupProjectMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		GroupId uuid.UUID
		ProjectId []uuid.UUID
	}
	var request Req
	var groups []database.GroupProjectMapping
	json.NewDecoder(r.Body).Decode(&request)
	DeleteGroupProjectByID(request.GroupId)
	for _,obj:=range request.ProjectId{
		var group database.GroupProjectMapping
		group.ID = uuid.New()
		group.GroupId=request.GroupId
		group.ProjectId=obj
		groups=append(groups,group)	
	}
	
	fmt.Println(groups)
	AddGroupProjectMappingToDB(groups)

	res1:=map[string]interface{}{"status":1,"msg":"New GroupProject Mapping  Added"}
	json.NewEncoder(w).Encode(&res1)
	
}

func AddNewGroupUserMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");

	type Req struct{
		GroupId uuid.UUID
		Email []string
		User string
	}
	
	var request Req
	var groups []database.GroupUserMapping
	json.NewDecoder(r.Body).Decode(&request)
	DeleteGroupUserByID(request.GroupId,request.User)
	for _,obj:=range request.Email{
		var group database.GroupUserMapping
		group.ID = uuid.New()
		group.GroupId=request.GroupId
		group.Email=obj
		groups=append(groups,group)	
	}
	
	fmt.Println(groups)
	AddGroupUserMappingToDB(groups)

	res1:=map[string]interface{}{"status":1,"msg":"New GroupUser Mapping  Added"}
	json.NewEncoder(w).Encode(&res1)
}

func GetGroupsByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Email string
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(r.Host,r.URL.Scheme,r.URL.RequestURI(),"get all")
	final:=GetGroupsByEmailDB(request.Email)
	groupMaps := []map[string]interface{}{}
    for _, group := range final {
        groupMap := map[string]interface{}{
            "Id":   group.ID,
            "Name": group.Name,
        }
        groupMaps = append(groupMaps, groupMap)
    }
	json.NewEncoder(w).Encode(&groupMaps)
	
}

func GetProjectsByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		ID uuid.UUID
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	final:=GetProjectsByGroupDB(request.ID)
	projectMaps := []map[string]interface{}{}
    for _, project := range final {
        projectMap := map[string]interface{}{
            "ID":   project.ID,
            "Name": project.Name,
			"Creator":project.Creator,
			"Role":project.Role,
        }
        projectMaps = append(projectMaps, projectMap)
    }
	json.NewEncoder(w).Encode(&projectMaps)
	
}

func GetAllEmailsByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		ID uuid.UUID
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	
	final:=GetEmailByGroupDB(request.ID)
	json.NewEncoder(w).Encode(&final)
	
}

