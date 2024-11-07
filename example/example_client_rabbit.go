/*
This is an example of celery controller for rabbitmq
in this example we just send a normal message in 'test' queue
*/

package example

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func clientRabbit() {
	// declare the url of the rabbitMQ
	rabbit_url := fmt.Sprintf("amqp://%s:%s@%s/", "guest", "guest", "localhost:5672")

	// make the backend and broker with rabbitMQ
	CeleryBackend := gocelery.NewAMQPCeleryBackend(rabbit_url)
	CeleryBroker := gocelery.NewAMQPCeleryBroker(rabbit_url, "test", true)

	// initialize celery client
	client, err := gocelery.NewCeleryClient(
		CeleryBroker,
		CeleryBackend,
		1, // number of worker for the client
	)
	if err != nil {
		log.Println("error in client , ", err)
	}
	////// Error handling

	// declare the task name for the message
	taskName := fmt.Sprintf("worker.%s", "test")

	for v := range 3 {
		_ = v
		// prepare the message for set to the queue
		msg := make(map[string]interface{})
		msg["a"] = rand.Intn(10)
		msg["b"] = rand.Intn(10)

		// send the message to the rabbitMQ
		_, err = client.DelayKwargs(taskName, msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	client.WaitForStopWorker()
}
