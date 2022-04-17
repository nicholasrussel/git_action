package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func SendResponse(w http.ResponseWriter, message string, status int) {
	var response MessageResponse
	response.Status = status
	response.Message = message

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	querry := "SELECT ID,name,age,address FROM users"
	rows, err := db.Query(querry)

	if err != nil {
		log.Println(err)
		return
	}

	var user User
	var users []User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Fatal(err.Error())
		} else {
			users = append(users, user)
		}
	}

	var response UsersResponse

	if len(users) < 1000 { //bisa apa saja ini cuman contohnya
		response.Status = 200
		response.Message = "Success"
		response.Data = users
	} else {
		response.Status = 400
		response.Message = "Error!"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()
	querry := "SELECT * FROM transactions"
	rows, err := db.Query(querry)

	if err != nil {
		SendResponse(w, "Internal Error", 400)
	}

	var transaction Transaction
	var transactions []Transaction
	for rows.Next() {
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity)
		if err != nil {
			SendResponse(w, "Internal Error", 400)
		} else {
			transactions = append(transactions, transaction)
		}
	}

	SendResponse(w, "Request Success", 200)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	query := "SELECT * FROM products"

	rows, err := db.Query(query)
	if err != nil {
		SendResponse(w, "Internal Error", 400)
	}

	var product Product
	var products []Product

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			products = append(products, product)
		}
	}

	SendResponse(w, "Request Success", 200)

}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	name := r.FormValue("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	_, errQuery := db.Exec("INSERT INTO users(name, age, address , email , password) VALUES (?,?,?,?,?)",
		name,
		age,
		address,
		email,
		password,
	)

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	SendResponse(w, "Insert Success", 200)
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	_, errQuery := db.Exec("INSERT INTO products(name, price) VALUES (?,?)",
		name,
		price,
	)

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	SendResponse(w, "Insert Success", 200)
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("userID"))
	productId, _ := strconv.Atoi(r.Form.Get("productID"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	_, errQuery := db.Exec("INSERT INTO transactions(userID, productID, quantity) VALUES (?,?,?)",
		userId,
		productId,
		quantity,
	)

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	SendResponse(w, "Insert Success", 200)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("userID"))
	productId, _ := strconv.Atoi(r.Form.Get("productID"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	rows, errQuery := db.Query("SELECT * FROM products WHERE id=?",
		productId,
	)

	if errQuery != nil {
		SendResponse(w, "Querry Error", 400)
	}

	if err != nil {
		SendResponse(w, "Error", 400)
		return
	}

	i := 0
	for rows.Next() {
		i++
	}

	if i == 0 {
		_, err = db.Exec("INSERT INTO products (id) VALUES (?)", productId)

		if err != nil {
			SendResponse(w, "Error", 400)
			return
		}
	}
	_, errQueries := db.Exec("INSERT INTO transactions (userid, productid, quantity) VALUES (?,?,?)", userId, productId, quantity)

	if errQueries != nil {
		SendResponse(w, "Internal Error", 400)
		return
	}

	SendResponse(w, "Insert Success", 200)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	userId := vars["userID"]

	result, errQuery := db.Exec("DELETE FROM users WHERE ID=?",
		userId,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return
	} else {
		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)
			return
		}
	}

	SendResponse(w, "Delete Success", 200)
}

func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	productId := vars["productID"]

	result, errQuery := db.Exec("DELETE FROM products WHERE ID =?",
		productId,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return
	} else {
		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)
			return
		}
	}

	SendResponse(w, "Delete Success", 200)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	transactionId := vars["transactionID"]

	result, errQuery := db.Exec("DELETE FROM transactions WHERE ID=?",
		transactionId,
	)
	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)
		return

	} else {

		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)
			return

		}
	}
	SendResponse(w, "Delete Success", 200)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	userId := vars["userID"]

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		log.Println(err)
		return
	}

	name := r.FormValue("name")
	age := r.Form.Get("age")
	address := r.Form.Get("address")

	result, errQuery := db.Exec("UPDATE users SET name=?, age=?, address=? WHERE id=?",
		name,
		age,
		address,
		userId,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)

		log.Println(errQuery)
		return
	} else {
		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)

			return
		}
	}

	SendResponse(w, "Update Success", 200)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	product_id := vars["product_id"]

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		log.Println(err)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	result, errQuery := db.Exec("UPDATE products SET name=?, price=? WHERE id=?",
		name,
		price,
		product_id,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {
		SendResponse(w, "Internal Error", 400)

		log.Println(errQuery)
		return
	} else {
		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)
			return
		}
	}

	SendResponse(w, "Update Success", 200)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {

	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	transactionId := vars["transaction_id"]

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Internal Error", 400)
		log.Println(err)
		return
	}

	userId, _ := strconv.Atoi(r.Form.Get("user_id"))
	productId, _ := strconv.Atoi(r.Form.Get("product_id"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	result, errQuery := db.Exec("UPDATE transactions SET user_id=?, product_id=?, quantity=? WHERE id=?",
		userId,
		productId,
		quantity,
		transactionId,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuery != nil {

		SendResponse(w, "Error", 400)
		log.Println(errQuery)
		return
	} else {
		if rowAffected == 0 {
			SendResponse(w, "Internal Error", 400)
			return
		}
	}

	SendResponse(w, "Update Success", 200)
}

func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	vars := mux.Vars(r)
	productId := vars["productID"]
	result, errQuerry := db.Exec("DELETE FROM transactions WHERE ID = ? ",
		productId,
	)
	result, errQuerry2 := db.Exec("DELETE FROM products WHERE ID = ?",
		productId,
	)

	rowAffected, _ := result.RowsAffected()

	if errQuerry != nil {
		SendResponse(w, "Internal Error", 400)
		return
	} else {
		if errQuerry2 != nil {
			SendResponse(w, "Internal Error", 400)
			return
		} else {
			if rowAffected == 0 {
				SendResponse(w, "Internal Error", 400)
				return
			}
		}
	}

	SendResponse(w, "Delete Success", 200)

}

func GetDetailTransaction(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	var transactionDetails []TransactionDetail
	query := "SELECT t.ID , u.ID, u.name, u.age, u.address, p.ID, p.name, p.price, t.quantity FROM transactions t JOIN users u ON t.userID = u.ID JOIN products p ON p.ID = t.productID"

	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE u.id = " + id[0]
	}

	rows, err := db.Query(query)

	var response TransactionDetailsResponse

	if err != nil {
		SendResponse(w, "Error Querry", 400)
		return
	}

	var transactionDetail TransactionDetail
	var user User
	var product Product

	for rows.Next() {
		if err := rows.Scan(&transactionDetail.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transactionDetail.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			transactionDetail.User = user
			transactionDetail.Product = product
			transactionDetails = append(transactionDetails, transactionDetail)
		}
	}

	if len(transactionDetails) != 0 {
		response.Status = 200
		response.Message = "Success Get Data"
		response.Data = transactionDetails
	} else {
		response.Status = 400
		response.Message = "Error Get Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		SendResponse(w, "Error Parsing", 400)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" {
		SendResponse(w, "please input Your Email ", 400)
		return
	}

	if password == "" {
		SendResponse(w, "please input your password", 400)
		return
	}
	rows, err := db.Query("SELECT email, password FROM users WHERE email= ?", email)

	if err != nil {
		SendResponse(w, "Error", 400)
		return
	}

	var user User
	var users []User

	for rows.Next() {
		if err := rows.Scan(&user.Email, &user.Password); err != nil {
			log.Println(err.Error())
		} else {
			users = append(users, user)
		}
	}

	if users[0].Password == password {
		SendResponse(w, "Login Success", 200)

	} else {
		SendResponse(w, "Login Failed", 400)
	}
}
