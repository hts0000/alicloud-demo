package main

import (
	example "alicloud-demo/example/ecs"
	"alicloud-demo/initialize"
)

func init() {
	initialize.Logger()
	initialize.AliCloud()
}

func main() {
	// ECS使用样例
	example.ECS()
}
