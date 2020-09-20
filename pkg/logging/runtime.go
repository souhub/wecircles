package logging

import "runtime"

func GetCurrentFileLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

func GetCurrentFile() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}
