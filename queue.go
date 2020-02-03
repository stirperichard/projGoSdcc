package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

//Definizione delle variabili
var (
	MsgQ              MessageQueue   //Message Queue utilizzata
	mutex             = sync.Mutex{} //Mutex
	id                int
	timeoutRetransmit int
	timeoutVisibility int
	semantic          int
	clientConnected   int
)

//Funzioni
func scelta_variabili(i int) {
	fmt.Print("Vuoi usare le variabili predefinite? y = YES, n = NO ")
	text1, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if text1 == "y\n" {
		timeoutRetransmit = TIMEOUT_RETRANSMIT
		timeoutVisibility = TIMEOUT_VISIBILITY
		semantic = SEMANTIC
	} else if text1 == "n\n" {
		fmt.Printf(" Digita il valore del Timeout retransmit? \n")
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal("Errore Scanf - Inserire un valore corretto.", err)
		}

		timeoutRetransmit = i

		fmt.Printf(" Digita il valore del Timeout visibility? \n")
		_, err = fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal("Errore Scanf - Inserire un valore corretto.", err)
		}

		timeoutVisibility = i

		fmt.Printf(" Digita il valore della semantic delivery\n 0 per at-least-once, 1 per timeout-based\n")
		_, err = fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal("Errore Scanf - Inserire un valore corretto.", err)
		}

		semantic = i
	} else {
		fmt.Println("Digitare un valore corretto")
	}
	fmt.Println("QUEUE CREATA IN ASCOLTO")
}

func (MsgQ *MessageQueue) GetSize(i int, ret *int) error {
	*ret = len(MsgQ.Messages)
	return nil
}

func (MsgQ *MessageQueue) newMessageQueue() *MessageQueue {
	MsgQ.Messages = []Message{}
	fmt.Print("Message Queue Creata \n")
	return MsgQ
}

func (MsgQ *MessageQueue) initMessageQueue() *MessageQueue {
	if MsgQ == nil {
		MsgQ.newMessageQueue()
	}
	fmt.Print("Message Queue Inizializzata \n")
	return MsgQ
}

func (MsgQ *MessageQueue) ClientConnected(a string, i *int) error {
	*i = clientConnected
	return nil
}

func (MsgQ *MessageQueue) push(m Message) {
	MsgQ.Messages = append(MsgQ.Messages, m)
}

func (MsgQ *MessageQueue) PushInQueue(m Message, reply *bool) error {
	m.Timeout_retransmit = timeoutRetransmit
	m.Timeout_visibility = timeoutVisibility
	m.Semantic = semantic
	m.Status = INQUEUE
	mutex.Lock()
	m.ID = id
	id++
	MsgQ.push(m)
	mutex.Unlock()
	fmt.Printf("Messaggio ricevuto: %s", m.Text)
	return nil
}

func (MsgQ *MessageQueue) pop() Message {
	var (
		value int
	)
	if len(MsgQ.Messages) != 0 {
		mutex.Lock()
		for i := 0; i < len(MsgQ.Messages); i++ {
			if MsgQ.Messages[i].Status == INQUEUE {
				value = i
				break
			}
		}
		MsgQ.Messages[value].Status = WAITINGACK //Attendo ACK+
		mutex.Unlock()
	}
	return MsgQ.Messages[value]
}

func (MsgQ *MessageQueue) PopFromQueue(a string, s *Message) error {
	if len(MsgQ.Messages) > 0 {
		var m = MsgQ.pop()
		*s = m
		go routine(m, MsgQ)
	}
	return nil
}

func routine(m Message, msg_q *MessageQueue) {
	time.Sleep(time.Duration(timeoutRetransmit) * time.Second)
	mutex.Lock()
	if len(msg_q.Messages) < 1 {
		mutex.Unlock()
		fmt.Println("MESSAGE QUEUE VUOTA")
		return
	}

	//Trovo il messaggio
	for i := 0; i < len(msg_q.Messages); i++ {
		if msg_q.Messages[i].ID == m.ID {
			//SEMANTICA At Least Once
			if semantic == ATLEASTONCE {
				fmt.Printf("MESSAGE ID: %d    STATUS: %d\n", msg_q.Messages[i].ID, msg_q.Messages[i].Status)
				if msg_q.Messages[i].Status == WAITINGACK || msg_q.Messages[i].Status == SENDED {
					msg_q.Messages[i].Status = INQUEUE
					fmt.Println("MESSAGGIO REINSERITO IN CODA - timeout retransmit scaduto")
				} else if msg_q.Messages[i].Status == ELABORATED || msg_q.Messages[i].Status == INQUEUE {
				}
			} else if semantic == TIMEOUTBASED { //SEMANTICA TIMEOUT BASED
				if msg_q.Messages[i].Status == WAITINGACK {
					msg_q.Messages[i].Status = INQUEUE
					fmt.Println("MESSAGGIO REINSERITO IN CODA")
				} else if msg_q.Messages[i].Status == SENDED {
					time.Sleep(time.Duration(timeoutVisibility) * time.Second)
					if msg_q.Messages[i].Status == SENDED {
						msg_q.Messages[i].Status = INQUEUE
						fmt.Println("MESSAGGIO REINSERITO IN CODA")
					}
				} else if msg_q.Messages[i].Status == ELABORATED || msg_q.Messages[i].Status == INQUEUE {
					fmt.Println("OK")
				}
			}
		}
	}
	mutex.Unlock()
	return
}

func (MsgQ *MessageQueue) ReceiveACK(a ACK, s *string) error {

	id = a.ID

	if semantic == ATLEASTONCE || semantic == TIMEOUTBASED {
		mutex.Lock()
		for i := 0; i < len(MsgQ.Messages); i++ {
			//CONTROLLO L'ID DEL MESSAGGIO DI CUI HO RICEVUTO ACK
			if MsgQ.Messages[i].ID == id {
				//SE IL SUO STATUS è INQUEUE e ho ricevuto ACK status->SENDED
				fmt.Printf("Message ID: %d with status: %d - RECEIVED ACK\n", MsgQ.Messages[i].ID, MsgQ.Messages[i].Status)
				if MsgQ.Messages[i].Status == WAITINGACK {
					fmt.Printf("FIRST ACK RECEIVED FOR MESSAGE FOR ID: %d\n", id)
					MsgQ.Messages[i].Status = SENDED
					mutex.Unlock()
					continue
				} else if MsgQ.Messages[i].Status == SENDED {
					//SE IL SUO STATUS è SENDED e ho ricevuto ACK status->ELABORATED
					fmt.Printf("SECOND ACK RECEIVED FOR MESSAGE FOR ID: %d\n", id)
					MsgQ.Messages[i].Status = ELABORATED
					mutex.Unlock()
					continue
				} else if MsgQ.Messages[i].Status == INQUEUE {
					//Scarta l'ACK
					fmt.Println("ACK SCARTATO")
					mutex.Unlock()
					continue
				} else {
					mutex.Unlock()
				}
			}
		}
		//FACCIO SCORRERE LA MESSAGE QUEUE, ELIMINANDO I MESSAGGI GIA ELABORATI
		for i := 0; i < len(MsgQ.Messages); i++ {
			mutex.Lock()
			if MsgQ.Messages[i].Status == ELABORATED {
				fmt.Printf("MESSAGE REMOVED FROM QUEUE WITH ID: %d\n", MsgQ.Messages[i].ID)
				MsgQ.Messages = MsgQ.Messages[1:len(MsgQ.Messages)]
				mutex.Unlock()
			} else {
				mutex.Unlock()
				continue
			}
		}
	}
	return nil
}

func (MsgQ *MessageQueue) GetInfoQueue(b int, e *[]Message) error {
	mutex.Lock()
	*e = MsgQ.Messages
	fmt.Println("RICHIESTA GET INFO DA CONSUMER")
	mutex.Unlock()
	return nil
}

func main() {
	var i int
	queue := new(MessageQueue)
	server := rpc.NewServer()
	err := server.RegisterName("MessageQueue", queue)
	if err != nil {
		log.Fatal("Format of service Queue is not correct: ", err)
	}

	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":12345")
	if e != nil {
		log.Fatal("Listen error:", e)
	}

	MsgQ.initMessageQueue()

	id = 0

	defer l.Close()

	scelta_variabili(i)

	clientConnected = 0
	for {
		clientConnected++
		server.Accept(l)
	}
}
