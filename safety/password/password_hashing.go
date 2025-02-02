package password

import "golang.org/x/crypto/bcrypt"

type UtilsPasswordHashing struct {
}

func (u *UtilsPasswordHashing) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func (u *UtilsPasswordHashing) CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
