/**
 * @Title  handle
 * @description  #
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package server

import (
	"myChat/logic"
	"net/http"
	"os"
	"path/filepath"
)

var rootDir string

func RegisterHandle(){
	inferRootDir()

	go logic.Broadcaster.Start()

	http.HandleFunc("/",homeHandleFunc)
	http.HandleFunc("/ws",WebSocketHandleFunc)
}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string{
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}