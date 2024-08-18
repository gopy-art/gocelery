package example

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

func clientRedis() {
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
}