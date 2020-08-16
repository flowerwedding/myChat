/**
 * @Title  receiveMessage
 * @description  #
 * @Author  沈来
 * @Update  2020/8/16 19:12
 **/
package logic

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func (b *broadcaster)ReceiveMessage(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil, )
	failOnError(err, "Failed to declare a queue")

	//key
	err = ch.QueueBind(q.Name, "*.message.*", "logs_topic", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil, )
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs {
			msg := &Message{}
			err := json.Unmarshal(d.Body, msg)
			if err != nil {
				log.Printf("Error decoding JSON: %s",err)
			}
			//log.Printf("Good: %s",string(d.Body))
			msg.User.Active = true
			if msg.To == "" {
				for _, user := range b.users{
					if user.UID == msg.User.UID {
						continue
					}

					Sending(nil,msg,"*."+user.Nickname+".*")
				}
			} else {
				if user, ok := b.users[msg.To]; ok {//私信
					Sending(nil,msg,"*."+user.Nickname+".*")
				} else {
					log.Println("user:",msg.To, "not exists!")
				}
			}
			if msg.Ats != nil {
				for _, str := range msg.Ats {//@
					if user, ok := b.users[str]; ok {//私信
						Sending(nil,NewNoticeMessage("你被@了"),"*."+user.Nickname+".*")
					} else {
						log.Println("user:",str, "not exists!")
					}
				}
			}

			OfflineProcessor.Save(msg)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}