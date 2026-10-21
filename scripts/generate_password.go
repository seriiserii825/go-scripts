package scripts

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
)

const (
	lowercase   = "abcdefghijklmnopqrstuvwxyz"
	uppercase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits      = "0123456789"
	punctuation = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

// GeneratePassword генерирует 10 паролей, даёт выбрать один в fzf и копирует
// его в буфер обмена (в уведомлении пароль не показывается).
func GeneratePassword() error {
	length, err := input.PositiveInt("Enter the length of the password: ")
	if err != nil {
		return err
	}
	symbols, err := input.Bool("Include symbols? (y/n): ")
	if err != nil {
		return err
	}

	passwords := make([]string, 10)
	for i := range passwords {
		passwords[i], err = newPassword(length, symbols)
		if err != nil {
			return err
		}
	}

	selected, ok, err := fzf.One(passwords)
	if err != nil {
		return err
	}
	if !ok {
		msg.Warning("No password selected.")
		return nil
	}
	fmt.Println("Selected password:", selected)
	return clipboard.WriteSecret(selected)
}

// newPassword собирает пароль из криптостойкого генератора (crypto/rand,
// как secrets в Python).
func newPassword(length int, symbols bool) (string, error) {
	chars := lowercase + digits + uppercase
	if symbols {
		chars += punctuation
	}
	limit := big.NewInt(int64(len(chars)))
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", err
		}
		password[i] = chars[n.Int64()]
	}
	return string(password), nil
}
