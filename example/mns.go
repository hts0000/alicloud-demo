package example

import (
	"fmt"
	"strings"
	"time"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
)

func MNS() {
	// 消息服务，支持队列模型(一对一)和订阅模型(一对多)
	// quereManager := alimns.NewMNSQueueClient(alicloud.AliCloudConfig)

	// 列出现有队列
	// quereManager.ListQueue()

	// 创建队列

	// 往队列发送消息

	// 从队列接受消息

	// 删除队列

	// 创建主题

	// 创建订阅

	// 往主题中发送消息

	// 订阅主题接受消息

	// 删除主题

	// 系统事件监听和触发：https://help.aliyun.com/document_detail/118597.htm?spm=a2c4g.116341.0.0.66f170952MG9Np#task-2458106
	// 首先需要创建消息队列：https://mns.console.aliyun.com/region/cn-shenzhen/queues
	// 然后设置系统事件发送到创建的消息队列中：https://cloudmonitor.console.aliyun.com/system-events
	// mnsClient := mns.NewMNSClient(alicloud.AliCloudConfig)
}

func Main() {
	// 设置HTTP接入域名。
	endpoint := "${HTTP_ENDPOINT}"
	// AccessKey ID，阿里云身份验证标识。获取方式，请参见本文前提条件中的获取AccessKey。
	accessKey := "${ACCESS_KEY}"
	// AccessKey Secret，阿里云身份验证密钥。获取方式，请参见本文前提条件中的获取AccessKey。
	secretKey := "${SECRET_KEY}"
	// 所属的Topic。
	topic := "${TOPIC}"
	// Topic所属实例ID，默认实例为空。
	instanceId := "${INSTANCE_ID}"
	// 您在控制台创建的Group ID（Consumer ID）。
	groupId := "${GROUP_ID}"

	client := mq_http_sdk.NewAliyunMQClient(endpoint, accessKey, secretKey, "")

	mqConsumer := client.GetConsumer(instanceId, topic, groupId, "")

	for {
		endChan := make(chan int)
		respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
		errChan := make(chan error)
		go func() {
			select {
			case resp := <-respChan:
				{
					// 处理业务逻辑。
					var handles []string
					fmt.Printf("Consume %d messages---->\n", len(resp.Messages))
					for _, v := range resp.Messages {
						handles = append(handles, v.ReceiptHandle)
						fmt.Printf("\tMessageID: %s, PublishTime: %d, MessageTag: %s\n"+
							"\tConsumedTimes: %d, FirstConsumeTime: %d, NextConsumeTime: %d\n"+
							"\tBody: %s\n"+
							"\tProps: %s\n",
							v.MessageId, v.PublishTime, v.MessageTag, v.ConsumedTimes,
							v.FirstConsumeTime, v.NextConsumeTime, v.MessageBody, v.Properties)
					}

					// NextConsumeTime前若不确认消息消费成功，则消息会重复消费。
					// 消息句柄有时间戳，同一条消息每次消费拿到的都不一样。
					ackerr := mqConsumer.AckMessage(handles)
					if ackerr != nil {
						// 某些消息的句柄可能超时了会导致确认不成功。
						fmt.Println(ackerr)
						for _, errAckItem := range ackerr.(errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem) {
							fmt.Printf("\tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
								errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
						}
						time.Sleep(time.Duration(3) * time.Second)
					} else {
						fmt.Printf("Ack ---->\n\t%s\n", handles)
					}

					endChan <- 1
				}
			case err := <-errChan:
				{
					// 没有消息。
					if strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
						fmt.Println("\nNo new message, continue!")
					} else {
						fmt.Println(err)
						time.Sleep(time.Duration(3) * time.Second)
					}
					endChan <- 1
				}
			case <-time.After(35 * time.Second):
				{
					fmt.Println("Timeout of consumer message ??")
					endChan <- 1
				}
			}
		}()

		// 长轮询消费消息。
		// 长轮询表示如果Topic没有消息则请求会在服务端挂住3s，3s内如果有消息可以消费则立即返回。
		mqConsumer.ConsumeMessage(respChan, errChan,
			3, //  一次最多消费3条（最多可设置为16条）。
			3, // 长轮询时间3秒（最多可设置为30秒）。
		)
		<-endChan
	}
}
