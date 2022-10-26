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
	ErrorUsernameIsTaken            = &Error{200, "username is already taken"}
	ErrorEmailIsTaken               = &Error{201, "email is already taken"}
	ErrorUsernameIsRequired         = &Error{202, "username is required"}
	ErrorEmailIsRequired            = &Error{203, "email is required"}
	ErrorPasswordIsRequired         = &Error{204, "password is required"}
	ErrorRepeatedPasswordIsRequired = &Error{205, "repeated_password is required"}

	ErrorInvalidUsernameLength    = &Error{206, "invalid username length"}
	ErrorInvalidEmailLength       = &Error{207, "invalid email length"}
	ErrorInvalidPasswordLength    = &Error{208, "invalid password length"}
	ErrorInvalidPasswordCharacter = &Error{209, "invalid password character"}

	ErrorMissingUpperInPassword = &Error{210, "missing upper character"}
	ErrorMissingLowerInPassword = &Error{211, "missing lower character"}
	ErrorMissingDigitInPassword = &Error{212, "missing digit character"}
	ErrorPasswordsNotEqual      = &Error{213, "passwords are not equal"}
)
