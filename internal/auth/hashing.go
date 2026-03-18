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

func VerifyPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}