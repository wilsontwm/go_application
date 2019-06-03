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
	"time"
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
	ResetPasswordCode *string `json:"resetPasswordCode"`
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
	fmt.Println(user.ActivationCode, user.ActivationCode == nil)
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

func (user *User) ForgetPassword() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	// Get the user by email
	user = GetUserByEmail(user.Email)

	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid email address.", errors)
	} else if user.ActivationCode != nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "The account has not been activated yet. Please activate the account first.", errors)
	} else {
		// Store the reset password code to the user
		hash := md5.New()
		hash.Write([]byte(fmt.Sprint(user.ID) + time.Now().String()))	
		resetPasswordCode := hex.EncodeToString(hash.Sum(nil))
		
		GetDB().Model(&user).Update("ResetPasswordCode", resetPasswordCode)
		
		resp = util.Message(true, http.StatusOK, "An email to reset password has been sent to you. Please check your inbox.", errors)		
		resp["data"] = user
	}

	return resp
}

func (user *User) ActivateAccount(code string) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	// Get the user by email
	user = GetUserByActivationCode(code)
	
	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid activation link.", errors)
	} else {
		// Reset the activation code of the user		
		GetDB().Model(&user).Update("ActivationCode", nil)
		
		resp = util.Message(true, http.StatusOK, "Thank you for signing up. Your account has been activated.", errors)	
	}

	return resp
}

func (user *User) ResetPassword(code string, password string) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	// Get the user by reset password code
	user = GetUserByResetPasswordCode(code)
	
	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid/expired reset password link.", errors)
	} else {
		// Reset the password of the user	
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)	
		GetDB().Model(&user).Update(map[string]interface{}{"ResetPasswordCode": nil, "Password": string(hashedPassword) })
		
		resp = util.Message(true, http.StatusOK, "Successfully reset the password.", errors)	
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

func GetUserByActivationCode(activationCode string) *User {
	user := &User{}
	GetDB().Table("users").Where("activation_code = ?", activationCode).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""
	return user
}

func GetUserByResetPasswordCode(resetPasswordCode string) *User {
	user := &User{}
	GetDB().Table("users").Where("reset_password_code = ?", resetPasswordCode).First(user)
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
