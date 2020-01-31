package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"time"
)

func main() {

	consumer, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}

	defer consumer.Close()
	for {
		i := scelta()

		if i == 0 {
			//GET INFO
			get_info(consumer)
		} else if i == 1 {
			//CONSUMER
			consume(consumer)
		} else {
			//ERRORE NEL VALORE DIGITATO
			fmt.Println("Digita un valore corretto")
		}
	}
}

func consume(consumer *rpc.Client) {
	for {
		var (
			args, risposta string
			in, out        int
			m              Message
		)

		//Prendo la grandezza della Queue
		msgCallSize := consumer.Go("MessageQueue.GetSize", in, &out, nil)

		msgCallSize = <-msgCallSize.Done
		if msgCallSize.Error != nil {
			log.Fatal("Error in MessageQueue.GetSize: ", msgCallSize.Error.Error())
		}

		fmt.Printf("LEN: %d\n", out)

		//Controllo se la dimensione della queue sia > 1, oppure vado avanti nel caso sia vuota
		if out != 0 {
			msgCall := consumer.Go("MessageQueue.PopFromQueue", args, &m, nil)
			msgCall = <-msgCall.Done
			if msgCall.Error != nil {
				log.Fatal("Error in MessageQueue.PopFromQueue: ", msgCall.Error.Error())
			}

			fmt.Printf("MESSAGGIO RICEVUTO CON ID: %d E CON TESTO: %s", m.ID, m.Text)

			consumeSeconds()

			var ack ACK
			ack.ID = m.ID

			msgACK := consumer.Go("MessageQueue.ReceiveACK", ack, &risposta, nil)
			msgACK = <-msgACK.Done
			if msgACK.Error != nil {
				log.Fatal("Error in MessageQueue.ReceiveACK: ", msgACK.Error.Error())
			}

			elaboraMessaggio()

			msgACK2 := consumer.Go("MessageQueue.ReceiveACK", ack, &risposta, nil)
			msgACK2 = <-msgACK2.Done
			if msgACK2.Error != nil {
				log.Fatal("Error in MessageQueue.ReceiveACK: ", msgACK2.Error.Error())
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func get_info(consumer *rpc.Client) {
	var (
		in int
		m  []Message
	)

	msgCallInfo := consumer.Go("MessageQueue.GetInfoQueue", in, &m, nil)

	msgCallInfo = <-msgCallInfo.Done
	if msgCallInfo.Error != nil {
		log.Fatal("Error in MessageQueue.GetInfoQueue: ", msgCallInfo.Error.Error())
	}

	for index, elem := range m {
		fmt.Printf("MESSAGGIO [" + strconv.Itoa(index) + "]" + "\nID :" + strconv.Itoa(elem.ID) +
			"\nTesto: " + elem.Text + "Status: " + strconv.Itoa(elem.Status) + "\n")
	}
}

func scelta() int {
	fmt.Printf("Cosa vuoi fare? \n Digita 0 per leggere la lista dei messaggi in coda oppure 1 per Consumer \n")
	var i int
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		log.Fatal("Error in Scanf: ", err)
	}
	return i
}

func consumeSeconds() {
	time.Sleep(5 * time.Second)
	return
}

func elaboraMessaggio() {
	time.Sleep(5 * time.Second)
	return
}
