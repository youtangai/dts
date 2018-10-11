package errors

import (
	"fmt"
	"log"
)

const (
	ErrorTag = "dts:ERROR:"
	PanicTag = "dts:PANIC:"
)

func baseErrorLog(message string, err error) error {
	text := fmt.Sprint(ErrorTag, message, err)
	log.Println(text)
	return fmt.Errorf(text)
}

func GrpcDialError(err error) error {
	return baseErrorLog("cannnot dial dts server", err)
}

func GetFilesError(dir string, err error) error {
	message := fmt.Sprintf("cannot get files in dir:%s", dir)
	return baseErrorLog(message, err)
}

func OpenFileError(filename string, err error) error {
	message := fmt.Sprintf("cannot open file:%s", filename)
	return baseErrorLog(message, err)
}

func GetStreamError(err error) error {
	return baseErrorLog("cannot get file transfer stream", err)
}

func FileReadError(err error) error {
	return baseErrorLog("cannot read file", err)
}

func StreamCloseError(err error) error {
	return baseErrorLog("error occuerd when closing stream", err)
}

func StreamSendError(err error) error {
	return baseErrorLog("cannot send filedata", err)
}

func StreamRecvError(err error) error {
	return baseErrorLog("cannot recv filedata", err)
}

func FileWriteError(err error) error {
	return baseErrorLog("cannot write file", err)
}

func ListenError(url string, err error) error {
	message := fmt.Sprintf("cannot listen url:%s", url)
	return baseErrorLog(message, err)
}

func GrpcServeError(err error) error {
	return baseErrorLog("cannot start grpc server", err)
}

func MkdirError(dir string, err error) error {
	message := fmt.Sprintf("cannot create dir:%s", dir)
	return baseErrorLog(message, err)
}

func basePanicLog(err error) {
	panic(err)
}
func NoMuchArgNumPanic(argNum int) {
	err := fmt.Errorf("no much arg num:%d", argNum)
	basePanicLog(err)
}
