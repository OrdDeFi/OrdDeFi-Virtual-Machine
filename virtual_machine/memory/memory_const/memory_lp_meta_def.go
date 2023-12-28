package memory_const

import (
	"OrdDeFi-Virtual-Machine/safe_number"
	"encoding/json"
	"errors"
	"strings"
)

type lpMetaSerialization struct {
	LTick string `json:"ltick"`
	RTick string `json:"rtick"`
	LAmt  string `json:"lamt"`
	RAmt  string `json:"ramt"`
	Total string `json:"total"`
}

type LPMeta struct {
	LTick string
	RTick string
	LAmt  *safe_number.SafeNum
	RAmt  *safe_number.SafeNum
	Total *safe_number.SafeNum
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
	s.Total = lpMeta.Total.String()
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
	result.Total = safe_number.SafeNumFromString(s.Total)
	return result, nil
}

func LPNameByTicks(tick1 string, tick2 string) *string {
	cmpRes := strings.Compare(tick1, tick2)
	if cmpRes < 0 {
		lpName := tick1 + "-" + tick2
		return &lpName
	} else if cmpRes > 0 {
		lpName := tick2 + "-" + tick1
		return &lpName
	}
	return nil
}
