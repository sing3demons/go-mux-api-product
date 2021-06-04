package controllers

import (
	"github/sing3demons/go_mux_api/config"
	"github/sing3demons/go_mux_api/models"
	"github/sing3demons/go_mux_api/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Users struct {
	DB *gorm.DB
}

type updateUserForm struct {
	Email  string                `form:"email"`
	Avatar *multipart.FileHeader `form:"avatar"`
	Name   string                `form:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email" `
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

func (u *Users) FindAll(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("role")

	if role != "Admin" {
		utils.JSON(w, http.StatusForbidden)(Map{"error": "forbindeb"})
		return
	}

	var users []models.User
	if err := u.DB.Find(&users).Error; err != nil {
		utils.JSON(w, http.StatusUnauthorized)(Map{"error": err.Error()})
		return
	}
	var serializedUsers []userResponse
	copier.Copy(&serializedUsers, &users)
	utils.JSON(w, http.StatusOK)(Map{"users": serializedUsers})

}

func (u *Users) FindOne(w http.ResponseWriter, r *http.Request) {
	user, err := u.findUserByID(r)
	if err != nil {
		utils.JSON(w, http.StatusNotFound)(Map{"error": err.Error()})
	}
	var serializedUser []userResponse
	copier.Copy(&serializedUser, &user)
	utils.JSON(w, http.StatusOK)(Map{"user": serializedUser})
}
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {}
func (u *Users) Update(w http.ResponseWriter, r *http.Request) {}
func (u *Users) Delete(w http.ResponseWriter, r *http.Request) {}
func (u *Users) Promote(w http.ResponseWriter, r *http.Request) {

}
func (u *Users) Demote(w http.ResponseWriter, r *http.Request) {}

func setUsersImage(w http.ResponseWriter, r *http.Request, user *models.User) {
	file, handler, err := r.FormFile("avatar")
	if file == nil || err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity)(Map{"error": err.Error()})
		return
	}
	defer file.Close()

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.Mkdir(path, 0755)
	filename := path + "/" + handler.Filename
	user.Avatar = os.Getenv("HOST") + "/" + filename

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	if err := config.GetDB().Save(user).Error; err != nil {
		log.Fatal(err.Error())
		return
	}
}

func (u *Users) findUserByID(r *http.Request) (*models.User, error) {
	id := mux.Vars(r)["id"]
	var user models.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
