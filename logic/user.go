/**
 * @Title  user
 * @description  用户处理
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"myChat/global"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID         int               `json:"uid"`
	Active      bool              `json:"active,omitempty"`
	Nickname    string            `json:"nickname"`
	EnterAt     time.Time         `json:"enter_at"`
	Addr        string            `json:"addr"`

	Token       string        `json:"token,omitempty"`
	isNew       bool

	conn *websocket.Conn
}

var System = &User{UID: -1,Nickname:"System"}

func NewUser(conn *websocket.Conn,token string, nickname string, addr string) *User {
	user := &User{
		Nickname:       nickname,
		Active:         false,
		Addr:           addr,
		EnterAt:        time.Now(),
		Token:          token,

		conn: conn,
	}

	if user.Token != "" {
		uid, err := parseTokenAndValidate(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = genToken(user.UID, user.Nickname)
		user.isNew = true
	}

	return user
}

func (u *User) SendMessage(ctx context.Context){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	//key
	err = ch.QueueBind(q.Name, "*."+u.Nickname+".*", "logs_topic", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil )
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs {
			good := &Message{}
			err := json.Unmarshal(d.Body, good)
			if err != nil {
				log.Printf("Error decoding JSON: %s",err)
			}

			txt := good.User.Nickname + "(" + good.MsgTime.Format("15:04:05" ) + ")" + ":" + good.Content

			_ = wsjson.Write(ctx, u.conn, txt)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func (u *User) CloseMessageChannel(){

}

func (u *User) ReceiveMessage(ctx context.Context) error{
	var (
	//	receiveStr string
		receiveMsg map[string]string
		err        error
	)

	for {
	//	err = wsjson.Read(ctx, u.conn, &receiveStr)//输入字符串 "XXX"
		err = wsjson.Read(ctx, u.conn, &receiveMsg)//输入JSON格式 {"content":"XXX"}
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}

			return err
		}

	//	sendMsg := NewMessage(u, receiveStr)
		sendMsg := NewMessage(u, receiveMsg["content"])
		sendMsg.Content = FilterSensitive(sendMsg.Content)

		sendMsg.Content = strings.TrimSpace(sendMsg.Content)
		if strings.HasPrefix(sendMsg.Content, "#") {
			sendMsg.To = strings.SplitN(sendMsg.Content, " ",2)[0][1:]
		}

		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)

		Broadcaster.Broadcast(sendMsg)
	}
}

func FilterSensitive(content string) string{
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}