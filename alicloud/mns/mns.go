package alimns

import (
	"alicloud-demo/alicloud"

	alimns "github.com/aliyun/aliyun-mns-go-sdk"
)

func NewMNSClient(config *alicloud.Config) alimns.MNSClient {
	return alimns.NewAliMNSClientWithConfig(alimns.AliMNSClientConfig{
		EndPoint:        config.MNSEndpoint,
		AccessKeyId:     config.AccessKeyID,
		AccessKeySecret: config.AccessKeySecret,
	})
}

func NewMNSQueueClient(config *alicloud.Config) alimns.AliQueueManager {
	return alimns.NewMNSQueueManager(NewMNSClient(config))
}
