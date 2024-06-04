package main

type Message struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

func main() {
	msg := Message{"Verification", "HEllo", "sarzhan.yernur@gmail.com"}
	sendEmail(msg)
	//// Connect to RabbitMQ
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")    //if err != nil {
	// log.Fatalf("Failed to connect to RabbitMQ: %v", err)    //}
	//defer conn.Close()    //
	//ch, err := conn.Channel()    //if err != nil {
	// log.Fatalf("Failed to open a channel: %v", err)    //}
	//defer ch.Close()    //
	//q, err := ch.QueueDeclare(    // "notification_queue", // name
	// false,                // durable    // false,                // delete when unused
	// false,                // exclusive    // false,                // no-wait
	// nil,                  // arguments    //)
	//if err != nil {    // log.Fatalf("Failed to declare a queue: %v", err)
	//}    //
	//msgs, err := ch.Consume(    // q.Name, // queue
	// "",     // consumer    // true,   // auto-ack
	// false,  // exclusive    // false,  // no-local
	// false,  // no-wait
	// nil,    // args    //)
	//if err != nil {    // log.Fatalf("Failed to register a consumer: %v", err)
	//}
	// Start the worker goroutine to process messages    //go processMessages(msgs)
	//
	//// Keep the main goroutine running    //select {}
}

////func processMessages(msgs <-chan amqp.Delivery) {
//  for msg := range msgs {//     var message Message
//     err := json.Unmarshal(msg.Body, &message)//     if err != nil {
//        log.Printf("Failed to unmarshal message: %v", err)//        continue
//     }//
//     go sendEmail(message)//     logMessage(message)
//  }//}
////func logMessage(msg Message) {
//  log.Printf("Message Type: %s, Recipient: %s, Message: %s", msg.Type, msg.Recipient, msg.Message)//}
