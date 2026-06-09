package service

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
)

func TestBuildContractSignParamsMatchesGeneratedCallback(t *testing.T) {
	contract := model.Contract{
		OutContractID: "OUT-CONTRACT-001",
		OpenID:        "openid-001",
		PlanID:        "plan-001",
		ContractID:    "CONTRACT-001",
	}
	key := "mock-sign-key"

	callbackXML := Service.Callback.BuildContractCallbackXML(contract, "mch-001", model.ContractStatusActive, key)

	var req model.ContractCallbackRequest
	if err := Service.XMLCodec.Unmarshal([]byte(callbackXML), &req); err != nil {
		t.Fatalf("unmarshal callback xml: %v", err)
	}

	if err := Service.CallbackRecord.VerifyContractCallback(req, key); err != nil {
		t.Fatalf("verify callback sign: %v", err)
	}
}

func TestBuildContractSignParamsRejectsTamperedCallback(t *testing.T) {
	contract := model.Contract{
		OutContractID: "OUT-CONTRACT-001",
		OpenID:        "openid-001",
		PlanID:        "plan-001",
		ContractID:    "CONTRACT-001",
	}
	key := "mock-sign-key"

	callbackXML := Service.Callback.BuildContractCallbackXML(contract, "mch-001", model.ContractStatusActive, key)

	var req model.ContractCallbackRequest
	if err := Service.XMLCodec.Unmarshal([]byte(callbackXML), &req); err != nil {
		t.Fatalf("unmarshal callback xml: %v", err)
	}

	req.ChangeType = "DELETE"
	if err := Service.CallbackRecord.VerifyContractCallback(req, key); err == nil {
		t.Fatal("expected verify callback sign to fail after tampering")
	}
}
