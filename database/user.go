package database

type User struct {
	ModelBase
	Username string `gorm:"unique"`
	Password string
	Email    string
	Role     RoleType
}

type RoleType string

const (
	ADMIN RoleType = "ADMIN"
	USER  RoleType = "USER"
)

// Create a new user
func CreateUser(user *User) error {
	return db.Create(user).Error
}

// Retrieve a user by ID
func GetUserByID(id string) (*User, error) {
	var user User
	err := db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(name string) (*User, error) {
	var user User
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckUserByName(name string) bool {
	var user User
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

// Update a user
func UpdateUser(user *User) error {
	return db.Save(user).Error
}

// Delete a user
func DeleteUser(user *User) error {
	return db.Delete(user).Error
}
