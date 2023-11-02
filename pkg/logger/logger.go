package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	filename string = "./logs/info.log"
	dirName  string = "logs"
)

//type Logger struct {
//	l *log.Logger
//}

func New() (*log.Logger, error) {
	if err := os.Mkdir(dirName, 0750); err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("New logger: %v", err)
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return nil, fmt.Errorf("New logger: %v", err)
	}
	// defer file.Close()

	return log.New(file, "APP ", log.LstdFlags|log.Lshortfile), nil
}
