package main

import (
	"alicloud-demo/example"
	"alicloud-demo/initialize"
)

func init() {
	initialize.Logger()
	initialize.AliCloud()
}

func main() {
	// ECS使用样例
	example.ECS()

	// MNS使用样例
	example.MNS()

	// RDS使用样例
	example.RDS()

	// BILL使用样例
	example.BILL()
}
