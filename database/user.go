package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
}

// Create a new user
func CreateUser(user *User) error {
	return db.Create(user).Error
}

// Retrieve a user by ID
func GetUserByID(id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update a user
func UpdateUser(user *User) error {
	return db.Save(user).Error
}

// Delete a user
func DeleteUser(user *User) error {
	return db.Delete(user).Error
}
