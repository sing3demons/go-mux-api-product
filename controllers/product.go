package controllers

import (
	"github/sing3demons/go_mux_api/models"
	"github/sing3demons/go_mux_api/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ProductController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	FindOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type productController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) ProductController {
	return &productController{DB: db}
}


type productRespons struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

type pagingRespons struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

func (p *productController) FindAll(w http.ResponseWriter, r *http.Request) {
	// JwtVerify(w, r)
	// AuthMiddleware(w, r)
	// id := r.Header.Get("sub")
	// fmt.Print(id)
	var products []models.Product

	pagination := pagination{
		ctx:     r,
		query:   p.DB,
		records: &products,
	}
	paging := pagination.pagingResource()
	// p.DB.Order("id desc").Find(&products)

	serializedProducts := []productRespons{}
	copier.Copy(&serializedProducts, &products)

	utils.JSON(w, http.StatusOK)(Map{"products": pagingRespons{Items: serializedProducts, Paging: paging}})
}

func (p *productController) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	product.Name = r.FormValue("name")
	product.Desc = r.FormValue("desc")
	product.Price, _ = strconv.Atoi(r.FormValue("price"))

	if err := p.DB.Create(&product).Error; err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}

	p.saveProductImage(r, &product)

	utils.JSON(w, http.StatusOK)(Map{"product": product})

}

func (p *productController) FindOne(w http.ResponseWriter, r *http.Request) {
	product, err := p.findProductByID(r)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}

	serializedProduct := []productRespons{}
	copier.Copy(&serializedProduct, &product)
	utils.JSON(w, http.StatusOK)(Map{"product": serializedProduct})
}

func (p *productController) Update(w http.ResponseWriter, r *http.Request) {
	product, err := p.findProductByID(r)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}

	product.Name = r.FormValue("name")
	product.Desc = r.FormValue("desc")
	product.Price, _ = strconv.Atoi(r.FormValue("price"))

	if err := p.DB.Save(&product).Error; err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}

	p.saveProductImage(r, product)

	utils.JSON(w, http.StatusOK)(Map{"message": "update success"})

}

func (p *productController) Delete(w http.ResponseWriter, r *http.Request) {
	product, err := p.findProductByID(r)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}

	p.DB.Unscoped().Delete(&product)

	utils.JSON(w, http.StatusOK)(Map{"message": "deleted..."})
}

func (p *productController) findProductByID(r *http.Request) (*models.Product, error) {
	var product models.Product
	id := mux.Vars(r)["id"]

	if err := p.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

func (p *productController) checkProduckImage(product *models.Product) {
	if product.Image != "" {
		product.Image = strings.Replace(product.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + product.Image)
	}

}

func (p *productController) saveProductImage(r *http.Request, product *models.Product) error {
	file, handler, err := r.FormFile("image")

	if err != nil || file == nil {
		return err
	}
	defer file.Close()

	p.checkProduckImage(product)

	path := "uploads/products/" + strconv.Itoa(int(product.ID))
	os.Mkdir(path, 0755)
	filename := path + "/" + handler.Filename
	product.Image = os.Getenv("HOST") + "/" + filename

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	// _, _ = io.WriteString(w, "File "+fileName+" Uploaded successfully")
	_, _ = io.Copy(f, file)

	if err := p.DB.Save(product).Error; err != nil {
		return err
	}
	return nil

}
