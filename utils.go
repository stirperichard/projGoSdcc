package main

import "sync"

//Definizione costanti
const ATLEASTONCE = 0
const TIMEOUTBASED = 1

const INQUEUE = 0
const WAITINGACK = 1
const SENDED = 2
const ELABORATED = 3

const TIMEOUT_VISIBILITY = 10
const TIMEOUT_RETRANSMIT = 10
const SEMANTIC = TIMEOUTBASED //è possibile modificarlo

//Definizione variabili
var MessageID int

// Definizione delle strutture
type Message struct {
	ID                 int
	Text               string
	Timeout_retransmit int
	Timeout_visibility int
	Semantic           int // Con (0) indidico At-Least-Once, con (1) indico Timeout-Based
	Status             int //Indico con (0) se in coda; (1) se attende ACK; (2) timeout-based; (3) ACK finale ricevuto;
	// (4) Ricevuto 1° ACK timeout-based (ricezione). L'elaborazione la consider con (3)
}

type MessageQueue struct {
	Messages []Message
	Lock     sync.RWMutex
}

type ACK struct {
	ID int
}
