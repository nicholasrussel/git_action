package main

import (
	controllers "git_action/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//get endpoint
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/getDetailTransactions", controllers.GetDetailTransaction).Methods("GET")

	//insert endpoint
	router.HandleFunc("/users/insertNewUser", controllers.InsertNewUser).Methods("POST")
	router.HandleFunc("/products/insertNewProduct", controllers.InsertNewProduct).Methods("POST")
	router.HandleFunc("/transactions/insertNewTransaction", controllers.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/transactions/insertTransaction", controllers.InsertTransaction).Methods("POST")

	//delete endpoint
	router.HandleFunc("/users/deleteUser/{userID}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/products/deleteProduct/{productID}", controllers.DeleteProducts).Methods("DELETE")
	router.HandleFunc("/transactions/deleteTransaction/{transactionID}", controllers.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/products/deleteSingleProduct/{productID}", controllers.DeleteSingleProduct).Methods("DELETE")

	//update endpoint
	router.HandleFunc("/users/updateUser/{userID}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/updateProduct/{productID}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/users/updateTransaction/{transactionID}", controllers.UpdateTransaction).Methods("PUT")

	//login
	router.HandleFunc("/users/login", controllers.Login).Methods("POST")

	log.Println("hello world!")
	http.ListenAndServe(":8000", router)

}
