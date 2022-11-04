package apierror

import (
	"github.com/centrifugal/centrifuge"
)

var (
	ErrorInternal                   = centrifuge.ErrorInternal
	ErrorUnauthorized               = centrifuge.ErrorUnauthorized
	ErrorPermissionDenied           = centrifuge.ErrorPermissionDenied
	ErrorMethodNotFound             = centrifuge.ErrorMethodNotFound
	ErrorAlreadySubscribed          = centrifuge.ErrorAlreadySubscribed
	ErrorBadRequest                 = centrifuge.ErrorBadRequest
	ErrorUsernameIsRequired         = &centrifuge.Error{401, "username is required", false}
	ErrorEmailIsRequired            = &centrifuge.Error{402, "email is required", false}
	ErrorPasswordIsRequired         = &centrifuge.Error{403, "password is required", false}
	ErrorRepeatedPasswordIsRequired = &centrifuge.Error{404, "repeated_password is required", false}
	ErrorInvalidUsernameLength      = &centrifuge.Error{405, "invalid username length", false}
	ErrorInvalidEmail               = &centrifuge.Error{406, "invalid email", false}
	ErrorInvalidEmailLength         = &centrifuge.Error{407, "invalid email length", false}
	ErrorInvalidPasswordLength      = &centrifuge.Error{408, "invalid password length", false}
	ErrorInvalidPasswordCharacter   = &centrifuge.Error{409, "invalid password character", false}
	ErrorMissingLetterInPassword    = &centrifuge.Error{410, "missing letter character", false}
	ErrorMissingDigitInPassword     = &centrifuge.Error{411, "missing digit character", false}
	ErrorPasswordsNotEqual          = &centrifuge.Error{412, "passwords are not equal", false}
	ErrorUsernameIsTaken            = &centrifuge.Error{413, "username is already taken", false}
	ErrorEmailIsTaken               = &centrifuge.Error{414, "email is already taken", false}
	ErrorInvalidCredentials         = &centrifuge.Error{415, "invalid credentials", false}
	ErrorUserNotFound               = &centrifuge.Error{416, "user not found", false}
	ErrorCallingYourselfForGame     = &centrifuge.Error{417, "can't call yourself for a game", false}
)
