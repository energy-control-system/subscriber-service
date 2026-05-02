package subscriber

import (
	"errors"
	"fmt"
	"time"
)

const (
	accountNumberMinLength = 4
	accountNumberMaxLength = 20

	nameMinLength = 1
	nameMaxLength = 50

	innLegalLength  = 10
	innPersonLength = 12

	passportSeriesLength      = 4
	passportNumberLength      = 6
	passportIssuedByMaxLength = 200

	minAge = 18
	maxAge = 120
)

var (
	ErrAccountNumberTooShort           = errors.New("account number must be at least 4 characters long")
	ErrAccountNumberTooLong            = errors.New("account number must be at most 20 characters long")
	ErrAccountNumberInvalidChar        = errors.New("account number must contain only uppercase Russian letters, digits, and hyphens")
	ErrAccountNumberNoHyphen           = errors.New("account number must contain at least one hyphen")
	ErrAccountNumberConsecutiveHyphens = errors.New("account number must not contain consecutive hyphens")
	ErrAccountNumberStartNotLetter     = errors.New("account number must start with an uppercase Russian letter")
	ErrAccountNumberEndNotDigit        = errors.New("account number must end with a digit")

	ErrSurnameEmpty         = errors.New("surname must not be empty")
	ErrSurnameInvalidLength = errors.New("surname length must be between 1 and 50 characters")
	ErrSurnameInvalidChar   = errors.New("surname must contain only Russian letters, spaces, and hyphens")

	ErrNameEmpty         = errors.New("name must not be empty")
	ErrNameInvalidLength = errors.New("name length must be between 1 and 50 characters")
	ErrNameInvalidChar   = errors.New("name must contain only Russian letters, spaces, and hyphens")

	ErrPatronymicInvalidLength = errors.New("patronymic length must be between 1 and 50 characters")
	ErrPatronymicInvalidChar   = errors.New("patronymic must contain only Russian letters, spaces, and hyphens")

	ErrPhoneEmpty         = errors.New("phone number must not be empty")
	ErrPhoneInvalidPrefix = errors.New("phone number must start with +7, 7, or 8")
	ErrPhoneInvalidLength = errors.New("phone number must contain exactly 10 digits after the prefix")
	ErrPhoneInvalidChar   = errors.New("phone number must contain only digits after the prefix")

	ErrEmailInvalidFormat = errors.New("email has invalid format")

	ErrINNInvalidLength = errors.New("INN must be 10 or 12 digits long")
	ErrINNInvalidChar   = errors.New("INN must contain only digits")

	ErrPassportSeriesInvalidLength    = errors.New("passport series must be exactly 4 digits")
	ErrPassportSeriesInvalidChar      = errors.New("passport series must contain only digits")
	ErrPassportNumberInvalidLength    = errors.New("passport number must be exactly 6 digits")
	ErrPassportNumberInvalidChar      = errors.New("passport number must contain only digits")
	ErrPassportIssuedByEmpty          = errors.New("passport issued-by must not be empty")
	ErrPassportIssuedByTooLong        = errors.New("passport issued-by must be at most 200 characters")
	ErrPassportIssueDateInvalidFormat = errors.New("passport issue date must be in DD.MM.YYYY format")
	ErrPassportIssueDateInFuture      = errors.New("passport issue date must not be in the future")

	ErrBirthDateInvalidFormat = errors.New("birth date must be in YYYY-MM-DD format")
	ErrBirthDateInvalidAge    = errors.New("subscriber age must be between 18 and 120 years")
)

// isUpperRussianLetter checks if the rune is an uppercase Russian letter (А-Я, Ё).
func isUpperRussianLetter(r rune) bool {
	return (r >= 'А' && r <= 'Я') || r == 'Ё'
}

// isLowerRussianLetter checks if the rune is a lowercase Russian letter (а-я, ё).
func isLowerRussianLetter(r rune) bool {
	return (r >= 'а' && r <= 'я') || r == 'ё'
}

// isRussianLetter checks if the rune is a Russian letter (both cases).
func isRussianLetter(r rune) bool {
	return isUpperRussianLetter(r) || isLowerRussianLetter(r)
}

// isDigit checks if the rune is an ASCII digit (0-9).
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// ValidateAccountNumber validates the subscriber's account number according to business rules:
//   - Length must be between 4 and 20 characters
//   - Only uppercase Russian letters (А-Я, Ё), digits (0-9), and hyphens (-) are allowed
//   - Must contain at least one character of each type (letter, digit, hyphen)
//   - Must start with an uppercase Russian letter
//   - Must end with a digit
//   - Consecutive hyphens are not allowed
func ValidateAccountNumber(accountNumber string) error {
	runes := []rune(accountNumber)
	length := len(runes)

	if length < accountNumberMinLength {
		return ErrAccountNumberTooShort
	}
	if length > accountNumberMaxLength {
		return ErrAccountNumberTooLong
	}

	if !isUpperRussianLetter(runes[0]) {
		return ErrAccountNumberStartNotLetter
	}

	if !isDigit(runes[length-1]) {
		return ErrAccountNumberEndNotDigit
	}

	hasHyphen := false
	prevHyphen := false

	for _, r := range runes {
		switch {
		case isDigit(r):
			prevHyphen = false
		case isUpperRussianLetter(r):
			prevHyphen = false
		case r == '-':
			if prevHyphen {
				return ErrAccountNumberConsecutiveHyphens
			}
			hasHyphen = true
			prevHyphen = true
		default:
			return ErrAccountNumberInvalidChar
		}
	}

	// Note: hasLetter is guaranteed by the first-character check above.
	// Note: hasDigit is guaranteed by the last-character check above.
	// Only hasHyphen requires explicit verification.
	if !hasHyphen {
		return ErrAccountNumberNoHyphen
	}

	return nil
}

// ValidateAddSubscriberRequest validates all fields of the AddSubscriberRequest.
func ValidateAddSubscriberRequest(req AddSubscriberRequest) error {
	if err := ValidateAccountNumber(req.AccountNumber); err != nil {
		return fmt.Errorf("account number: %w", err)
	}

	if req.Surname == "" {
		return ErrSurnameEmpty
	}
	if sLen := len([]rune(req.Surname)); sLen < nameMinLength || sLen > nameMaxLength {
		return ErrSurnameInvalidLength
	}
	for _, r := range req.Surname {
		if !isRussianLetter(r) && r != '-' && r != ' ' {
			return ErrSurnameInvalidChar
		}
	}

	if req.Name == "" {
		return ErrNameEmpty
	}
	if nLen := len([]rune(req.Name)); nLen < nameMinLength || nLen > nameMaxLength {
		return ErrNameInvalidLength
	}
	for _, r := range req.Name {
		if !isRussianLetter(r) && r != '-' && r != ' ' {
			return ErrNameInvalidChar
		}
	}

	if req.Patronymic != "" {
		if pLen := len([]rune(req.Patronymic)); pLen < nameMinLength || pLen > nameMaxLength {
			return ErrPatronymicInvalidLength
		}
		for _, r := range req.Patronymic {
			if !isRussianLetter(r) && r != '-' && r != ' ' {
				return ErrPatronymicInvalidChar
			}
		}
	}

	phone := []rune(req.PhoneNumber)
	if len(phone) == 0 {
		return ErrPhoneEmpty
	}
	if phone[0] != '+' && phone[0] != '7' && phone[0] != '8' {
		return ErrPhoneInvalidPrefix
	}
	startIdx := 1
	if phone[0] == '+' {
		if len(phone) < 2 || phone[1] != '7' {
			return ErrPhoneInvalidPrefix
		}
		startIdx = 2
	}
	if len(phone[startIdx:]) != 10 {
		return ErrPhoneInvalidLength
	}
	for _, r := range phone[startIdx:] {
		if !isDigit(r) {
			return ErrPhoneInvalidChar
		}
	}

	atCount := 0
	atIdx := -1
	for i, r := range req.Email {
		if r == '@' {
			atCount++
			atIdx = i
		}
	}
	if atCount != 1 || atIdx < 1 {
		return ErrEmailInvalidFormat
	}
	local := req.Email[:atIdx]
	domain := req.Email[atIdx+1:]
	if len(local) == 0 || len(domain) < 3 {
		return ErrEmailInvalidFormat
	}
	hasDot := false
	for _, r := range domain {
		if r == '.' {
			hasDot = true
		}
	}
	if !hasDot {
		return ErrEmailInvalidFormat
	}

	inn := []rune(req.INN)
	if len(inn) != innLegalLength && len(inn) != innPersonLength {
		return ErrINNInvalidLength
	}
	for _, r := range inn {
		if !isDigit(r) {
			return ErrINNInvalidChar
		}
	}

	series := []rune(req.Passport.Series)
	if len(series) != passportSeriesLength {
		return ErrPassportSeriesInvalidLength
	}
	for _, r := range series {
		if !isDigit(r) {
			return ErrPassportSeriesInvalidChar
		}
	}

	passNum := []rune(req.Passport.Number)
	if len(passNum) != passportNumberLength {
		return ErrPassportNumberInvalidLength
	}
	for _, r := range passNum {
		if !isDigit(r) {
			return ErrPassportNumberInvalidChar
		}
	}

	if req.Passport.IssuedBy == "" {
		return ErrPassportIssuedByEmpty
	}
	if len([]rune(req.Passport.IssuedBy)) > passportIssuedByMaxLength {
		return ErrPassportIssuedByTooLong
	}

	issueDate, err := time.Parse("02.01.2006", req.Passport.IssueDate)
	if err != nil {
		return ErrPassportIssueDateInvalidFormat
	}
	if issueDate.After(time.Now()) {
		return ErrPassportIssueDateInFuture
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return ErrBirthDateInvalidFormat
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	if age < minAge || age > maxAge {
		return ErrBirthDateInvalidAge
	}

	return nil
}
