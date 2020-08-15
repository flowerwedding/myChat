/**
 * @Title  home
 * @description  其他的router函数
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"myChat/logic"
	"net/http"
)

// @Summary  模板
// @Router  / [get]
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

// @Summary  获取用户列表
// @Produce  json
// @Success  200 {object} logic.User "成功"
// @Router  /user_list [get]
func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")//gin框架的c.JSON也是这个原理
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		_, _ = fmt.Fprint(w, `[]`)
	} else {
		_, _ = fmt.Fprint(w, string(b))
	}
}