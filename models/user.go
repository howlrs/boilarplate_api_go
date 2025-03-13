package models

import "golang.org/x/crypto/bcrypt"

const (
	COLLECTION_USER = "users"
)

type User struct {
	ID       string `json:"id" firestore:"id"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
}

func (u *User) ToCollection(isTest bool) string {
	if isTest {
		return "test_" + COLLECTION_USER
	}
	return COLLECTION_USER
}

// パスワードのハッシュ化
func (u *User) ToEncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// パスワードの検証
func (u *User) IsVerifyPassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword))
}
