package bcrypt

import "golang.org/x/crypto/bcrypt"

func HashValue(rawValue string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(rawValue), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(password)
	return hashPassword, nil
}

func ValidateHash(rawValue, hashValue string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(rawValue))
	return err
}
