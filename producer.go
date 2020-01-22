package projGoSdcc

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
		// Asynchronous call
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text1, _ := reader.ReadString('\n')
		fmt.Println(text1)

		var reply bool
		var reply_s string

		msgCall := producer.Go("Queue.PushInQueue", text1, reply, nil)
		msgCall = <-msgCall.Done
		if msgCall.Error != nil {
			log.Fatal("Error in Queue.PushInQueue: ", msgCall.Error.Error())
		}
		if reply == true {
			reply_s = "OK"
		} else {
			reply_s = "ERROR"
		}
		fmt.Printf("Queue.PushInQueue: %s", reply_s)
	}
}
