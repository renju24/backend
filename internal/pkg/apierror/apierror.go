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
	ErrorUsernameIsRequired         = &Error{200, "username is required"}
	ErrorUsernameIsTaken            = &Error{201, "username is already taken"}
	ErrorEmailIsRequired            = &Error{202, "email is required"}
	ErrorEmailIsTaken               = &Error{203, "email is already taken"}
	ErrorPasswordIsRequired         = &Error{204, "password is required"}
	ErrorRepeatedPasswordIsRequired = &Error{205, "repeated_password is required"}

	ErrorInvalidUsernameLength    = &Error{206, "invalid username length"}
	ErrorInvalidEmail             = &Error{207, "invalid email"}
	ErrorInvalidEmailLength       = &Error{208, "invalid email length"}
	ErrorInvalidPasswordLength    = &Error{209, "invalid password length"}
	ErrorInvalidPasswordCharacter = &Error{210, "invalid password character"}

	ErrorMissingLetterInPassword = &Error{211, "missing letter character"}
	ErrorMissingUpperInPassword  = &Error{212, "missing upper character"}
	ErrorMissingLowerInPassword  = &Error{213, "missing lower character"}
	ErrorMissingDigitInPassword  = &Error{214, "missing digit character"}
	ErrorPasswordsNotEqual       = &Error{215, "passwords are not equal"}
)
