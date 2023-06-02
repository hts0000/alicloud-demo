package alibill

import (
	"alicloud-demo/alicloud"

	alibill "github.com/alibabacloud-go/bssopenapi-20171214/client"
	"github.com/alibabacloud-go/darabonba-openapi/client"
)

func NewBILLClient(config *alicloud.Config) (*alibill.Client, error) {
	return alibill.NewClient(&client.Config{
		AccessKeyId:     &config.AccessKeyID,
		AccessKeySecret: &config.AccessKeySecret,
		RegionId:        &config.RegionID,
	})
}
