package known_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
)

func TestGetMethodDesc_TonfunFactory(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("tonfun.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "tonfun_factory" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

}

func TestOperationDesc_TonfunFactory(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("tonfun.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "tonfun_factory" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

	var testCases = []*struct {
		name     string
		boc      string
		expected string
	}{
		{
			name:     `create_ton_fun_success`,
			boc:      `te6cckEBAgEAbgABCAoIhVwBAMmAGZVLanby5XtjYZyV+1h/YRPWh4RHlLwJTwPqcq0yP08wA6nDLQpSEiOyddsSjrjJkRm51M4513WzYwoRA9FV5r1WACV8I4JKP9O+CiVULVC0LiiGKtt71IrJnkdoGYv13e7NQJ/1VLY=`,
			expected: `{"data":{"creator":"EQDMqltTt5cr2xsM5K_aw_sInrQ8Ijyl4Ep4H1OVaZH6efgU","ton_fun_address":"EQDqcMtClISI7J12xKOuMmRGbnUzjnXdbNjChED0VXmvVaex","jetton_master_address":"EQBK-EcElH-nfBRKqFqhaFxRDFW296kVkzyO0DMX67vdmko-"}}`,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}

func TestOperationDesc_Tonfun(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("tonfun.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "tonfun" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

	var testCases = []*struct {
		name     string
		boc      string
		expected string
	}{
		{
			name: `buy_jetton_success`,
			boc: `te6cckEBAwEASAABCCO0l2EBAVtDuLh8BxvWIClwvMWAAIdVaP0xUNB/vAwsflbhFLo6AflbG1QXPURgp2Si3SMQAgAcUXoaJR94CvBHPi7d8HUYE6SL`,
			expected: `{"data":{"amount_in":"999000000","amount_out":"7835257993739461","sender_address":"EQAEOqtH6YqGg_3gYWPytwil0dAPytjaoLnqIwU7JRbpGLnL","reserve_data":{"reserve_0":"101496017399","reserve_1":"788208267189678197"}}}`,
		},
		{
			name: `sell_jetton_success`,
			boc: `te6cckEBAwEASAABCE6IQpABAVtw4OZJQBUsZB3rpYmAGZVLanby5XtjYZyV+1h/YRPWh4RHlLwJTwPqcq0yP08wAgAcUXZhbKN4CwwdXlhOrTqzMtQV`,
			expected: `{"data":{"amount_in":"3956474816582342","amount_out":"501982601","sender_address":"EQDMqltTt5cr2xsM5K_aw_sInrQ8Ijyl4Ep4H1OVaZH6efgU","reserve_data":{"reserve_0":"100497017399","reserve_1":"796043525183417658"}}}`,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}
