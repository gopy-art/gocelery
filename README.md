[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/go-1.24.1-blue)](https://go.dev/dl/) 
[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/wotk_with-rabbitMQ-yellow)](https://go.dev/dl/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/work_with-redis-red)](https://go.dev/dl/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/unit_test-failing-green)](https://go.dev/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopy-art/gocelery)](https://goreportcard.com/report/github.com/gopy-art/gocelery)
[![Go Reference](https://pkg.go.dev/badge/github.com/rabbitmq/amqp091-go.svg)](https://pkg.go.dev/github.com/gopy-art/gocelery)

# gocelery
Gocelery is a task queue implementation for Go modules used to asynchronously execute work outside the HTTP request-response cycle. Celery is an implementation of the task queue concept.

- [gocelery](#gocelery)
	- [How it works?](#how-it-works)
	- [Distributed Systems you can use](#distributed-systems-you-can-use)
	- [How to use it?](#how-to-use-it)
	- [Implement controller code example](#implement-controller-code-example)
	- [Implement worker code example](#implement-worker-code-example)
	- [Structure of messages in broker](#structure-of-messages-in-broker)

## How it works?
this package has two side for does its work! <br>
1 - controller <br>
2 - worker

`controller` : it implements functions which work with broker and backend , and the responsebilty of them are to declare and insert data to the broker!

`worker` : it implements functions which work with broker and backend, and the responsebilty of them are to read from broker and set result to the backend!

<strong>Note</strong> : we have two concept in this package, broker and backend. <strong>broker</strong> is a system that we can share and read our data from that. <strong>backend</strong> is a system that we can set our results from package in that, for example the result of our workers!

## Distributed Systems you can use
1. `RabbitMQ` = RabbitMQ is an open source message agent software that implements the Advanced Message Queuing protocol.<br>
The files that impelement rabbitMQ celery in this package are : 
    * amqp_backend.go
    * amqp_broker.go
    * amqp.go

2. `Redis` = Redis is a very high-performance extensible key-value database management system written in C ANSI. It is part of the NoSQL movement and aims to provide the highest possible performance. <br>
The files that impelement redis celery in this package are : 
    * redis_backend.go
    * redis_broker.go

## How to use it?
if you want to use this package you can follow this steps :
   * clone this module 
   * make `gocelery` folder in your project
   * copy and paste all of the files to the gocelery folder
   * run `go mod tidy` command in your terminal
now enjoy!

## Implement controller code example
you can write controller for redis and rabbitMQ, and as you read at the top, the controller will set and declare the messages and data to the broker and backend. <br>
code example with <strong>RabbitMQ</strong> : 
```go
	// declare the url of the rabbitMQ
	rabbit_url := fmt.Sprintf("amqp://%s:%s@%s/", username, password, url)

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
```

code example with <strong>Redis</strong> : 
```go
    // create redis connection pool
	redisPool := &redis.Pool{
		MaxIdle:     3,                 // maximum number of idle connections in the pool
		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// initialize celery client
	client, err := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1, // number of workers
	)
	if err != nil {
		log.Println("error in client , ", err)
	}
	////// Error handling

	// declare the task name for the message
	taskName := fmt.Sprintf("worker.%s", "test")

	for v := range 10 {
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
```

## Implement worker code example
you can write worker for redis and rabbitMQ too, and as you know worker will read from redis or rabbitMQ and set the result to the redis or rabbitMQ.<br>
code example with <strong>RabbitMQ</strong> :
```go
	// exampleAddTask is integer addition task
	// with named arguments
	type ExampleAddTask struct {
		TaskID string
		a      int
		b      int
	}

	// this function is for reading the argument that has been passed to the message from 'Kwargs'
	func (a *ExampleAddTask) ParseKwargs(kwargs map[string]interface{}) error {
		kwargA, ok := kwargs["a"]
		if !ok {
			return fmt.Errorf("undefined kwarg a")
		}
		kwargAFloat, ok := kwargA.(float64)
		if !ok {
			return fmt.Errorf("malformed kwarg a")
		}
		a.a = int(kwargAFloat)
		kwargB, ok := kwargs["b"]
		if !ok {
			return fmt.Errorf("undefined kwarg b")
		}
		kwargBFloat, ok := kwargB.(float64)
		if !ok {
			return fmt.Errorf("malformed kwarg b")
		}
		a.b = int(kwargBFloat)
		return nil
	}

	func (a *ExampleAddTask) ParseId(id string) error {
		a.TaskID = id
		return nil
	}

	// The main function that will be execute
	func (a *ExampleAddTask) RunTask() (interface{}, error) {
		result := a.a + a.b
		fmt.Printf("Task with uuid %v has result %v \n", a.TaskID, result)
		return result, nil
	}

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
```

code example with <strong>Redis</strong> :
```go
	// exampleAddTask is integer addition task
	// with named arguments
	type ExampleAddTask struct {
		TaskID string
		a      int
		b      int
	}

	// this function is for reading the argument that has been passed to the message from 'Kwargs'
	func (a *ExampleAddTask) ParseKwargs(kwargs map[string]interface{}) error {
		kwargA, ok := kwargs["a"]
		if !ok {
			return fmt.Errorf("undefined kwarg a")
		}
		kwargAFloat, ok := kwargA.(float64)
		if !ok {
			return fmt.Errorf("malformed kwarg a")
		}
		a.a = int(kwargAFloat)
		kwargB, ok := kwargs["b"]
		if !ok {
			return fmt.Errorf("undefined kwarg b")
		}
		kwargBFloat, ok := kwargB.(float64)
		if !ok {
			return fmt.Errorf("malformed kwarg b")
		}
		a.b = int(kwargBFloat)
		return nil
	}

	func (a *ExampleAddTask) ParseId(id string) error {
		a.TaskID = id
		return nil
	}

	// The main function that will be execute
	func (a *ExampleAddTask) RunTask() (interface{}, error) {
		result := a.a + a.b
		fmt.Printf("Task with uuid %v has result %v \n", a.TaskID, result)
		return result, nil
	}

	// create redis connection pool
	redisPool := &redis.Pool{
		MaxIdle:     3,                 // maximum number of idle connections in the pool
		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// initialize celery client
	worker := gocelery.NewCeleryWorker(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		6, // number of workers
	)

	////// Error handling

	ch := make(chan int)
	// register task
	worker.Register("worker.test", &ExampleAddTask{})

	// start workers (non-blocking call)
	worker.StartWorker()

	// wait for client request
	<-ch

	// stop workers gracefully (blocking call)
	worker.StopWorker()
```

<strong>NOTE</strong> : you can also use both distributed systems (redis, rabbitMQ) at the same time.
for example you can use rabbitMQ for broker and redis for backend!

## Structure of messages in broker
the message type in broker is json and its in this structure : 
```go
	// TaskMessage is celery-compatible message
	type TaskMessage struct {
		ID      string                 `json:"id"`
		Task    string                 `json:"task"`
		Args    []interface{}          `json:"args"`
		Kwargs  map[string]interface{} `json:"kwargs"`
		Retries int                    `json:"retries"`
		ETA     *string                `json:"eta"`
		Expires *time.Time             `json:"expires"`
	}
```

<strong>Note</strong> : Celery must be configured to use json instead of default pickle encoding. This is because Go currently has no stable support for decoding pickle objects. Pass below configuration parameters to use json.

I hope you can enjoy using this package!