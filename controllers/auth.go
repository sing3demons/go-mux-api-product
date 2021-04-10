package controllers

import (
	"app/models"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name" `
}

func (a *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.Password = user.GenerateEncryptedPassword()

	if err := a.DB.Create(&user).Error; err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	serializedAuth := authResponse{}
	copier.Copy(&serializedAuth, &user)

	JSON(w, http.StatusOK)(Map{"user": serializedAuth})
}
