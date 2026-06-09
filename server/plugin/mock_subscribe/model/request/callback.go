package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type CallbackRecordSearch struct {
	commonReq.PageInfo
	StartCreatedAt  *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt    *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	MchID           string     `json:"mchId" form:"mchId"`
	OutContractCode string     `json:"outContractCode" form:"outContractCode"`
	ContractCode    string     `json:"contractCode" form:"contractCode"`
	CallbackType    string     `json:"callbackType" form:"callbackType"`
	SignValid       *bool      `json:"signValid" form:"signValid"`
}
