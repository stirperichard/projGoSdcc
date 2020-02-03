# projGoSdcc
**Progetto SDCC in Go**

#How to Installazione
Il progetto è composto da 4 file principali, 2 file utilizzati per il testing e il Makefile.

##I file principali:
queue.go\
utils.go\
producer.go\
consumer.go

##I file di testing:
prova_consumer.go\
prova_producer.go

##Files non di testing:
###Build
Per avviare l’applicativo bisogna come primo passo fare il build dei file queue.go, consumer.go, producer.go, inserendo il file utils.go. \
È stato creato il file Makefile per velocizzare il processo di build e run.\
Utilizzando il compilatore, dopo essersi portati nella directory del progetto, eseguo il build con il comando:

````
make build
````

###Run queue

Dopo aver compilato i 3 files, aprendo altri due terminali (oltre a quello già aperto) ed essersi portati nella directory del progetto, da un terminale lancio il comando:
````
make run queue
````
per mandare in run il server contenente la coda di messaggi. \

###Run producer

Una volta settata la coda, rispondendo alle domande poste sul terminale, da un altro terminale, lancio il comando:
````
make run producer
````

###Run consumer

Ed infine dal terminale rimanente lancio 
````
make run consumer
````

##Per i file di test:
Digitare da terminale i comandi

###Build

Utilizzando il compilatore, dopo essermi portato nella directory del progetto, eseguo il build con il comando:
````
make build
````

###Run queue

Dopo aver compilato i 3 files, aprendo altri due terminali (oltre a quello già aperto) ed essersi portati nella directory del progetto, da un terminale lancio il comando:
````
make run queue
````
per mandare in run il server contenente la coda di messaggi.

###Run producer test
Una volta settata la coda, rispondendo alle domande poste sul terminale, da un altro terminale, lancio il comando:
````
make run producer_test
````
per mandare in run il producer test.

###Run consumer test

Ed infine dal terminale rimanente lancio il comando
````
make run consumer_test
````
per mandare in run il consumer test.

##I parametri:
Per quanto riguarda i files di testing il tempo che intercorre tra l’arrivo del messaggio e l’invio dell’ACK (primo dei due previsti) viene generato in maniera casuale in un intervallo compreso tra 1 e 10 secondi. Stesso vale per il tempo di elaborazione, cioè il tempo che intercorre tra l’invio del primo ACK e il tempo di invio del secondo ACK.

Il producer invia 50 messaggi di tipo stringa di 10 caratteri generati casualmente alla coda. 

I parametri timeout_retransmit, timeout_visibility e semantic (semantica utilizzata) possono essere impostati manualmente sul file utils.go (seguendo le costanti pre-definite) oppure dal terminale della queue all’avvio, dopo aver risposto n alla prima domanda posta chiedente se si vogliano utilizzare i parametri predefiniti o parametri digitati dal terminale.


## Author

* **Richard Stirpe** - *Esercizio SDCC in GO* - [stirperichard](https://github.com/stirperichard)