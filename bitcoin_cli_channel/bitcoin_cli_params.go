package bitcoin_cli_channel

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type GlobalBitcoinCliParams struct {
	Params []string
}

func (s *GlobalBitcoinCliParams) LoadConfigPath(path string) {
	result := make([]string, 0)
	s.Params = result
	if path == "" {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("GlobalBitcoinCliParams LoadConfigPath error:", err)
		return
	}
	defer file.Close()
	data := make([]byte, 1024)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("GlobalBitcoinCliParams Read error:", err)
		return
	}
	contentString := string(data[:count])
	fmt.Println("bitcoin-cli params:", contentString)

	contentString = strings.Replace(contentString, "\n", " ", -1)
	contentString = strings.Replace(contentString, "\t", " ", -1)

	params := strings.Split(contentString, " ")
	for _, eachParam := range params {
		if eachParam != "" {
			result = append(result, eachParam)
		}
	}
	s.Params = result
	return
}

var instance *GlobalBitcoinCliParams
var once sync.Once

func GlobalParams() *GlobalBitcoinCliParams {
	once.Do(func() {
		instance = &GlobalBitcoinCliParams{}
		instance.Params = make([]string, 0)
	})
	return instance
}
