package test

import (
	"OrdDeFi-Virtual-Machine/tx_utils"
	"testing"
)

func TestAddressValidator(t *testing.T) {
	if tx_utils.IsValidateBitcoinAddress("abcd") {
		t.Errorf("abcd should be an invalid address")
	}
	if tx_utils.IsValidateBitcoinAddress("bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h") == false {
		t.Errorf("bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h should be a valid address")
	}
	if tx_utils.IsValidateBitcoinAddress("bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h") == false {
		t.Errorf("bc1quhruqrghgcca950rvhtrg7cpd7u8k6svpzgzmrjy8xyukacl5lkq0r8l2d should be a valid address")
	}
	if tx_utils.IsValidateBitcoinAddress("18xKwnk85zMj1ajEkZkQokJhrVaj89jr25") == false {
		t.Errorf("18xKwnk85zMj1ajEkZkQokJhrVaj89jr25 should be a valid address")
	}
	if tx_utils.IsValidateBitcoinAddress("18xKwnk85zMj1ajEkZkQokJhrVaj89jr26") {
		t.Errorf("18xKwnk85zMj1ajEkZkQokJhrVaj89jr26 should be an invalid address")
	}
	if tx_utils.IsValidateBitcoinAddress("bc1pv2s3twphzh3xf9d39tsa3nmt676hzgvv4u2qa98t3tywmu6njrnsgxehp8") == false {
		t.Errorf("bc1pv2s3twphzh3xf9d39tsa3nmt676hzgvv4u2qa98t3tywmu6njrnsgxehp8 should be a valid address")
	}
	if tx_utils.IsValidateBitcoinAddress("bc1pv2s3twphzh3xf9d39tsa3nmt676hzgvv4u2qa98t3tywmu6njrnsgxehp9") {
		t.Errorf("bc1pv2s3twphzh3xf9d39tsa3nmt676hzgvv4u2qa98t3tywmu6njrnsgxehp9 should be an invalid address")
	}
}
