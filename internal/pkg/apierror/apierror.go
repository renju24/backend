package apierror

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	ErrorInvalidBody        = &Error{100, "invalid JSON body"}
	ErrorInternal           = &Error{101, "internal server error"}
	ErrorInvalidCredentials = &Error{103, "invalid credentials"}
	ErrorUnauthorized       = &Error{104, "invalid token"}
)

var (
	ErrorUsernameIsRequired         = &Error{201, "username is required"}
	ErrorEmailIsRequired            = &Error{202, "email is required"}
	ErrorPasswordIsRequired         = &Error{203, "password is required"}
	ErrorRepeatedPasswordIsRequired = &Error{204, "repeated_password is required"}

	ErrorInvalidUsernameLength    = &Error{205, "invalid username length"}
	ErrorInvalidEmailLength       = &Error{206, "invalid email length"}
	ErrorInvalidPasswordLength    = &Error{207, "invalid password length"}
	ErrorInvalidPasswordCharacter = &Error{208, "invalid password character"}

	ErrorMissingUpperInPassword = &Error{209, "missing upper character"}
	ErrorMissingLowerInPassword = &Error{210, "missing lower character"}
	ErrorMissingDigitInPassword = &Error{211, "missing digit character"}
	ErrorPasswordsNotEqual      = &Error{212, "passwords are not equal"}
)
