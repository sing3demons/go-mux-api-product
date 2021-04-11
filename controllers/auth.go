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

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name" `
}

func (a *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
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

func (a *Auth) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	var user models.User
	if err := a.DB.First(&user, id).Error; err != nil {
		JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
		return
	}
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	JSON(w, http.StatusOK)(Map{"user": serializedUser})
}

func (a *Auth) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var form updateUserForm
	id := r.Header.Get("id")
	var user models.User

	if err := a.DB.First(&user, id).Error; err != nil {
		JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
		return
	}

	form.Name = r.FormValue("name")
	if form.Name == "" {
		form.Name = user.Name
	}
	
	form.Email = r.FormValue("email")
	if form.Email == "" {
		form.Email = user.Email
	}

	a.DB.Model(&user).Updates(map[string]interface{}{"name": form.Name, "email": form.Email})

	setUsersImage(r, &user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	JSON(w, http.StatusOK)(Map{"user": serializedUser})

}

func (a *Auth) UpdateImageProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	var user models.User

	a.DB.First(&user, id)

	setUsersImage(r, &user)
	JSON(w, http.StatusCreated)(Map{"message": "success"})

}
