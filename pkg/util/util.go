package util

import "github.com/google/uuid"

func UUID() string {
	return uuid.New().String()
}

func SHA1(s string) string {
	return uuid.NewSHA1(uuid.New(), []byte(s)).String()
}
