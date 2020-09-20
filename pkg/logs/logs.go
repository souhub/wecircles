package logs

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func warn(msg string) {
	// ログ例2（出力される）
	log.WithFields(log.Fields{
		"file": GetCurrentFile(),
		"line": GetCurrentFileLine(),
	}).Warn(msg)
}

func fatal(msg string) {
	// ログ例3（出力される）
	log.WithFields(log.Fields{
		"file": GetCurrentFile(),
		"line": GetCurrentFileLine(),
	}).Fatal(msg)
}

func init() {
	// JSONフォーマット
	log.SetFormatter(&log.JSONFormatter{})

	// 標準エラー出力でなく標準出力とする
	log.SetOutput(os.Stdout)

	// Warningレベル以上を出力
	log.SetLevel(log.WarnLevel)
}
