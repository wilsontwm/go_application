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

func (user User) IndexCompany() (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	// Get the companies for the user
	result := &[]CompanyResult{}
	db := GetDB()
	db.Raw("SELECT C.name, C.id as company_id, R.is_admin FROM companies C JOIN company_users CU ON CU.company_id = C.id JOIN roles R ON R.id = CU.role_id WHERE CU.user_id = ? ORDER BY C.name ASC", user.ID).Scan(&result)
	defer db.Close()
	
	resp = util.Message(true, http.StatusOK, "", errors)
	resp["companies"] = result

	return resp
}

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

func (company *Company) ShowCompany(id, userId uuid.UUID) (map[string] interface{}) {
	var errors []string
	var resp map[string] interface{}

	company = GetCompany(id, userId)

	if company == nil {
		resp = util.Message(false, http.StatusUnprocessableEntity, "No available result.", errors)
	} else {		
		resp = util.Message(true, http.StatusOK, "", errors)
		resp["data"] = company
	}
	
	return resp
}

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

func GetCompany(companyId, userId uuid.UUID) *Company {
	// Only retrieve the company if user is in current company
	company := &Company{}
	db := GetDB()
	db.Raw("SELECT * FROM companies C JOIN company_users CU ON CU.company_id = C.id WHERE CU.user_id = ? AND C.id = ? LIMIT 1", userId, companyId).Scan(company)
	defer db.Close()

	if company.ID == uuid.Nil {
		return nil
	}

	return company
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