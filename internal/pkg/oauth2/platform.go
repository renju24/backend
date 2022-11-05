package oauth2

import "errors"

var ErrUnknownPlatform = errors.New("unknown platform")

type Platform string

const (
	Web     Platform = "web"
	Android Platform = "android"
	IOS     Platform = "ios"
)

func ParsePlatform(s string) (Platform, error) {
	switch s {
	case "web":
		return Web, nil
	case "android":
		return Android, nil
	}
	return "", ErrUnknownPlatform
}
