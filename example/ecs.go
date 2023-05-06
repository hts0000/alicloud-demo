package example

import (
	"alicloud-demo/alicloud"
	aliecs "alicloud-demo/alicloud/ecs"
	"encoding/base64"
	"time"

	"github.com/alibabacloud-go/ecs-20140526/v3/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

func ECS() {
	// 创建ecs客户端
	ecsClient, err := aliecs.NewECSClient(alicloud.AliCloudConfig)
	if err != nil {
		zap.S().Errorf("create ecs client failed, err: %v", err)
		return
	}

	// 获取地区内所有实例
	// document: https://help.aliyun.com/document_detail/25506.html?spm=a2c4g.25498.0.0.6a522752NTkLQY#resultMapping
	descResp, err := ecsClient.DescribeInstancesWithOptions(&client.DescribeInstancesRequest{
		RegionId:   &alicloud.AliCloudConfig.RegionID,
		PageNumber: tea.Int32(1),
		PageSize:   tea.Int32(50),
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("get instances list failed, err: %v", err)
		return
	}
	instanceIDs := []*string{}
	for _, instance := range descResp.Body.Instances.Instance {
		zap.S().Infof("实例序列号: %s, 实例状态: %s, 实例ID: %s, 实例名称: %s, 实例的操作系统名称: %s, ",
			*instance.SerialNumber, *instance.Status, *instance.InstanceId, *instance.InstanceName, *instance.OSName)
		zap.S().Infof("实例创建时间: %s, 过期时间: %s", *instance.CreationTime, *instance.ExpiredTime)

		instanceIDs = append(instanceIDs, instance.InstanceId)
	}

	// 启动地区内所有停止的实例
	// document: https://help.aliyun.com/document_detail/155374.htm?spm=a2c4g.25485.0.0.40a310b424cwQO#doc-api-Ecs-StartInstances
	startResp, err := ecsClient.StartInstancesWithOptions(&client.StartInstancesRequest{
		RegionId:          &alicloud.AliCloudConfig.RegionID,
		InstanceId:        instanceIDs,
		BatchOptimization: tea.String("SuccessFirst"),
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("start stoped instances failed, err: %v", err)
		return
	}
	for _, info := range startResp.Body.InstanceResponses.InstanceResponse {
		zap.S().Infof("<启动实例> - 实例ID: %s, 实例操作结果错误码: %s, 实例操作返回错误信息: %q, 实例当前状态: %s, 操作前实例的状态: %s",
			*info.InstanceId, *info.Code, *info.Message, *info.CurrentStatus, *info.PreviousStatus)
	}

	// 停止地区内所有运行的实例
	// document: https://help.aliyun.com/document_detail/155372.htm?spm=a2c4g.25485.0.0.40a39626Q9TbZt#doc-api-Ecs-StopInstances
	stopResp, err := ecsClient.StopInstancesWithOptions(&client.StopInstancesRequest{
		RegionId:          &alicloud.AliCloudConfig.RegionID,
		InstanceId:        instanceIDs,
		BatchOptimization: tea.String("SuccessFirst"),
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("stop started instances failed, err: %v", err)
		return
	}
	for _, info := range stopResp.Body.InstanceResponses.InstanceResponse {
		zap.S().Infof("<停止实例> - 实例ID: %s, 实例操作结果错误码: %s, 实例操作返回错误信息: %q, 实例当前状态: %s, 操作前实例的状态: %s",
			*info.InstanceId, *info.Code, *info.Message, *info.CurrentStatus, *info.PreviousStatus)
	}

	// 重启地区内所有运行的实例
	// document: https://help.aliyun.com/document_detail/155373.htm?spm=a2c4g.25485.0.0.40a33578m4m8bz#doc-api-Ecs-RebootInstances
	rebootResp, err := ecsClient.RebootInstancesWithOptions(&client.RebootInstancesRequest{
		RegionId:          &alicloud.AliCloudConfig.RegionID,
		ForceReboot:       tea.Bool(false),
		BatchOptimization: tea.String("SuccessFirst"),
		InstanceId:        instanceIDs,
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("reboot started instances failed, err: %v", err)
		return
	}
	for _, info := range rebootResp.Body.InstanceResponses.InstanceResponse {
		zap.S().Infof("<重启实例> - 实例ID: %s, 实例操作结果错误码: %s, 实例操作返回错误信息: %q, 实例当前状态: %s, 操作前实例的状态: %s",
			*info.InstanceId, *info.Code, *info.Message, *info.CurrentStatus, *info.PreviousStatus)
	}

	// 释放过期实例
	// document: https://help.aliyun.com/document_detail/25507.htm?spm=a2c4g.25485.0.0.1d3310b430tcXJ#t9861.html
	releaseResp, err := ecsClient.DeleteInstanceWithOptions(&client.DeleteInstanceRequest{
		InstanceId:            instanceIDs[0],
		Force:                 tea.Bool(false),
		TerminateSubscription: tea.Bool(false), // 是否释放已到期的包年包月实例
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("release expired instances failed, err: %v", err)
		return
	}
	zap.S().Infof("<释放实例> - 实例ID: %s, 实例操作结果错误码: %s", *releaseResp.Body.RequestId, releaseResp.StatusCode)

	// 在指定实例上运行命令，执行命令需要借助云助手：https://help.aliyun.com/document_detail/64601.html?spm=a2c4g.87011.0.0.5342202aO1eeT2
	// 2017年之后使用公共镜像创建的ECS实例默认安装了云助手
	// RunCommand会创建并执行云助手命令
	// document: https://help.aliyun.com/document_detail/141751.html?spm=a2c4g.64841.0.0.cf53739fkNqGes
	command := "echo hello world! 实例id: {{ACS::InstanceId}} 实例名: {{ACS::InstanceName}} 自定义参数: {{diy}}"
	commandBase64 := base64.StdEncoding.EncodeToString([]byte(command))
	zap.S().Infof("<创建云命令> - 命令内容: %q, 命令内容base64加密: %q", command, commandBase64)
	runCmdResp, err := ecsClient.RunCommand(&client.RunCommandRequest{
		RegionId:        &alicloud.AliCloudConfig.RegionID,
		Name:            tea.String("hello cloud command!"),
		Description:     tea.String("打印hello world以及当前实例的信息!"),
		Type:            tea.String("RunShellScript"),
		CommandContent:  &commandBase64, // 命令Base64编码后的内容
		EnableParameter: tea.Bool(true), // 开启可以使用自定义参数和内置参数
		WorkingDir:      tea.String("/root"),
		Timeout:         tea.Int64(60), // 秒
		ContentEncoding: tea.String("Base64"),
		RepeatMode:      tea.String("Once"),
		Parameters: map[string]interface{}{
			"diy": "diydiydiy",
		},
		KeepCommand: tea.Bool(true), // 执行完该命令后，是否保留下来
		Username:    tea.String("root"),
		InstanceId:  instanceIDs,
	})
	if err != nil {
		zap.S().Errorf("create command failed, err: %v", err)
		return
	}
	zap.S().Infof("<运行云命令> - 命令ID: %s, 请求ID: %s, 执行ID: %s", *runCmdResp.Body.CommandId, *runCmdResp.Body.RequestId, *runCmdResp.Body.InvokeId)

	// 查看执行结果，可以查看近4周的执行信息
	invokeCmdResp, err := ecsClient.DescribeInvocationResults(&client.DescribeInvocationResultsRequest{
		RegionId:   &alicloud.AliCloudConfig.RegionID,
		InvokeId:   tea.String("t-sz03hivpti0qpkw"),
		InstanceId: instanceIDs[0],
		CommandId:  tea.String("c-sz03hivpthqr474"),
	})
	if err != nil {
		zap.S().Errorf("invoke command failed, err: %v", err)
		return
	}

	count := 0
	for count < len(invokeCmdResp.Body.Invocation.InvocationResults.InvocationResult) {
		for _, res := range invokeCmdResp.Body.Invocation.InvocationResults.InvocationResult {
			if res.ExitCode == nil {
				continue
			}
			count++
			content := make([]byte, len(*res.Output))
			n, err := base64.StdEncoding.Decode(content, []byte(*res.Output))
			if err != nil {
				zap.S().Errorf("decode command failed, err: %v", err)
				return
			}
			zap.S().Infof("<执行结果> - 实例ID: %s, 命令执行状态: %s, 执行输出结果: %s", *res.InstanceId, *res.InvocationStatus, content[:n])
		}
		time.Sleep(time.Second * 1)
	}

	// 为实例的云盘创建快照
	// 首先查询一个实例下面的云盘ID
	// document: https://help.aliyun.com/document_detail/25514.html?spm=a2c4g.25512.0.0.689b64b9ZwjgW3
	diskResp, err := ecsClient.DescribeDisks(&client.DescribeDisksRequest{
		RegionId:   &alicloud.AliCloudConfig.RegionID,
		InstanceId: instanceIDs[0],
		DiskType:   tea.String("all"),
	})
	if err != nil {
		zap.S().Errorf("get instance disk list failed, err: %v", err)
		return
	}
	for _, disk := range diskResp.Body.Disks.Disk {
		zap.S().Infof("<获取实例云盘> 实例ID: %s, 云盘ID: %s, 云盘名称: %s, 云盘状态: %s, 云盘类型: %s",
			*instanceIDs[0], *disk.DiskId, *disk.DiskName, *disk.Status, *disk.Type)
	}

	// 为云盘创建快照
	// document: https://help.aliyun.com/document_detail/25524.html?spm=a2c4g.25523.0.0.4ceb6117J9j6Cr
	createSpResp, err := ecsClient.CreateSnapshot(&client.CreateSnapshotRequest{
		DiskId:        diskResp.Body.Disks.Disk[0].DiskId,
		SnapshotName:  tea.String("snapshot-create-from-golang-api"),
		Description:   tea.String("使用go api 创建的云盘快照"),
		RetentionDays: tea.Int32(30), // 保存天数
	})
	if err != nil {
		zap.S().Errorf("create disk snapshot failed, err: %v", err)
		return
	}
	zap.S().Infof("<创建快照> - 快照ID: %s, 请求ID: %s", *createSpResp.Body.SnapshotId, *createSpResp.Body.RequestId)
}
