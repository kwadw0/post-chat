package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashpass), nil
}