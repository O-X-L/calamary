package log

import (
	"fmt"
	"time"

	"github.com/superstes/calamary/cnf"
)

func log(lvl string, pkg string, msg string) {
	var base string
	if cnf.C.Service.Debug {
		base = fmt.Sprintf("%s | %s | %s\n", lvl, pkg, msg)
	} else {
		base = fmt.Sprintf("%s | %s\n", lvl, msg)
	}

	if cnf.LOG_TIME {
		fmt.Printf("%s | %s", time.Now().Format(cnf.LOG_TIME_FORMAT), base)
	} else {
		fmt.Printf("%s\n", base)
	}
}

func logConn(lvl string, pkg string, src string, dst string, msg string) {
	var base string
	if cnf.C.Service.Debug {
		base = fmt.Sprintf("%s | %s | %s => %s | %s\n", lvl, pkg, src, dst, msg)
	} else {
		base = fmt.Sprintf("%s | %s => %s | %s\n", lvl, src, dst, msg)
	}

	if cnf.LOG_TIME {
		fmt.Printf("%s | %s", time.Now().Format(cnf.LOG_TIME_FORMAT), base)

	} else {
		fmt.Printf("%s", base)
	}
}

func ErrorS(pkg string, msg string) {
	log("ERROR", pkg, msg)
}

func Error(pkg string, err error) {
	log("ERROR", pkg, fmt.Sprintf("%s", err))
}

func ConnErrorS(pkg string, src string, dst string, msg string) {
	logConn("ERROR", pkg, src, dst, msg)
}

func ConnError(pkg string, src string, dst string, err error) {
	logConn("ERROR", pkg, src, dst, fmt.Sprintf("%s", err))
}

func Debug(pkg string, msg string) {
	if cnf.C.Service.Debug {
		log("DEBUG", pkg, msg)
	}
}

func ConnDebug(pkg string, src string, dst string, msg string) {
	if cnf.C.Service.Debug {
		logConn("DEBUG", pkg, src, dst, msg)
	}
}

func Info(pkg string, msg string) {
	log("INFO", pkg, msg)
}

func ConnInfo(pkg string, src string, dst string, msg string) {
	logConn("INFO", pkg, src, dst, msg)
}

func Warn(pkg string, msg string) {
	log("WARN", pkg, msg)
}

func Fatal(pkg string, msg string) {
	var base string
	if cnf.C.Service.Debug {
		base = fmt.Sprintf("FATAL | %s | %s\n", pkg, msg)
	} else {
		base = fmt.Sprintf("FATAL | %s\n", msg)
	}

	if cnf.LOG_TIME {
		panic(fmt.Sprintf("%s | %s", time.Now().Format(cnf.LOG_TIME_FORMAT), base))
	} else {
		panic(fmt.Sprintf("%s\n", base))
	}
}
