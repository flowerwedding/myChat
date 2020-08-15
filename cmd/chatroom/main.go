/**
 * @Title  main
 * @description  main.main
 * @Author  沈来
 * @Update  2020/8/9 14:51
 **/
package main

import (
	"fmt"
	"log"
	"myChat/server"
	"net/http"
)

var (
	addr = ":2000"
	banner = `
     
    ----------                    /\    ----------
    |           |          |     /  \        |
    |           |          |    /    \       |
    |           |----------|   /------\      |
    |           |          |  /        \     |
    ----------  |          | /          \    |

Go 编程之旅 ———— 聊天室, start on %s
`
)

// @title  聊天室
// @version  1.0
// @description 《Go语言编程之旅》项目练习
// @termsOfService  https://github.com/flowerwedding/myChat
func main(){
	fmt.Printf(banner + "\n" ,addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}