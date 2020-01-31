GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run

# File names
UTILS_FILE=utils.go
PRODUCER_FILE = producer.go
CONSUMER_FILE = consumer.go
QUEUE_FILE = queue.go
PRODUCER_TEST_FILE = prova_producer.go
CONSUMER_TEST_FILE = prova_consumer.go



build:
	$(GOBUILD) $(QUEUE_FILE) $(UTILS_FILE)
	$(GOBUILD) $(PRODUCER_FILE) $(UTILS_FILE)
	$(GOBUILD) $(CONSUMER_FILE) $(UTILS_FILE)

run queue:
	$(GORUN) $(QUEUE_FILE) $(UTILS_FILE)

run producer:
	$(GORUN) $(PRODUCER_FILE) $(UTILS_FILE)

run consumer:
	$(GORUN) $(CONSUMER_FILE) $(UTILS_FILE)

run producer_test:
	$(GORUN) $(PRODUCER_TEST_FILE) $(UTILS_FILE)

run consumer_test:
	$(GORUN) $(CONSUMER_TEST_FILE) $(UTILS_FILE)