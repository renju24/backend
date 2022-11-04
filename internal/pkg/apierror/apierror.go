package apierror

import (
	"github.com/centrifugal/centrifuge"
)

var (
	ErrorInvalidBody                = &centrifuge.Error{400, "invalid JSON body", false}
	ErrorInternal                   = &centrifuge.Error{401, "internal server error", true}
	ErrorInvalidCredentials         = &centrifuge.Error{402, "invalid credentials", false}
	ErrorUnauthorized               = &centrifuge.Error{403, "invalid token", false}
	ErrorUserNotFound               = &centrifuge.Error{404, "user not found", false}
	ErrorUsernameIsRequired         = &centrifuge.Error{405, "username is required", false}
	ErrorUsernameIsTaken            = &centrifuge.Error{406, "username is already taken", false}
	ErrorEmailIsRequired            = &centrifuge.Error{407, "email is required", false}
	ErrorEmailIsTaken               = &centrifuge.Error{408, "email is already taken", false}
	ErrorPasswordIsRequired         = &centrifuge.Error{409, "password is required", false}
	ErrorRepeatedPasswordIsRequired = &centrifuge.Error{410, "repeated_password is required", false}
	ErrorInvalidUsernameLength      = &centrifuge.Error{411, "invalid username length", false}
	ErrorInvalidEmail               = &centrifuge.Error{412, "invalid email", false}
	ErrorInvalidEmailLength         = &centrifuge.Error{413, "invalid email length", false}
	ErrorInvalidPasswordLength      = &centrifuge.Error{414, "invalid password length", false}
	ErrorInvalidPasswordCharacter   = &centrifuge.Error{415, "invalid password character", false}
	ErrorMissingLetterInPassword    = &centrifuge.Error{416, "missing letter character", false}
	ErrorMissingDigitInPassword     = &centrifuge.Error{417, "missing digit character", false}
	ErrorPasswordsNotEqual          = &centrifuge.Error{418, "passwords are not equal", false}
	ErrorCallingYourselfForGame     = &centrifuge.Error{419, "can't call yourself for a game", false}
)
