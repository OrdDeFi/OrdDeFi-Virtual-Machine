package memory_const

import (
	"OrdDeFi-Virtual-Machine/safe_number"
	"encoding/json"
	"errors"
)

type coinMetaSerialization struct {
	Max     string `json:"max"`
	Lim     string `json:"lim"`
	AddrLim string `json:"alim"`
	Desc    string `json:"desc"`
	Icon    string `json:"icon"`
}

type CoinMeta struct {
	Max     *safe_number.SafeNum
	Lim     *safe_number.SafeNum
	AddrLim *safe_number.SafeNum
	Desc    string
	Icon    string
}

func (coinMeta CoinMeta) JsonString() (*string, error) {
	if coinMeta.Max == nil || coinMeta.Lim == nil {
		return nil, errors.New("error: Max or Lim is nil")
	}
	s := new(coinMetaSerialization)
	s.Max = coinMeta.Max.String()
	s.Lim = coinMeta.Lim.String()
	if coinMeta.AddrLim != nil {
		s.AddrLim = coinMeta.AddrLim.String()
	} else {
		s.AddrLim = s.Max
	}
	s.Desc = coinMeta.Desc
	s.Icon = coinMeta.Icon
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	result := string(jsonData)
	return &result, nil
}

func CoinMetaFromJsonString(jsonString string) (*CoinMeta, error) {
	var s coinMetaSerialization
	err := json.Unmarshal([]byte(jsonString), &s)
	if err != nil {
		return nil, err
	}
	result := new(CoinMeta)
	result.Max = safe_number.SafeNumFromString(s.Max)
	result.Lim = safe_number.SafeNumFromString(s.Lim)
	result.AddrLim = safe_number.SafeNumFromString(s.AddrLim)
	result.Desc = s.Desc
	result.Icon = s.Icon
	return result, nil
}
