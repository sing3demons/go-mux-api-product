package seeds

import (
	"fmt"
	"github/sing3demons/go_mux_api/config"
	"github/sing3demons/go_mux_api/models"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
)

func Setup() {
	// config.GetDB().Migrator().DropTable(&models.Product{})
	config.GetDB().AutoMigrate(&models.Product{})
	fmt.Println("Creating products...")
	numberOfProduct := 100000

	products := make([]models.Product, numberOfProduct)
	for i := 0; i < numberOfProduct; i++ {
		product := models.Product{
			Name:  faker.Name(),
			Desc:  faker.Word(),
			Price: int(rand.Int63()),
			Image: "https://i.pravatar.cc/100?" + strconv.Itoa(i),
		}
		products[i] = product
	}

	config.GetDB().CreateInBatches(products, 1000)

	fmt.Println("success")

}
