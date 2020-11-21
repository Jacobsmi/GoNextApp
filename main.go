package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strconv"
)

// Create a user struct to create instances of users
type User struct{
	First string
	Last string
}
type Response struct {
	Success bool
}
var DB *sql.DB

func GetAllUsers(w http.ResponseWriter, r *http.Request)  {

}

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	//var userAdded Response
	var userAdded Response
	// Create an instance of a user to hold data for new user being added
	var u User
	// Decode json data that is being sent in the request
	decodeErr := json.NewDecoder(r.Body).Decode(&u)
	// Catch error in decoding
	if decodeErr != nil {
		fmt.Println(decodeErr)
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		panic("There was an error when decoding data")
	}
	execString := fmt.Sprintf("INSERT INTO users(first, last) VALUES('%s', '%s')", u.First, u.Last)
	_, execErr := DB.Exec(execString)
	if execErr != nil {
		fmt.Println(execErr)
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		panic("There was an error inserting user into DB in CreateUser")
	}
	// Return json {Success: true} if everything works
	userAdded.Success = true
	json.NewEncoder(w).Encode(userAdded)
}

func CreateDatabase()  {
	// Create a connection with the env file so we can get vars from it
	envErr := godotenv.Load(".env")
	// Catch errors opening .env file
	if envErr != nil{
		fmt.Println(envErr)
		panic("Error opening .env in Create Database")
	}
	// Convert the port from a string to an int
	port, stringToIntErr := strconv.Atoi(os.Getenv("DB_PORT"))
	// Catch any errors converting string to an int
	if stringToIntErr != nil{
		fmt.Println(stringToIntErr)
		panic("There was an error converting port to an int")
	}
	// Create a connection string from env vars by formatting a string with vars
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	// Connect to the DB
	var dbErr error
	DB, dbErr = sql.Open("postgres", connectionString)
	// Catch any errors connecting to DB
	if dbErr != nil {
		fmt.Println(dbErr)
		panic("There was a fatal error while opening the database")
	}
	_, createTableError := DB.Exec("CREATE TABLE IF NOT EXISTS users(" +
		"first VARCHAR(30) NOT NULL," +
		"last VARCHAR(30) NOT NULL)")
	if createTableError != nil {
		fmt.Println(createTableError)
		panic("Error creating user table in CreateDatabase")
	}
}

func handleRequests(){
	r := mux.NewRouter()
	r.HandleFunc("/users", GetAllUsers)
	r.HandleFunc("/createuser", CreateUser)
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Println(err)
		panic("There was an error when starting the server")
	}
}

func main()  {
	CreateDatabase()
	handleRequests()
}
