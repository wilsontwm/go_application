package models

import (
	//"database/sql"
	"github.com/dgrijalva/jwt-go"
	util "app/utils"
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Name string `json:"name";gorm:"not null"`
	Email string `json:"email";gorm:"unique;not null"`
	Password string `json:"password";gorm:"not null"`
	Token *string `json:"token";sql:"-"`
	ActivationCode *string `json:"activationCode"`
}

// Validate the incoming details
func (user *User) ValidateSignup() (map[string] interface{}, bool) {
	var errors []string
	var resp map[string] interface{}
	
	// Email must be unique
	temp := &User{}

	// Check for errors and duplicate emails
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		resp = util.Message(false, http.StatusInternalServerError, "Connection error. Please retry.", errors)
		return resp, false
	}

	if temp.Email != "" {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Email address has already been taken.", errors)
		return resp, false
	}

	resp = util.Message(true, http.StatusOK, "Input has been validated.", errors)
	return resp, true
}

func (user *User) Create() (map[string] interface{}) {
	var errors []string

	// Validate the account first
	if resp, ok := user.ValidateSignup(); !ok {
		return resp;
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Token = nil

	GetDB().Create(user)

	if user.ID <= 0 {
		resp := util.Message(false, http.StatusInternalServerError, "Failed to create account, connection error.", errors)
		return resp;
	}

	// Store the activation code to the user
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(user.ID)))	
	activationCode := hex.EncodeToString(hash.Sum(nil))
	
	GetDB().Model(&user).Update("ActivationCode", activationCode)
	
	user.Password = "" // delete the password

	resp := util.Message(true, http.StatusOK, "You have successfully signed up.", errors)
	resp["data"] = user

	return resp
}

func (user *User) ResendActivation() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	// Get the user by email
	user = GetUserByEmail(user.Email)

	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid email address.", errors)
	} else if user.ActivationCode == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "The account has already been activated.", errors)
	} else {
		resp = util.Message(true, http.StatusOK, "The activation link has been emailed to you. Please check your inbox.", errors)		
		resp["data"] = user
	}

	return resp
}

/*
func Login(email, password string) (map[string] interface{}) {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Invalid email address or password.")
		}

		return u.Message(false, "Connection error. Please retry.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid email address or password.")
	}

	account.Password = "" // delete the password

	// Create new JWT token for the newly registered account
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	response := u.Message(true, "You have logged in.")
	response["account"] = account

	return response
}
*/

func GetUserByEmail(email string) *User {
	user := &User{}
	GetDB().Table("users").Where("email = ?", email).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""
	return user
}

func GetUser(u uint) *User {
	user := &User{}
	GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""
	return user
}
