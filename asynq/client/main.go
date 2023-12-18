package main

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

// Task payload for any email related tasks.
type EmailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

// client.go
func main() {
	for i := 0; i < 1000; i++ {
		test()
	}
}

func test() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	// Create a task with typename and payload.
	payload, err := json.Marshal(EmailTaskPayload{UserID: 42})
	if err != nil {
		log.Fatal(err)
	}
	t1 := asynq.NewTask("email:welcome", payload)

	t2 := asynq.NewTask("email:reminder", payload)

	// Process the task immediately.
	info, err := client.Enqueue(t1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

	// Process the task 24 hours later.
	info, err = client.Enqueue(t2,
		asynq.ProcessIn(60*time.Second), //相对当前时间多久之后
		//asynq.ProcessAt(), 在某个时间点
		//asynq.Timeout() //任务执行的最大时间，超过则认为是超时重试
		//asynq.Deadline() //任务执行的相对时间，超过则认为是超时，重试
		//asynq.MaxRetry(),//任务重试的最大次数

		//Group返回一个选项，用于指定用于任务的组。
		//给定队列中具有相同组的任务将聚合为一个任务，然后传递给处理程序。
		asynq.Group("testGroup"),

		//asynq.Queue() ////Queue返回一个选项，用于指定将任务排入队列的队列。

		//Retention返回一个选项，用于指定任务的保留期的持续时间。
		//如果提供此选项，则在成功处理后，任务将被存储为已完成的任务。
		//已完成的任务将在指定的持续时间过后删除。
		//asynq.Retention()

		//TaskID返回指定任务ID的选项。
		//asynq.TaskID()

		//asynq.Unique()

	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)
}
