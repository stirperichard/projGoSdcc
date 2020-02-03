package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func main() {

	producer, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}

	defer producer.Close()

	for i := 0; i < 1000; i++ {

		var m Message
		text1 := StringWithCharset(10, charset)
		text1 = text1 + "\n"

		m.Text = text1
		var reply bool

		msgCall := producer.Go("MessageQueue.PushInQueue", m, &reply, nil)
		msgCall = <-msgCall.Done
		if msgCall.Error != nil {
			log.Fatal("Error in MessageQueue.PushInQueue: ", msgCall.Error.Error())
		}

		fmt.Printf("MessageQueue.PushInQueue: OK \n")
		time.Sleep(20 * time.Millisecond)
	}
}
