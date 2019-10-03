package policy

import (
	"app/models"
	"github.com/satori/go.uuid"
)

// Check if the user can create post
func CreatePost(userId, companyId uuid.UUID) bool {
	// Check if the user belongs to the company
	company := models.GetCompany(companyId, userId)

	return company != nil
}

// Check if the post belongs to the company
func ShowPost(userId, postId, companyId uuid.UUID) bool {
	db := models.GetDB()
	defer db.Close()

	// Check if the post belongs to the company
	post := models.Post{}
	db.Table("posts").
		Where("id = ? AND company_id = ?", postId, companyId).
		Scan(&post)

	return post.ID != uuid.Nil
}

// Check if the user can update / delete post
func UpdateDeletePost(userId, postId, companyId uuid.UUID) bool {
	db := models.GetDB()
	defer db.Close()

	// Check if the post belongs to the user
	post := models.Post{}
	db.Table("posts").
		Where("id = ? AND author_id = ? AND company_id = ?", postId, userId, companyId).
		Scan(&post)

	return post.ID != uuid.Nil
}
