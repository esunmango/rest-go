package main

import(
	"log"
	"net/http"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
	_ "github.com/denisenkom/go-mssqldb"
	"encoding/json"
	//"io/ioutil"
	"database/sql"
	"fmt"
)

// Book Struct (Model)
type Book struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

//Init Books vas as slice Book strcuct
var books []Book

// Get All Books
func getBooks(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(books)

}

// Get single books
func getBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r) //Get params
	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create New Book
func createBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
	books = append(books,book)
	json.NewEncoder(w).Encode(book)
}

// Update Book
func updateBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for index, item := range books{
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe
			books = append(books,book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}
// Delete Book
func deleteBook(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for index, item := range books{
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}



func test(w http.ResponseWriter, r *http.Request){

	condb, errdb := sql.Open("mssql", "server=SQL03\\DB03;user id=sa;password=avanceytec;database=AXTEST;")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	//

	rows, err := condb.Query("SELECT TOP 1 ct.*,dp.NAME FROM CustTable ct left join DIRPARTYTABLE dp on ct.PARTY = dp.RECID")
	if err != nil {
		log.Fatal(err)
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if (ok) {
				v = string(b)
			} else {
				v = val
			}

			fmt.Println(col, v)
		}
	}
	json.NewEncoder(w).Encode(values)
}

func main(){
	//Init Router
	r := mux.NewRouter()

	//Mock Data - @todo - implement DB
	books = append(books,Book{ID:"1",Isbn:"448743",Title:"Book of nothing",Author:&Author{Firstname:"John",Lastname:"Mango"}})
	books = append(books,Book{ID:"2",Isbn:"423454",Title:"Book of All",Author:&Author{Firstname:"Paco",Lastname:"koo"}})

	//Rout Handlers / Endpoints
	r.HandleFunc("/api/books",getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}",getBook).Methods("GET")
	r.HandleFunc("/api/books",createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}",updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}",deleteBook).Methods("DELETE")

	r.HandleFunc("/api/test",test).Methods("GET")

	log.Fatal(http.ListenAndServe(":9797",r))
}
