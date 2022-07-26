package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("starting the server...")
	http.HandleFunc("/get-employee", GetEmployeeHandler)
	http.HandleFunc("/post-employee", PostEmployeeHandler)

	http.ListenAndServe("127.0.0.1:8080", nil)

}

type Employee struct {
	Eid  string `json:"eid"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Handler for POST Request
func PostEmployeeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods is supported.")
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	employee := Employee{}

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		fmt.Println(err)
	}
	fmt.Println(employee)
	b, err := json.Marshal(employee)
	if err != nil {
		fmt.Println(err)
	}

	UtilWrite(string(b) + "\n")

	fmt.Fprintf(w, "Employee %v's Data saved successfully!", employee.Name)
}

// Handler for GET Request
func GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(UtilRead()); err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "Application/json")
}

// Util function to write a file
func UtilWrite(str string) {
	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	IfError(err)

	defer file.Close()

	_, err = file.WriteString(str)
	IfError(err)
}

// Util function to read from file
func UtilRead() []Employee {

	f, err := os.Open("data.txt")
	IfError(err)

	defer f.Close()

	scanner := bufio.NewScanner(f)

	data := []Employee{}

	for scanner.Scan() {
		emp := Employee{}
		json.Unmarshal([]byte(scanner.Text()), &emp)
		data = append(data, emp)
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data

}

func IfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
