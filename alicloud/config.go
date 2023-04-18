package alicloud

type Config struct {
	AccessKeyID     string `env:"ACCESS_KEY_ID" require:"true"`
	AccessKeySecret string `env:"ACCESS_KEY_SECRET" require:"true"`
	RegionID        string `env:"REGION_ID" require:"true"`
	ECSEndpoint     string `env:"ECS_ENDPOINT" require:"false"`
	MNSEndpoint     string `env:"MNS_ENDPOINT" require:"true"`
}

var AliCloudConfig *Config
