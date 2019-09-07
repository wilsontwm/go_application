package models

import (
	"errors"
	util "app/utils"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Company struct {
	Base
	Name string `gorm:"not null;"`
	Slug string `gorm:"not null;"`
	Description string
	Email string
	Phone string
	Fax string
	Address string
	Roles []Role `gorm:"foreignkey:CompanyID"`
	Users []User `gorm:"many2many:company_users"`
	CompanyUsers []CompanyUser `gorm:"foreignkey:CompanyID"`
}

type CompanyResult struct {
	Name string
	CompanyID uuid.UUID
	IsAdmin bool
}

// Validate the incoming details for creation of company
func (company *Company) Validate() (map[string] interface{}, bool) {
	var errors []string
	var resp map[string] interface{}
	
	// Slug must be unique
	temp := &Company{}

	// Check for errors and duplicate slug
	db := GetDB()
	err := db.Table("companies").Where("slug = ? and id <> ?", company.Slug, company.ID).First(temp).Error
	defer db.Close()
	
	if err != nil && err != gorm.ErrRecordNotFound {
		resp = util.Message(false, http.StatusInternalServerError, "Connection error. Please retry.", errors)
		return resp, false
	}

	if temp.ID != uuid.Nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "Slug has already been taken.", errors)
		return resp, false
	}

	resp = util.Message(true, http.StatusOK, "Input has been validated.", errors)
	return resp, true
}

// Get a list of the companies
func (user User) IndexCompany() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	// Get the companies for the user
	result := &[]CompanyResult{}
	db := GetDB()
	db.Raw("SELECT C.name, C.id as company_id, R.is_admin FROM companies C JOIN company_users CU ON CU.company_id = C.id JOIN roles R ON R.id = CU.role_id WHERE CU.user_id = ? AND C.deleted_at is NULL ORDER BY C.name ASC", user.ID).Scan(&result)
	defer db.Close()
	
	resp = util.Message(true, http.StatusOK, "", errors)
	resp["companies"] = result

	return resp
}

// Create the company
func (user User) CreateCompany(company *Company) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	// Validate the input first
	if resp, ok := company.Validate(); !ok {
		return resp;
	}

	if err := CreateCompanyTransaction(user, company); err != nil {
		resp = util.Message(false, http.StatusInternalServerError, err.Error(), errors)
		return resp
	}
		
	resp = util.Message(true, http.StatusOK, "You have successfully created a company. Invite people to your company now.", errors)
	resp["data"] = company

	return resp
}

// Get the company
func (company *Company) ShowCompany(id, userId uuid.UUID) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	company = GetCompany(id, userId)

	if company == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "No available result.", errors)
	} else {		
		resp = util.Message(true, http.StatusOK, "", errors)
		user := GetUser(userId)
		resp["data"] = company
		resp["isAdmin"] = user.IsAdmin(company)
	}
	
	return resp
}

// Update the company
func (company *Company) EditCompany() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	// Validate the input first
	if resp, ok := company.Validate(); !ok {
		return resp;
	}
	
	db := GetDB()
	db.Model(&company).Updates(company)
	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully updated company details.", errors)
	resp["data"] = company

	return resp
}

// Delete the company
func (company *Company) DeleteCompany() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	
	db := GetDB()
	db.Delete(&company)
	defer db.Close()

	resp = util.Message(true, http.StatusOK, "You have successfully deleted the company.", errors)

	return resp
}

// Send the invitation to emails to join the company
func (company *Company) InviteToCompany(email string, message string, senderId uuid.UUID) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	db := GetDB()
	// Check if email is already an user in the company for non-soft deleted
	companyUser := CompanyUser{}
	db.Raw("SELECT user_id, company_id, role_id FROM company_users CU JOIN users U ON U.id = CU.user_id WHERE U.email = ?", email).Scan(&companyUser)
	
	companyInvitationRequest := CompanyInvitationRequest{}
	db.Table("company_invitation_requests").Where("company_id = ? and email = ?", company.ID, email).First(&companyInvitationRequest)

	// If email is not in the company and not in the invitation list, create the invitation
	if(companyUser.UserID == uuid.Nil && companyInvitationRequest.Email == "") {
		companyInvitationRequest := CompanyInvitationRequest{
			CompanyID: company.ID,
			Email: email,
			Message: message,
			SenderID: &senderId,
		}

		db.Create(&companyInvitationRequest)
		resp = util.Message(true, http.StatusOK, "You have successfully invited " + email + " to the company.", errors)
		resp["data"] = companyInvitationRequest
	} else {
		resp = util.Message(false, http.StatusOK, "The user with the email " + email + " is already part of the company.", errors)
	}

	defer db.Close()

	return resp
}

// Get the company invitation list of the company
func (company *Company) GetCompanyInvitationList(page int) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	const resultsPerPage int = 25

	db := GetDB()
	companyInvitationRequests := []CompanyInvitationRequest{}

	if page <= 0 {
		db.Where("company_id = ?", company.ID).Order("created_at desc").Find(&companyInvitationRequests)
	} else {
		offset := resultsPerPage * ( page - 1 )
		db.Where("company_id = ?", company.ID).Order("created_at desc").Offset(offset).Limit(resultsPerPage).Find(&companyInvitationRequests)
	}

	defer db.Close()
	
	message := "You have successfully retrieved the invited emails to the company."
	if len(companyInvitationRequests) == 0 {
		message = "No more results."
	}

	resp = util.Message(true, http.StatusOK, message, errors)
	resp["data"] = companyInvitationRequests

	return resp
}

// Get the users in the company
func (company *Company) GetUserList(page int) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	const resultsPerPage int = 25

	db := GetDB()
	users := []User{}

	if page <= 0 {
		db.Table("users").
		Joins("JOIN company_users on company_users.user_id = users.id").
		Select("users.*").
		Where("company_users.company_id = ?", company.ID).
		Order("users.name asc").
		Find(&users)
	} else {
		offset := resultsPerPage * ( page - 1 )
		db.Table("users").
		Joins("JOIN company_users on company_users.user_id = users.id").
		Select("users.*").
		Where("company_users.company_id = ?", company.ID).
		Order("users.name asc").
		Offset(offset).
		Limit(resultsPerPage).
		Find(&users)
	}

	defer db.Close()
	
	message := "You have successfully retrieved the users of the company."
	if len(users) == 0 {
		message = "No more results."
	}

	resp = util.Message(true, http.StatusOK, message, errors)
	resp["data"] = users

	return resp
}

// Return the company if the user belongs to the company
func GetCompany(companyId, userId uuid.UUID) *Company {
	// Only retrieve the company if user is in current company
	company := &Company{}
	db := GetDB()
	db.Raw("SELECT * FROM companies C JOIN company_users CU ON CU.company_id = C.id WHERE CU.user_id = ? AND C.id = ? AND deleted_at is NULL LIMIT 1", userId, companyId).Scan(company)
	defer db.Close()

	if company.ID == uuid.Nil {
		return nil
	}

	return company
}

// Get the company based on ID
func GetCompanyByID(id uuid.UUID) *Company {
	comp := &Company{}
	db := GetDB()
	db.Table("companies").Where("id = ?", id).First(comp)
	defer db.Close()
	
	if comp.ID == uuid.Nil {
		return nil
	}

	return comp
}

// Get the unique slug/URL for the company name
func GetUniqueSlug(companyId uuid.UUID, slug string) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}
	companies := &[]Company{}

	db := GetDB()
	db.Raw("SELECT * FROM companies C WHERE C.id <> ? AND C.slug = ?", companyId, slug).Scan(companies)
	defer db.Close()

	message := "The slug is still available."
	isUniqueSlug := true
	if len(*companies) > 0 {
		message = "The slug has been taken."
		isUniqueSlug = false
	}

	resp = util.Message(true, http.StatusOK, message, errors)
	resp["is_unique"] = isUniqueSlug
	
	return resp
}

// The database transaction to create company
func CreateCompanyTransaction(user User, company *Company) error {
	db := GetDB()

	defer db.Close()
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	
	defer func() {
	  if r := recover(); r != nil {
		tx.Rollback()
	  }
	}()
  
	if err := tx.Error; err != nil {
	  return err
	}
  
	if err := tx.Create(&company).Error; err != nil {
	   tx.Rollback()
	   return err
	}	

	// Attach the admin & user roles as well
	admin := Role{Name: "Admin", IsAdmin: true, CompanyID: company.ID}
	normalUser := Role{Name: "User", CompanyID: company.ID}

	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		return err
	}	

	if err := tx.Create(&normalUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if admin.ID == uuid.Nil {
		tx.Rollback()
		err := errors.New("The admin role is not created in the company.")
		return err
	}

	// Associate the user to the company
	companyUser := CompanyUser{
		UserID: user.ID,
		CompanyID: company.ID,
		RoleID: admin.ID,
	}

	if err := tx.Create(&companyUser).Error; err != nil {
	   tx.Rollback()
	   return err
	}
	
	return tx.Commit().Error
}