package utils

import (
	"io"
)

func Close(closer io.Closer) {
	err := closer.Close()
	NoError(err)
}

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintError(err error) {
	if err != nil {
		panic(err)
	}
}
