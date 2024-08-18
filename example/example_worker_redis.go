package example

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func workerRedis() {
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
}
