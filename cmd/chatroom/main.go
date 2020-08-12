/**
 * @Title  main
 * @description  #
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

func main(){
	fmt.Printf(banner + "\n" ,addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}