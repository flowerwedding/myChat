/**
 * @Title  home
 * @description  #
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(w http.ResponseWriter,req *http.Request){
	tpl, err := template.ParseFiles(rootDir + "/template/home.html")
	if err != nil {
		_, _ = fmt.Fprint(w, "模板解析错误！")
		return
	}

	err = tpl.Execute(w, nil)
	if  err != nil {
		_, _ = fmt.Fprint(w, "模板解析错误！")
		return
	}
}