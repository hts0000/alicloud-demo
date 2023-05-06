package example

import (
	"alicloud-demo/alicloud"
	alirds "alicloud-demo/alicloud/rds"
	"time"

	"github.com/alibabacloud-go/rds-20140815/v3/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

func RDS() {
	// 创建rds客户端
	rdsClient, err := alirds.NewRDSClient(alicloud.AliCloudConfig)
	if err != nil {
		zap.S().Errorf("create rds client failed, err: %v", err)
		return
	}

	// 获取区域内所有实例
	// document: https://help.aliyun.com/document_detail/610396.htm?spm=a2c4g.610369.0.0.235a5813ELdlzz
	descResp, err := rdsClient.DescribeDBInstancesWithOptions(&client.DescribeDBInstancesRequest{
		Engine:           tea.String("MYSQL"),       // 数据库类型
		ZoneId:           tea.String("cn-shenzhen"), // 可用区ID
		ResourceGroupId:  tea.String(""),            // 资源组ID
		DBInstanceStatus: tea.String("Running"),     // 实例状态
		SearchKey:        tea.String(""),            //可基于实例ID或者实例备注模糊搜索
		DBInstanceId:     tea.String(""),            // 实例ID
		DBInstanceType:   tea.String("Primary"),     // 实例类型
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("describe rds instances failed, err: %v", err)
		return
	}
	for _, instance := range descResp.Body.Items.DBInstance {
		zap.S().Infof("数据库id: %s, 数据库类型: %s, 数据库描述: %s, 数据库状态: %s",
			*instance.DBInstanceId, *instance.Engine, *instance.DBInstanceDescription, *instance.DBInstanceStatus)
	}

	// 查询只读实例复制延迟
	// document: https://help.aliyun.com/document_detail/610483.htm?spm=a2c4g.610369.0.0.c25d5d7aSqVm8U
	delayResp, err := rdsClient.DescribeReadDBInstanceDelayWithOptions(&client.DescribeReadDBInstanceDelayRequest{
		DBInstanceId:   tea.String("xxxxxxx"), // 主实例ID
		ReadInstanceId: tea.String("xxxxxx"),  // 只读实例ID
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("describe rds read only instances delay failed, err: %v", err)
		return
	}
	zap.S().Infof("复制延迟: %d", delayResp.Body.DelayTime)

	// 查询实例的SQL审计日志
	// document: https://help.aliyun.com/document_detail/610533.html?spm=a2c4g.610530.0.0.437da906OvEguB
	recordResp, err := rdsClient.DescribeSQLLogRecordsWithOptions(&client.DescribeSQLLogRecordsRequest{
		DBInstanceId: tea.String(""), // 实例id
		StartTime: tea.String(time.Now().AddDate(0, 0, -14).
			UTC().Format("2006-01-02T15:04:05Z")), // 查询开始时间，可查询当前日期前15天内的数据。格式：yyyy-MM-ddTHH:mm:ssZ（UTC时间）
		EndTime: tea.String(time.Now().Format("2006-01-02T15:04:05Z")), // 查询结束时间
		Form:    tea.String("Stream"),                                  // 触发审计文件的生成或者返回SQL记录列表
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("describe sql record failed, err: %v", err)
		return
	}
	for _, record := range recordResp.Body.Items.SQLRecord {
		zap.S().Infof("sql语句: %s, 执行sql的耗时(微秒): %f, 执行sql的ip: %s, 执行操作的账户名: %s",
			*record.SQLText, *record.TotalExecutionTimes, *record.HostAddress, *record.AccountName)
	}

	// 查询实例的慢SQL审计日志
	slowLogResp, err := rdsClient.DescribeSlowLogsWithOptions(&client.DescribeSlowLogsRequest{
		DBInstanceId: tea.String(""), // 实例ID
		StartTime: tea.String(time.Now().AddDate(0, 0, -14).
			UTC().Format("2006-01-02T15:04:05Z")), // 查询开始时间，可查询当前日期前15天内的数据。格式：yyyy-MM-ddTHH:mm:ssZ（UTC时间）
		EndTime: tea.String(time.Now().Format("2006-01-02T15:04:05Z")), // 查询结束时间
	}, &service.RuntimeOptions{
		ReadTimeout: tea.Int(1000),
		MaxAttempts: tea.Int(3),
	})
	if err != nil {
		zap.S().Errorf("describe slow sql failed, err: %v", err)
		return
	}
	for _, record := range slowLogResp.Body.Items.SQLSlowLog {
		zap.S().Infof("sql语句: %s, 执行sql平均的耗时(秒): %d, 慢sql唯一标识: %s, 数据库名称: %s",
			*record.SQLText, *record.AvgExecutionTime, *record.SQLHASH, *record.DBName)
	}
}
