package models

type Gender struct {
	GenderId int
	Sex string
}

var genders []Gender

func init() {
	// Initialize all the genders
	sexes := []string{
		"Unspecified",
		"Male",
		"Female",
	}

	for i, s := range sexes {
		gender := CreateGender(i+1, s)

		genders = append(genders, *gender)
	}
}

func CreateGender(id int, sex string) *Gender {
    return &Gender{GenderId: id, Sex: sex}
}

func GetGenders() []Gender {
	return genders
}