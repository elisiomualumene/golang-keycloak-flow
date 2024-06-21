package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gorilla/mux"
)

type Response struct {
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", good).Methods("GET")

	r.HandleFunc("/create", CreateKeycloakUser).Methods("GET")
	r.HandleFunc("/login", LoginKeycloakUser).Methods("POST")


	fmt.Println("Server is starting on port 3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func good(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(Response{
		Message: "logado com sucesso",
	})
}

func LoginKeycloakUser(w http.ResponseWriter, r *http.Request) {
	var data SignInRequest

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := gocloak.NewClient("http://localhost:8080")
	ctx := context.Background()

	token, err := client.Login(ctx, "rest-golang", "gUGw81NmB4n8X5k5wiy9XHvTlEdwrGir", "master", data.Username, data.Password)
	if err != nil {
		fmt.Println("Error logging in admin:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:  "Server is up and running",
		Response: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateKeycloakUser(w http.ResponseWriter, r *http.Request) {
	client := gocloak.NewClient("http://localhost:8080")
	ctx := context.Background()

	token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		fmt.Println("Error logging in admin:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token == nil {
		fmt.Println("Received nil token")
		http.Error(w, "Received nil token", http.StatusInternalServerError)
		return
	}

	user := gocloak.User{
		FirstName: gocloak.StringP("Elisio"),
		LastName:  gocloak.StringP("Mualumene"),
		Email:     gocloak.StringP("elisiomualumene@gmail.com"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("elisiomualumene"),
	}

	userID, err := client.CreateUser(ctx, token.AccessToken, "master", user)
	if err != nil {
		fmt.Println("Error creating user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Created user ID:", userID)

	password := "admin"
	err = client.SetPassword(ctx, token.AccessToken, userID, "master", password, false)
	if err != nil {
		fmt.Println("Error setting password:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:  "Server is up and running",
		Response: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
