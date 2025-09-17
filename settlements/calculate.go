/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package settlements

import (
	"fmt"
)

// PreCalculateRequest 试算税费请求
type PreCalculateRequest struct {
	ProjectCode string  `json:"projectCode"` // 项目编号
	Name        string  `json:"name"`        // 用工人员姓名
	IDCard      string  `json:"idCard"`      // 用工人员身份证号
	Amount      float64 `json:"amount"`      // 试算金额(元)
}

// PreCalculateResponse 试算税费响应
type PreCalculateResponse struct {
	ProjectCode                  string  `json:"projectCode"`                            // 项目编号
	Name                         string  `json:"name"`                                   // 用工人员姓名
	IDCard                       string  `json:"idCard"`                                 // 用工人员身份证号
	Amount                       float64 `json:"amount"`                                 // 金额(元)
	TotalAmount                  float64 `json:"totalAmount"`                            // 订单总金额(元)
	ReceivedAmount               float64 `json:"receivedAmount"`                         // 到账金额(元)
	CrowdsourcingName            string  `json:"crowdsourcingName"`                      // 众包企业名称
	IncomeTax                    float64 `json:"inComeTax,omitempty"`                    // 预扣个税(元)
	ValueAddedTax                float64 `json:"valueAddedTax,omitempty"`                // 预扣增值税(元)
	AdditionalTax                float64 `json:"additionalTax,omitempty"`                // 预扣附加税(元)
	ServiceCharge                float64 `json:"serviceCharge,omitempty"`                // 应付企业服务费(元)
	ServiceChargeBearWay         int     `json:"serviceChargeBearWay,omitempty"`         // 企业服务费承担方 1:企业 2:个人
	TaxBearWay                   int     `json:"taxBearWay,omitempty"`                   // 税金承担方 1:企业 2:个人
	IncomeTaxRate                float64 `json:"incomeTaxRate,omitempty"`                // 所得税率，小数，0.3000代表30%
	IncomeValueAddedRate         float64 `json:"incomeValueAddedRate,omitempty"`         // 个人增值税税率，小数，0.3000代表30%
	AdditionalTaxRate            float64 `json:"additionalTaxRate,omitempty"`            // 附加税税率，小数，0.3000代表30%
	ServiceChargeRate            float64 `json:"serviceChargeRate,omitempty"`            // 企业服务费费率，小数，0.3000代表30%
	PlatformServiceCharge        float64 `json:"platformServiceCharge,omitempty"`        // 平台服务费(元)
	PlatformServiceChargeRate    float64 `json:"platformServiceChargeRate,omitempty"`    // 平台服务费费率，小数，0.3000代表30%
	PlatformServiceChargeBearWay int     `json:"platformServiceChargeBearWay,omitempty"` // 平台服务费承担方 1:企业 2:个人
}

// PreCalculate 试算税费
func (s *SettlementService) PreCalculate(req *PreCalculateRequest) (*PreCalculateResponse, error) {
	var resp PreCalculateResponse
	err := s.client.DoRequest("/settlement/taxApi/preCalculate", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("试算税费失败: %w", err)
	}

	return &resp, nil
}
