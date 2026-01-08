package bcrypt

import basebcrypt "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := basebcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ComparePasswords(password, providedPassword string) bool {
	err := basebcrypt.CompareHashAndPassword([]byte(password), []byte(providedPassword))
	return err == nil
}
