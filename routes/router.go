package routes

import (
	"backend/controllers"
	"backend/middleware"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func findBackendServer(path string) string {
	var backendServers = map[string]string{
		"/users":   "http://localhost:5001",  
		"/projects": "http://localhost:5002", 
		"/groups":  "http://localhost:5003",  
	}
	for prefix, backendURL := range backendServers {
		if strings.HasPrefix(path, prefix) {
			return backendURL
		}
	}
	return ""
}

func startServer(port string, name string) {
	r := mux.NewRouter()
	
	r.Handle("/validate", middleware.VerifyToken((http.HandlerFunc(controllers.Validate)))).Methods("POST")
	r.Handle("/projects/addProject", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewProject)))).Methods("POST")
	r.Handle("/users/addRole", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewRole)))).Methods("POST")
	r.Handle("/groups/addGroup", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewGroup)))).Methods("POST")
	r.Handle("/projects/getProjectsCreated", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllProjectsCreated)))).Methods("POST")
	r.HandleFunc("/users/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/users/login", controllers.Login).Methods("POST")
	r.HandleFunc("/users/getAllRoles",controllers.GetAllRoles).Methods("POST")
	r.Handle("/users/getAllUsers", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllUsers)))).Methods("POST")
	r.Handle("/projects/getAllProjects", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllProjects)))).Methods("POST")
	r.Handle("/projects/deleteProject", middleware.VerifyToken((http.HandlerFunc(controllers.DeleteProject)))).Methods("POST")
	r.HandleFunc("/users/deleteUser", controllers.DeleteUser).Methods("POST")
	r.Handle("/groups/addGroupProjectMapping", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewGroupProjectMapping)))).Methods("POST")
	r.Handle("/groups/addGroupUserMapping", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewGroupUserMapping)))).Methods("POST")
	r.Handle("/projects/getProjectsInvolvedByRole", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllProjectsInvlovedbyRole)))).Methods("POST")
	r.Handle("/projects/getAllGroupProjects", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllGroupProjects)))).Methods("POST")
	r.Handle("/groups/getEmailsByGroup", middleware.VerifyToken((http.HandlerFunc(controllers.GetAllEmailsByGroup)))).Methods("POST")
	r.Handle("/groups/getGroupsByEmail", middleware.VerifyToken((http.HandlerFunc(controllers.GetGroupsByEmail)))).Methods("POST")
	r.Handle("/groups/getProjectsByGroup", middleware.VerifyToken((http.HandlerFunc(controllers.GetProjectsByGroup)))).Methods("POST")
	r.HandleFunc("/users/addUsersByCSV", controllers.AddUsersFromCSV).Methods("POST")
	r.Handle("/projects/addProjectsByCSV", middleware.VerifyToken((http.HandlerFunc(controllers.AddProjectsFromCSV)))).Methods("POST")
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	ipPort := "0.0.0.0:" + port
	fmt.Print("Server running on " + ipPort + "\n")
	serverErr := http.ListenAndServe(ipPort, handlers.CORS(headers, methods, origins)(r))
	if serverErr != nil {
		panic(serverErr.Error())
	}
}


func startProxyServer(port string, name string) {	
	reverseProxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		targetURL,_:= url.Parse(findBackendServer(r.URL.Path))
		r.Host = targetURL.Host
        r.URL.Host = targetURL.Host
        r.URL.Scheme = targetURL.Scheme
        r.RequestURI = ""
		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            _, _ = fmt.Fprint(w, err)
            return
        }
		w.WriteHeader(http.StatusOK)
        io.Copy(w, originServerResponse.Body)
	})
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	corsHandler := corsOptions.Handler(reverseProxy)
	ipPort := "0.0.0.0:" + port
	fmt.Print("Proxy Server running on " + ipPort + "\n")
	serverErr := http.ListenAndServe(ipPort,corsHandler)
	if serverErr != nil {
		panic(serverErr.Error())
	}
}
func InitializeServers(){
	fmt.Println(56545665)
	port0:=os.Getenv("PORT")
	port1 := os.Getenv("USER_PORT")
	port2 := os.Getenv("PROJECTS_PORT")
	port3 := os.Getenv("GROUPS_PORT")
	
	//for users
	go startServer(port1, "Server 1")
	//for projects
	go startServer(port2, "Server 2")
	//for groups
	go startServer(port3, "Server 3")

	go startProxyServer(port0,"main")
	
	select {}
}
func InitializeRouter() {
	InitializeServers()
	
	

	// r.Handle("/addProjectMap", middleware.VerifyToken((http.HandlerFunc(controllers.AddNewProjectMapping)))).Methods("POST")
	
	fmt.Println(23145)
	
}