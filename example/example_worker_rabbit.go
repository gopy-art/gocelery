/*
This is an example of celery worker for rabbitmq implementation
in this example we just run 3 worker to listen on 'test' queue for get messages and handle them
*/

package example

import (
	"fmt"

	"github.com/gopy-art/gocelery"
)

func workerRabbit() {
	// declare the url of the rabbitMQ
	rabbit_url := fmt.Sprintf("amqp://%s:%s@%s/", "guest", "guest", "localhost:5672")

	// make the backend and broker with rabbitMQ
	CeleryBackend := gocelery.NewAMQPCeleryBackend(rabbit_url)
	CeleryBroker := gocelery.NewAMQPCeleryBroker(rabbit_url, "test", true)

	// initialize celery client
	worker := gocelery.NewCeleryWorker(
		CeleryBroker,
		CeleryBackend,
		3, // number of worker for the client
	)

	ch := make(chan int)
	// register task
	worker.Register("worker.test", &ExampleAddTask{})

	// start workers (non-blocking call)
	worker.StartWorker()

	// wait for client request
	<-ch

	// stop workers gracefully (blocking call)
	worker.StopWorker()
}
