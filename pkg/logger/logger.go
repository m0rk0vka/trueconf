package logger

import (
	"fmt"
	"log"
	"os"
)

const filename string = "./logs/info.log"

//type Logger struct {
//	l *log.Logger
//}

func New() (*log.Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, fmt.Errorf("New logger: %v", err)
	}
	// defer file.Close()

	return log.New(file, "APP ", log.LstdFlags|log.Lshortfile), nil
}
