package service

import (
	"bytes"
	"encoding/xml"
)

type xmlCodec struct{}

func (s *xmlCodec) Marshal(v any) (string, error) {
	buf := bytes.NewBufferString(xml.Header)
	encoder := xml.NewEncoder(buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *xmlCodec) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
