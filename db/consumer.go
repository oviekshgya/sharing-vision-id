package db

import (
	"fmt"
	"log"
)

func StartConsumerPayment() {
	msgs, err := RabbitChannel.Consume(
		"payment", // Queue Name
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("failed RabbitMQ:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Println("Received a message:", d.Body)
			//userService := service.UserService{
			//	DB: DBMongo,
			//}
			//var transaksi pkg.JSONRequestReadPayment
			//err := json.Unmarshal(d.Body, &transaksi)
			//if err != nil {
			//	return
			//}
			//result, errRes := userService.UpdatePayment(transaksi)
			//if errRes != nil {
			//	log.Printf("Failed user: %v", errRes.Error())
			//}
		}
	}()

	fmt.Println("RabbitMQ Consumer Payment runnig...")
	<-forever
}
