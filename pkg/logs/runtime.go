package logs

import "runtime"

// ファイル名を取得するメソッド
func GetCurrentFile() string {
	_, file, _, _ := runtime.Caller(1)

	return file
}

// 行数を取得するメソッド
func GetCurrentFileLine() int {
	_, _, line, _ := runtime.Caller(1)

	return line
}
