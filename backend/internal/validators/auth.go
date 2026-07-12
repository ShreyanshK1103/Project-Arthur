package validators

import (
	"errors"
	"net/mail"
	"unicode"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return errors.New("invalid email address")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)
	for _, c := range password {

		switch {
		case unicode.IsUpper(c):
			hasUpper = true

		case unicode.IsLower(c):
			hasLower = true

		case unicode.IsDigit(c):
			hasDigit = true

		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain an uppercase letter")
	}

	if !hasLower {
		return errors.New("password must contain a lowercase letter")
	}

	if !hasDigit {
		return errors.New("password must contain a digit")
	}

	if !hasSpecial {
		return errors.New("password must contain a special character")
	}

	return nil
}
