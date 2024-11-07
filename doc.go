/*
# gocelery
Gocelery is a task queue implementation for Go modules used to asynchronously execute work outside the HTTP request-response cycle. 
Celery is an implementation of the task queue concept.

# How it works?
this package has two side for does its work! <br>
1 - controller <br>
2 - worker

`controller` : it implements functions which work with broker and backend , and the responsebilty of them are to declare and insert data to the broker!

`worker` : it implements functions which work with broker and backend, and the responsebilty of them are to read from broker and set result to the backend!

Note: we have two concept in this package, broker and backend. <strong>broker</strong> is a system that we can share and read our data from that. <strong>backend</strong> is a system that we can set our results from package in that, for example the result of our workers!

# Distributed Systems you can use
1. RabbitMQ = RabbitMQ is an open source message agent software that implements the Advanced Message Queuing protocol.

The files that impelement rabbitMQ celery in this package are : 
    * amqp_backend.go
    * amqp_broker.go
    * amqp.go

2. Redis = Redis is a very high-performance extensible key-value database management system written in C ANSI. It is part of the NoSQL movement and aims to provide the highest possible performance.

The files that impelement redis celery in this package are : 
    * redis_backend.go
    * redis_broker.go

# How to use it?

if you want to use this package you can follow this steps :
   * clone this module 
   * make `gocelery` folder in your project
   * copy and paste all of the files to the gocelery folder
   * run `go mod tidy` command in your terminal
now enjoy!
*/

package gocelery