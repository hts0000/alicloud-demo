package initialize

import (
	"alicloud-demo/alicloud"

	"github.com/codingconcepts/env"
	"go.uber.org/zap"
)

func AliCloud() {
	alicloud.AliCloudConfig = &alicloud.Config{}
	if err := env.Set(alicloud.AliCloudConfig); err != nil {
		zap.S().Fatalf("get config from env failed, err: %v", err)
		return
	}
	zap.S().Infof("get config from env success, config: %#v", alicloud.AliCloudConfig)
}
