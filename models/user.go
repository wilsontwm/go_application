package models

import (
	util "app/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

const (
	formattedDate = "01/02/2006"
)

type Token struct {
	UserId uuid.UUID
	Expiry time.Time
	jwt.StandardClaims
}

type User struct {
	Base
	Name                  string     `json:"name" gorm:"not null"`
	Email                 string     `json:"email" gorm:"unique;not null"`
	Password              string     `json:"password" gorm:"not null"`
	ProfilePicture        string     `json:"profilePicture"`
	Token                 string     `json:"token" gorm:"-"`
	ActivationCode        *string    `json:"activationCode"`
	ResetPasswordCode     *string    `json:"resetPasswordCode"`
	ResetPasswordExpiryDT *time.Time `json:"resetPasswordExpiryDateTime"`
	Phone                 string     `json:"phone"`
	City                  string     `json:"city"`
	Country               int        `json:"country" gorm:"default:'0'"`
	Gender                int        `json:"gender" gorm:"default:'0'"`
	Birthday              *time.Time `json:"birthday"`
	BirthdayString        string     `json:"birthday_string" gorm:"-"`
	Bio                   string     `json:"bio" sql:"type:text"`
}

func (user *User) Login(email string, password string) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Get the user by email
	db := GetDB()
	db.Table("users").Where("email = ?", email).First(&user)
	// Also get the companies that the user is assigned to
	companies := []Company{}
	db.Table("companies").
		Joins("JOIN company_users ON company_users.company_id = companies.id").
		Select("companies.*").
		Where("company_users.user_id = ?", user.ID).
		Order("company_users.last_visited desc").
		Find(&companies)

	defer db.Close()

	if user.Email == "" {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid email address or password.", errors)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		// If password does not match
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid email address or password.", errors)
		} else {
			// Password matches
			user.Password = "" // remove the password

			// Create new JWT token for the newly registered account
			expiry := time.Now().Add(time.Hour * 2) // Only valid for 2 hours
			tk := &Token{UserId: user.ID, Expiry: expiry}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
			tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
			user.Token = tokenString

			resp = util.Message(true, http.StatusOK, "You have successfully logged in.", errors)
			resp["data"] = user
			resp["companies"] = companies
			resp["selectedCompany"] = nil
			if len(companies) > 0 {
				selectedCompany := companies[0]
				resp["selectedCompany"] = selectedCompany
				user.SelectCompany(&selectedCompany)
			}
		}
	}

	return resp
}

// Validate the incoming details for signup
func (user *User) ValidateSignup() (map[string]interface{}, bool) {
	var errors []string
	var resp map[string]interface{}

	// Email must be unique
	temp := &User{}

	// Check for errors and duplicate emails
	db := GetDB()
	err := db.Table("users").Where("email = ?", user.Email).First(temp).Error

	defer db.Close()

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

func (user *User) Create() map[string]interface{} {
	var errors []string

	// Validate the account first
	if resp, ok := user.ValidateSignup(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Token = ""

	db := GetDB()
	db.Create(user)

	if user.ID == uuid.Nil {
		resp := util.Message(false, http.StatusInternalServerError, "Failed to create account, connection error.", errors)
		return resp
	}

	// Store the activation code to the user
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(user.ID)))
	activationCode := hex.EncodeToString(hash.Sum(nil))

	db.Model(&user).Update("ActivationCode", activationCode)

	defer db.Close()

	user.Password = "" // delete the password

	resp := util.Message(true, http.StatusOK, "You have successfully signed up. An activation email will be sent to you.", errors)
	resp["data"] = user

	return resp
}

func (user *User) ResendActivation() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

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

func (user *User) ForgetPassword() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

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
		// Add one hour to the expiry date for reseting the password
		resetPasswordExpiryDT := time.Now().Local().Add(time.Hour * 1)

		db := GetDB()
		db.Model(&user).Update(map[string]interface{}{
			"ResetPasswordCode":     resetPasswordCode,
			"ResetPasswordExpiryDT": resetPasswordExpiryDT,
		})

		defer db.Close()

		resp = util.Message(true, http.StatusOK, "An email to reset password has been sent to you. Please check your inbox.", errors)
		resp["data"] = user
	}

	return resp
}

func (user *User) ActivateAccount(code string) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Get the user by activation code
	user = GetUserByActivationCode(code)

	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid activation link.", errors)
	} else {
		// Reset the activation code of the user
		db := GetDB()
		db.Model(&user).Update("ActivationCode", nil)

		defer db.Close()

		resp = util.Message(true, http.StatusOK, "Thank you for signing up. Your account has been activated.", errors)
	}

	return resp
}

func (user *User) ResetPassword(code string, password string) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Get the user by reset password code
	user = GetUserByResetPasswordCode(code)

	if user == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Invalid/expired reset password link.", errors)
	} else {
		// Reset the password of the user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		db := GetDB()
		db.Model(&user).Update(map[string]interface{}{
			"ResetPasswordCode":     nil,
			"ResetPasswordExpiryDT": nil,
			"Password":              string(hashedPassword),
		})

		defer db.Close()

		resp = util.Message(true, http.StatusOK, "Successfully reset the password.", errors)
	}

	return resp
}

func (user *User) EditProfile() map[string]interface{} {
	var errors []string

	db := GetDB()
	db.Model(&user).Update(map[string]interface{}{
		"Name":     user.Name,
		"Phone":    user.Phone,
		"City":     user.City,
		"Country":  user.Country,
		"Gender":   user.Gender,
		"Birthday": user.Birthday,
		"Bio":      user.Bio,
	})

	defer db.Close()

	resp := util.Message(true, http.StatusOK, "Successfully updated profile.", errors)
	resp["data"] = user

	return resp
}

func (user *User) UploadPicture() map[string]interface{} {
	var errors []string

	db := GetDB()
	db.Model(&user).Update(map[string]interface{}{
		"ProfilePicture": user.ProfilePicture,
	})

	defer db.Close()

	resp := util.Message(true, http.StatusOK, "Successfully uploaded profile picture.", errors)
	resp["data"] = user

	return resp
}

func (user *User) DeletePicture() map[string]interface{} {
	var errors []string

	db := GetDB()
	db.Model(&user).Update(map[string]interface{}{
		"ProfilePicture": user.ProfilePicture,
	})

	defer db.Close()

	resp := util.Message(true, http.StatusOK, "Successfully removed profile picture.", errors)
	resp["data"] = user

	return resp
}

func (user *User) EditPassword() map[string]interface{} {
	var errors []string
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	password := string(hashedPassword)

	db := GetDB()
	db.Model(&user).Update(map[string]interface{}{
		"Password": password,
	})

	defer db.Close()

	resp := util.Message(true, http.StatusOK, "Successfully updated password.", errors)

	return resp
}

// Get the list of company invitation requests for the user
func (user *User) GetCompanyInvitationList() map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	companyInvitationRequests := []CompanyInvitationRequestOutput{}

	db := GetDB()
	db.Table("company_invitation_requests").
		Joins("JOIN companies ON company_invitation_requests.company_id = companies.id").
		Joins("JOIN users on company_invitation_requests.sender_id = users.id").
		Select("company_invitation_requests.*, companies.name as company_name, users.name as sender_name, users.email as sender_email, TO_CHAR(company_invitation_requests.created_at, '"+util.DateSQLFormat+"') as timestamp").
		Where("company_invitation_requests.email = ?", user.Email).
		Order("company_invitation_requests.created_at desc").
		Find(&companyInvitationRequests)

	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully retrieved all the company invitation requests.", errors)
	resp["data"] = companyInvitationRequests

	return resp
}

// Set the last visisted company's datetime
func (user *User) SelectCompany(company *Company) map[string]interface{} {
	var errors []string
	var resp map[string]interface{}

	// Update the last visited timestamp of the user at the company
	db := GetDB()
	db.Model(CompanyUser{}).Where("company_id = ? AND user_id = ?", company.ID, user.ID).Updates(map[string]interface{}{
		"LastVisited": time.Now(),
	})

	defer db.Close()

	resp = util.Message(true, http.StatusOK, company.Name+" has been selected.", errors)
	resp["selectedCompany"] = &company
	return resp
}

// Return a flag to show if user is admin of a company
func (user *User) IsAdmin(company *Company) bool {
	result := CompanyResult{}
	db := GetDB()
	db.Raw("SELECT C.name, C.id as company_id, R.is_admin FROM companies C JOIN company_users CU ON CU.company_id = C.id JOIN roles R ON R.id = CU.role_id WHERE CU.user_id = ? AND CU.company_id = ? AND C.deleted_at is NULL ORDER BY C.name ASC", user.ID, company.ID).First(&result)
	defer db.Close()

	if result.CompanyID == uuid.Nil {
		return false
	}

	return result.IsAdmin
}

func getUser(user *User) *User {
	if user.Email == "" {
		return nil
	}

	user.Password = ""

	return user
}

func GetUserByEmail(email string) *User {
	user := &User{}
	db := GetDB()
	db.Table("users").
		Select("users.*, TO_CHAR(users.birthday, '"+util.DateSQLFormat+"') as birthday_string").
		Where("email = ?", email).
		First(user)
	defer db.Close()

	return getUser(user)
}

func GetUserByActivationCode(activationCode string) *User {
	user := &User{}
	db := GetDB()
	db.Table("users").
		Select("users.*, TO_CHAR(users.birthday, '"+util.DateSQLFormat+"') as birthday_string").
		Where("activation_code = ?", activationCode).
		First(user)
	defer db.Close()

	return getUser(user)
}

func GetUserByResetPasswordCode(resetPasswordCode string) *User {
	user := &User{}
	now := time.Now().Local()
	db := GetDB()
	db.Table("users").
		Select("users.*, TO_CHAR(users.birthday, '"+util.DateSQLFormat+"') as birthday_string").
		Where("reset_password_code = ?", resetPasswordCode).
		Where("reset_password_expiry_dt > ?", now).
		First(user)
	defer db.Close()

	return getUser(user)
}

func GetUser(u uuid.UUID) *User {
	user := &User{}
	db := GetDB()
	db.Table("users").
		Select("users.*, TO_CHAR(users.birthday, '"+util.DateSQLFormat+"') as birthday_string").
		Where("id = ?", u).
		First(user)
	defer db.Close()

	return getUser(user)
}
