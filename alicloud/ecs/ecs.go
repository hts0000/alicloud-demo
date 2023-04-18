package aliecs

import (
	"alicloud-demo/alicloud"

	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliecs "github.com/alibabacloud-go/ecs-20140526/v3/client"
)

func NewECSClient(config *alicloud.Config) (*aliecs.Client, error) {
	return aliecs.NewClient(&client.Config{
		AccessKeyId:     &config.AccessKeyID,
		AccessKeySecret: &config.AccessKeySecret,
		RegionId:        &config.RegionID,
	})
}
