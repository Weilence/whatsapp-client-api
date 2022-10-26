package utils

import (
	"io"
	"log"
)

func Close(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
