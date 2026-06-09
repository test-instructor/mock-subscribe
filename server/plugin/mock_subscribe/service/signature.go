package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type signature struct{}

func (s *signature) ShouldVerify(merchantVerifySign bool) bool {
	return merchantVerifySign
}

func (s *signature) RequireSignIfNeeded(merchantVerifySign bool, sign string) bool {
	if !s.ShouldVerify(merchantVerifySign) {
		return true
	}
	return strings.TrimSpace(sign) != ""
}

func (s *signature) VerifyIfNeeded(merchantVerifySign bool, params map[string]string, key string) error {
	if !s.ShouldVerify(merchantVerifySign) {
		return nil
	}
	return s.Verify(params, key)
}

func (s *signature) Sign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || strings.TrimSpace(v) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys)+1)
	for _, keyName := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", keyName, params[keyName]))
	}
	parts = append(parts, fmt.Sprintf("key=%s", key))
	sum := md5.Sum([]byte(strings.Join(parts, "&")))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func (s *signature) Verify(params map[string]string, key string) error {
	sign := params["sign"]
	if strings.TrimSpace(sign) == "" {
		return errors.New("缺少签名")
	}
	if s.Sign(params, key) != strings.ToUpper(sign) {
		return errors.New("签名校验失败")
	}
	return nil
}
