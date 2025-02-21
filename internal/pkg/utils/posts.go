package utils

import "gorm.io/gorm"

func GetUserLikedPosts(db *gorm.DB, userID uint) map[uint]bool {
	var likedPostIDs []uint
	db.Table("likes").Where("author_id = ?", userID).Pluck("post_id", &likedPostIDs)

	likedPosts := make(map[uint]bool)
	for _, id := range likedPostIDs {
		likedPosts[id] = true
	}
	return likedPosts
}

func GetUserLikedPost(db *gorm.DB, userID uint, postID uint) bool {
	var count int64
	db.Table("likes").
		Where("author_id = ?", userID).
		Where("post_id = ?", postID).
		Count(&count)

	return count > 0
}
