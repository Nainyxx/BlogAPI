// Input validation helpers for user-provided fields (name, login, email, image URL).
package utils

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"
	"unicode"
)

func IsEmailValid(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func IsNameValid(name string) error {
	if ContainsOnlyLetter(name) != nil {
		return errors.New("invalid name")
	}
	return nil
}

func ContainsOnlyLetter(str string) error {
	for _, v := range str {
		if !unicode.IsLetter(v) {
			return errors.New("invalid name")
		}
	}
	return nil
}

func IsLoginValid(login string) error {
	for _, v := range login {
		if !(unicode.IsLetter(v) || unicode.IsNumber(v) || string(v) == "_" || string(v) == "-") {
			return errors.New("invalid login")
		}
	}
	return nil
}

func IsImageURLValid(imageURL string) error {
	if imageURL == "" {
		return errors.New("image URL cannot be empty")
	}

	u, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL must use http or https scheme")
	}

	path := strings.ToLower(u.Path)
	isImage := false

	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			isImage = true
			break
		}
	}

	if !isImage {
		return errors.New("URL must point to a valid image file (.jpg, .png, .webp, etc.)")
	}

	return nil
}
