/**
 * @Title  receiveEnter
 * @description  #
 * @Author  沈来
 * @Update  2020/8/16 19:11
 **/
package logic

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func (b *broadcaster)ReceiveUser(v string){
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
	err = ch.QueueBind(q.Name, v, "logs_topic", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil, )
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func(){
		for d:= range msgs {
			good := &User{}
			err := json.Unmarshal(d.Body, good)
			if err != nil {
				log.Printf("Error decoding JSON: %s",err)
			}
			//log.Printf("Good: %s",string(d.Body))

			if v == "entering.*.*" {
				b.users[good.Nickname] = good
				OfflineProcessor.Send(good)
			} else {
				delete(b.users, good.Nickname)

			}

		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}