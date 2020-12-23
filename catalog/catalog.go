package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float64 `json:"price,string"`
}

type Products struct {
	Products []Product `json:"products"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func loadProducts() []Product {
	response, err := http.Get(productsUrl + "/products")
	if err != nil {
		fmt.Println("Erro de HTTP")
	}

	data, _ := ioutil.ReadAll(response.Body)

	var products Products
	json.Unmarshal(data, &products)

	return products.Products
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("template/catalog.html"))
	t.Execute(w, products)
}

func showProducts(w http.ResponseWriter, r *http.Request) {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", listProducts)
	r.HandleFunc("/products/{id}", showProducts)
	http.ListenAndServe(":8082", r)
}
