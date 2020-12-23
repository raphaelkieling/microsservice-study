package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"checkout/queue"

	"github.com/gorilla/mux"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Order struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	ProductId string `json:"product_id"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func displayCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	response, err := http.Get(productsUrl + "/products/" + vars["id"])
	if err != nil {
		fmt.Println("Erro de HTTP")
	}

	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	json.Unmarshal(data, &product)

	t := template.Must(template.ParseFiles("template/view.html"))
	t.Execute(w, product)
}

func finish(w http.ResponseWriter, r *http.Request) {
	order := Order{
		Name:      r.FormValue("name"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
		ProductId: r.FormValue("product_id"),
	}

	data, _ := json.Marshal(order)
	fmt.Println(string(data))

	connection := queue.Connect()
	queue.Notify(data, "checkout_ex", "", connection)

	w.Write([]byte("Processou!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/finish", finish)
	r.HandleFunc("/{id}", displayCheckout)
	http.ListenAndServe(":8083", r)
}
