package middleware

import (
	"app/config"
	"app/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type Map map[string]interface{}

type formLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	var form formLogin
	json.NewDecoder(r.Body).Decode(&form)

	var user models.User
	copier.Copy(&user, &form)
	db := config.GetDB()
	if err := db.Where("email = ?", form.Email).First(&user).Error; err != nil {
		JSON(w, http.StatusUnauthorized)(Map{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		JSON(w, http.StatusUnauthorized)(Map{"error": err.Error()})
		return
	}

	jwt, _ := jwtSign(user)
	serializedUser := jwt
	JSON(w, http.StatusOK)(Map{"token": serializedUser})

}

func jwtSign(user models.User) (string, error) {
	// Create token
	at := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := at.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Local().Unix()

	// Generate encoded token and send it as response.
	token, err := at.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", err
	}
	return token, nil

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if len(tokenString) == 0 {
			JSON(w, http.StatusUnauthorized)(Map{"error": "Missing Authorization Header"})
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := verifyToken(tokenString)

		if err != nil {
			JSON(w, http.StatusUnauthorized)(Map{"error": err.Error()})
			return
		}

		claims, ok := token.(jwt.MapClaims)
		if !ok {
			JSON(w, http.StatusUnauthorized)(Map{"error": "Missing Authorization Header"})
			return
		}

		id := fmt.Sprintf("%v", claims["id"])
		var user models.User
		db := config.GetDB()
		if err := db.First(&user, id).Error; err != nil {
			fmt.Println(err.Error())
		}

		role := user.Role

		r.Header.Set("id", id)
		r.Header.Set("sub", role)

		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func JSON(w http.ResponseWriter, statusCode int) func(v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}
