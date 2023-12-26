package memory_const

import (
	"OrdDeFi-Virtual-Machine/safe_number"
	"encoding/json"
	"errors"
)

type lpMetaSerialization struct {
	LTick string `json:"ltick"`
	RTick string `json:"rtick"`
	LAmt  string `json:"lamt"`
	RAmt  string `json:"ramt"`
}

type LPMeta struct {
	LTick string
	RTick string
	LAmt  *safe_number.SafeNum
	RAmt  *safe_number.SafeNum
}

func (lpMeta LPMeta) JsonString() (*string, error) {
	if lpMeta.LAmt == nil || lpMeta.RAmt == nil {
		return nil, errors.New("error: LAmt or Ramt is nil")
	}
	s := new(lpMetaSerialization)
	s.LTick = lpMeta.LTick
	s.RTick = lpMeta.RTick
	s.LAmt = lpMeta.LAmt.String()
	s.RAmt = lpMeta.RAmt.String()
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	result := string(jsonData)
	return &result, nil
}

func LPMetaFromJsonString(jsonString string) (*LPMeta, error) {
	var s lpMetaSerialization
	err := json.Unmarshal([]byte(jsonString), &s)
	if err != nil {
		return nil, err
	}
	result := new(LPMeta)
	result.LTick = s.LTick
	result.RTick = s.RTick
	result.LAmt = safe_number.SafeNumFromString(s.LAmt)
	result.RAmt = safe_number.SafeNumFromString(s.RAmt)
	return result, nil
}
