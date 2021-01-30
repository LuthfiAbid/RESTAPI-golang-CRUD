package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/LuthfiAbid/golang_interview/api/auth"
	"github.com/LuthfiAbid/golang_interview/api/models"
	"github.com/LuthfiAbid/golang_interview/api/responses"
	"github.com/LuthfiAbid/golang_interview/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Username, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server)SignIn(username, password string) (string, error){
	var err error
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err!=nil{return "", err}
	err = models.VerifyPassword(user.Password, password)
	if err!=nil && err==bcrypt.ErrMismatchedHashAndPassword{
		return "", err
	}

	return auth.CreateToken(user.ID)
}