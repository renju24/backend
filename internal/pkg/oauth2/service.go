package oauth2

import "errors"

var ErrUnknownService = errors.New("unknown service")

type Service string

const (
	Google Service = "google"
	Yandex Service = "yandex"
	VK     Service = "vk"
)

func ParseService(s string) (Service, error) {
	switch s {
	case "google":
		return Google, nil
	case "yandex":
		return Yandex, nil
	case "vk":
		return VK, nil
	}
	return "", ErrUnknownService
}
