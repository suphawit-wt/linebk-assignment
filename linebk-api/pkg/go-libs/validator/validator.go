package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Instance holds the universal translator configuration for field validations.
type Instance struct {
	Uni *ut.UniversalTranslator
}

// New initializes and registers custom validators and translations for validation checking.
func New() (*Instance, error) {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)

	enTrans, found := uni.GetTranslator("en")
	if !found {
		return nil, fmt.Errorf("validator error: Cannot get an 'en' translator")
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("username", usernameValidator); err != nil {
			return nil, fmt.Errorf("register username validation error: %w", err)
		}

		if err := v.RegisterValidation("passwd", passwordValidator); err != nil {
			return nil, fmt.Errorf("register passwd validation error: %w", err)
		}

		if err := en_translations.RegisterDefaultTranslations(v, enTrans); err != nil {
			return nil, fmt.Errorf("validator error: %w", err)
		}

		if err := v.RegisterTranslation(
			"username",
			enTrans,
			registerUsernameTranslation,
			translateUsernameError,
		); err != nil {
			return nil, fmt.Errorf("register username translation error: %w", err)
		}

		if err := v.RegisterTranslation(
			"passwd",
			enTrans,
			registerPasswdTranslation,
			translatePasswdError,
		); err != nil {
			return nil, fmt.Errorf("register passwd translation error: %w", err)
		}
	}

	return &Instance{
		Uni: uni,
	}, nil
}

func usernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9_|-]{1,28}[a-zA-Z0-9])?$`)

	reservedUsernames := map[string]bool{
		"admin":         true,
		"administrator": true,
		"root":          true,
		"system":        true,
		"support":       true,
		"help":          true,
		"api":           true,
		"config":        true,
		"null":          true,
		"undefined":     true,
		"anonymous":     true,
		"partner":       true,
		"official":      true,
	}

	if len(username) < 3 || len(username) > 30 {
		return false
	}

	if strings.Contains(username, "@") {
		return false
	}

	if !usernameRegex.MatchString(username) {
		return false
	}

	if strings.Contains(username, "__") || strings.Contains(username, "--") ||
		strings.Contains(username, "_-") ||
		strings.Contains(username, "-_") {
		return false
	}

	if reservedUsernames[strings.ToLower(username)] {
		return false
	}

	return true
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 64 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func registerUsernameTranslation(trans ut.Translator) error {
	return trans.Add(
		"username",
		"{0} must be 3-30 characters, start/end with a letter or number, contain no '@' or spaces, and cannot be a reserved system name",
		true,
	)
}

func translateUsernameError(trans ut.Translator, fe validator.FieldError) string {
	t, _ := trans.T("username", fe.Field())
	return t
}

func registerPasswdTranslation(trans ut.Translator) error {
	return trans.Add(
		"passwd",
		"{0} must be between 8 and 64 characters long and include an uppercase letter, lowercase letter, number, and special character",
		true,
	)
}

func translatePasswdError(trans ut.Translator, fe validator.FieldError) string {
	t, _ := trans.T("passwd", fe.Field())
	return t
}
