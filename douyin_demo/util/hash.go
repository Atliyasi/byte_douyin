package util

import "golang.org/x/crypto/bcrypt"

// PwdHash 对面密码进行bcrypt加密
func PwdHash(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil

}

// PwdVerify 验证密码
func PwdVerify(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
