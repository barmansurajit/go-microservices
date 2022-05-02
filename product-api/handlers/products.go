package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	data "github.com/barmansurajit/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
	}
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (*Products) getProducts(w http.ResponseWriter) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
}
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	product := new(data.Product)
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request", http.StatusBadRequest)
	}
	p.l.Printf("Product: %#v", product)
	data.AddProduct(product)
}
func (p Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
