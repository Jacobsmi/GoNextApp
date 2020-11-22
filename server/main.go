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
	ID int
	First string
	Last string
}
type Response struct {
	Success bool
}
var DB *sql.DB

// Handles error throwing
func throwError(e error, w http.ResponseWriter, s int){
	fmt.Println(e)
	if s != 0{
		http.Error(w, e.Error(), s)
	}
	panic(e.Error())
}

// Returns all users in the database in the form of JSON
func GetAllUsers(w http.ResponseWriter, r *http.Request)  {
	// Create an empty list of users
	var users []User
	// Query the database
	rows, queryErr := DB.Query("SELECT * FROM users")
	// Catch and handle errors
	if queryErr != nil {
		throwError(queryErr, w, http.StatusInternalServerError)
	}
	// For all the rows produced by the query get the user information, put it into a temp user, and then add it to
	// the users list
	for rows.Next(){
		var u User
		if scanErr := rows.Scan(&u.ID, &u.First, &u.Last); scanErr != nil{
			throwError(scanErr, w, http.StatusInternalServerError)
		}
		users = append(users, u)
	}
	// Encode the users list and send it back as JSON
	if encErr := json.NewEncoder(w).Encode(users); encErr != nil{
		throwError(encErr, w, http.StatusInternalServerError)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request)  {
	//var userAdded Response
	var userAdded Response
	// Create an instance of a user to hold data for new user being added
	var u User
	// Decode json data that is being sent in the request
	// Catch error in decoding
	if decodeErr := json.NewDecoder(r.Body).Decode(&u); decodeErr != nil {
		throwError(decodeErr, w, http.StatusBadRequest)
	}
	// Format string to be executed on the database
	execString := fmt.Sprintf("INSERT INTO users(first, last) VALUES('%s', '%s')", u.First, u.Last)
	// Execute the query and catch any errors
	if _, execErr := DB.Exec(execString); execErr != nil {
		throwError(execErr, w, http.StatusInternalServerError)
	}
	// Return json {Success: true} if everything works
	userAdded.Success = true
	if encodeErr:= json.NewEncoder(w).Encode(userAdded); encodeErr != nil{
		throwError(encodeErr, w, http.StatusInternalServerError)
	}
}

func CreateDatabase()  {
	// Create a connection with the env file so we can get vars from it and catch errors opening .env file
	if envErr := godotenv.Load("./server/.env"); envErr != nil{
		throwError(envErr, nil, 0)
	}
	// Convert the port from a string to an int
	port, stringToIntErr := strconv.Atoi(os.Getenv("DB_PORT"))
	// Catch any errors converting string to an int
	if stringToIntErr != nil{
		throwError(stringToIntErr, nil, 0)
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
		throwError(dbErr, nil, 0)
	}
	_, createTableError := DB.Exec("CREATE TABLE IF NOT EXISTS users(" +
		"id SERIAL," +
		"first VARCHAR(30) NOT NULL," +
		"last VARCHAR(30) NOT NULL)")

	if createTableError != nil {
		throwError(createTableError, nil, 0)
	}
}

func handleRequests(){
	r := mux.NewRouter()
	r.HandleFunc("/users", GetAllUsers)
	r.HandleFunc("/createuser", CreateUser)
	// Start the server and catch any errors starting the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		throwError(err, nil, 0)
	}
}

func main()  {
	CreateDatabase()
	handleRequests()
}
