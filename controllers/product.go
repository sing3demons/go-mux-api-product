package controllers

import (
	"app/models"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type H map[string]interface{}

type Product struct {
	DB *gorm.DB
}

type productRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

func (p *Product) FindAll(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	p.DB.Find(&products)

	serializedProducts := []productRespons{}
	copier.Copy(&serializedProducts, &products)

	JSON(w, http.StatusOK)(H{"products": serializedProducts})
}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product models.Product

	// json.NewDecoder(r.Body).Decode(&form)

	product.Name = r.FormValue("name")
	product.Desc = r.FormValue("desc")
	price, _ := strconv.Atoi(r.FormValue("price"))
	product.Price = price

	if err := p.DB.Create(&product).Error; err != nil {
		JSON(w, http.StatusNotFound)(H{"error": err.Error()})
	}

	p.saveProductImage(r, &product)

	// json.NewEncoder(w).Encode(product)
	JSON(w, http.StatusOK)(H{"product": product})

}

func (p *Product) FindOne(w http.ResponseWriter, r *http.Request) {
	product, err := p.findProductByID(r)
	if err != nil {
		JSON(w, http.StatusNotFound)(H{"error": err.Error()})
	}

	serializedProduct := []productRespons{}
	copier.Copy(&serializedProduct, &product)
	JSON(w, http.StatusOK)(H{"product": serializedProduct})
}

func (p *Product) findProductByID(r *http.Request) (*models.Product, error) {
	var product models.Product
	id := mux.Vars(r)["id"]

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

func (p *Product) saveProductImage(r *http.Request, product *models.Product) error {
	file, handler, err := r.FormFile("image")

	if err != nil || file == nil {
		panic(err)
	}
	defer file.Close()

	path := "uploads/products/" + strconv.Itoa(int(product.ID))
	os.Mkdir(path, 0755)
	filename := path + "/" + handler.Filename
	product.Image = os.Getenv("HOST") + "/" + filename

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// _, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)

	if err := p.DB.Save(product).Error; err != nil {
		return err
	}
	return nil

}

func JSON(w http.ResponseWriter, statusCode int) func(v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}
