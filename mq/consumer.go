package mq

import "fmt"

var done chan bool

// StartConsume 开始监听队列，获取消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1. 通过channel.Consume获得消息通道
	msgs, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil {

		fmt.Println(err.Error())
		return
	}
	done = make(chan bool)
	go func() {
		// 2. 循环从channel获取新的消息
		for msg := range msgs {
			// 3. 每次获取新的消息都会调用callback
			procssSuc := callback(msg.Body)
			if procssSuc {
				// 1. TODO 将任务写到另外一个队列，用来异常重试
			}
		}
	}()

	<-done
	channel.Close()
}
