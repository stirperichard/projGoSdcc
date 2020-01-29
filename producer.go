package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {

	producer, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}

	defer producer.Close()

	for {

		fmt.Print("Enter text: ")
		text1, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		var reply bool

		msgCall := producer.Go("MessageQueue.PushInQueue", text1, &reply, nil)
		msgCall = <-msgCall.Done
		if msgCall.Error != nil {
			log.Fatal("Error in MessageQueue.PushInQueue: ", msgCall.Error.Error())
		}

		fmt.Printf("MessageQueue.PushInQueue: OK \n")
	}
}
