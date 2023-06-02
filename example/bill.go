package example

import (
	"alicloud-demo/alicloud"
	alibill "alicloud-demo/alicloud/bill"

	"github.com/alibabacloud-go/bssopenapi-20171214/client"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

func BILL() {
	billClient, err := alibill.NewBILLClient(alicloud.AliCloudConfig)
	if err != nil {
		zap.S().Errorf("create bill client failed, err: %v", err)
		return
	}

	// 查询用户账户余额信息
	// document: https://help.aliyun.com/document_detail/472991.html#doc-api-BssOpenApi-QueryAccountBalance
	balanceResp, err := billClient.QueryAccountBalance()
	if err != nil {
		zap.S().Errorf("get balance failed, err: %v", err)
		return
	}
	balance := *balanceResp.Body.Data
	zap.S().Infof("可用额度: %v, 信控余额: %v, 网商银行信用额度: %v, 币种(人民币/美元/日元): %v, 现金余额: %v",
		*balance.AvailableAmount, *balance.CreditAmount, *balance.MybankCreditAmount, *balance.Currency, *balance.AvailableCashAmount)

	// 查询用户某个账期内账单总览信息
	// document: https://help.aliyun.com/document_detail/472986.html#doc-api-BssOpenApi-QueryBillOverview
	billOverviewResp, err := billClient.QueryBillOverview(&client.QueryBillOverviewRequest{
		BillingCycle: tea.String("2023-01"),
	})
	if err != nil {
		zap.S().Errorf("query bill overview failed, err: %v", err)
		return
	}
	for _, bill := range billOverviewResp.Body.Data.Items.Item {
		zap.S().Infof("应付金额: %.2f, 产品名称: %v, 原始金额: %v", *bill.PretaxAmount, *bill.ProductName, *bill.PretaxGrossAmount)
	}

	// 查询分账账单
	// document: https://help.aliyun.com/document_detail/473031.html?spm=a2c4g.473049.0.0.9f4242ebmylNmP
	describeResp, err := billClient.DescribeSplitItemBill(&client.DescribeSplitItemBillRequest{
		BillingCycle: tea.String("2023-03"),
		ProductType:  tea.String("ecs"),
	})
	if err != nil {
		zap.S().Errorf("describe split item bill failed, err: %v", err)
		return
	}
	for _, bill := range describeResp.Body.Data.Items {
		zap.S().Infof("实例昵称: %v, 产品明细: %v, 用量: %v%s, 原始金额: %.2f, 应付金额: %.2f",
			bill.NickName, bill.ProductDetail, bill.Usage, bill.UsageUnit, bill.PretaxGrossAmount, bill.PretaxAmount)
	}
}
