package encrypt

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 对密码进行加密（加盐哈希）
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 是 10，可以根据安全需求调整（越高越慢）
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码是否与哈希匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
