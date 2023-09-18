package controllers

import (
	"backend/database"
	"backend/loader"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)
func containsAllElements(a, b []string) bool {
    aMap := make(map[string]bool)
    for _, element := range a {
        aMap[element] = true
    }
    for _, element := range b {
        if !aMap[element] {
            return false
        }
    }
    return true
}

func AddNewProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Name string
		Creator string
		Role []string
		Group string
	}
	var project database.Project
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	project.ID = uuid.New()
	project.Name=request.Name
	project.Creator=request.Creator
	project.Role=pq.StringArray([]string(request.Role))
	fmt.Printf("%T",project.Role)
	project.Group=request.Group
	AddProjectToDB(project)
	res:=map[string]interface{}{"status":1,"msg":"New Project Added"}
	json.NewEncoder(w).Encode(&res)
}

func GetAllProjectsCreated(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Email string
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request.Email)
	final:=GetAllProjectsByEmail(request.Email)
	var projectMaps []map[string]interface{}
    for _, project := range final {
        projectMap := map[string]interface{}{
            "ID":          project.ID,
            "Name":        project.Name,
            "Creator":     project.Creator,
			"Role":		   project.Role,
        }
        projectMaps = append(projectMaps, projectMap)
    }
	json.NewEncoder(w).Encode(&projectMaps)
}

func GetAllProjectsInvlovedbyRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Role string
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	// targetURL ,_:= url.Parse("http://localhost:5001")
	// r.Host = targetURL.Host
	fmt.Println(r.Host,r.URL.Scheme,r.URL.RequestURI(),"get all")
	final:=GetAllProjectsByRole(request.Role)
	var projectMaps []map[string]interface{}
    for _, project := range final {
        projectMap := map[string]interface{}{
            "ID":          project.ID,
            "Name":        project.Name,
            "Creator":     project.Creator,
			"Role":		   project.Role,
        }
        projectMaps = append(projectMaps, projectMap)
    }
	json.NewEncoder(w).Encode(&projectMaps)

}

func GetAllProjects(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	
	var projects[] database.Project
	result := loader.DB.Find(&projects)
	if result.Error!=nil{
		panic("Error while fetching all users")
	}
	var projectMaps []map[string]interface{}
    for _, project := range projects {
        projectMap := map[string]interface{}{
            "ID":          project.ID,
            "Name":        project.Name,
        }
        projectMaps = append(projectMaps, projectMap)
    }
	json.NewEncoder(w).Encode(&projectMaps)

}

func GetAllGroupProjects(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	var projects[] database.Project
	var res []uuid.UUID
	type Req struct{
		Email string
	}
	var request Req
	json.NewDecoder(r.Body).Decode(&request)
	groups:=GetGroupsByEmailDB(request.Email)
	
    for _, group := range groups {
        projects=GetProjectsByGroupDB(group.ID)
		for _,p :=range projects{
			res=append(res,p.ID)
		}

    }
	json.NewEncoder(w).Encode(&res)

}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	type Req struct{
		Id uuid.UUID
	}
	var request Req
	
	json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request.Id,10)
	DeleteProjectDB(request.Id)
}

func AddProjectsFromCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true");
	roles:=GetAllRolesDB()
	users:=GetAllUsersDB()
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
		fmt.Println(record)
		if err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, "Error reading CSV: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var project database.Project
		project.ID=uuid.New()
		project.Name=record[0]
		project.Creator=record[1]
		project.Role=pq.StringArray([]string(strings.Split(record[2], ";")))
		
		if record[0]!="name"{
			if containsElement(users,record[1]) &&containsAllElements(roles,project.Role){
				AddProjectToDB(project)
				len+=1
			}else{
				count+=1
			}	
		}
	}
	response:=map[string]interface{}{"Error":count,"Success":len}
	json.NewEncoder(w).Encode(&response)
	
}

