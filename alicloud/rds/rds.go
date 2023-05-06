package alirds

import (
	"alicloud-demo/alicloud"

	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alirds "github.com/alibabacloud-go/rds-20140815/v3/client"
)

func NewRDSClient(config *alicloud.Config) (*alirds.Client, error) {
	return alirds.NewClient(&client.Config{
		AccessKeyId:     &config.AccessKeyID,
		AccessKeySecret: &config.AccessKeySecret,
		RegionId:        &config.RegionID,
	})
}
