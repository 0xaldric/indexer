package known_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
)

func TestOperationDesc_DedustV2Vault(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("dedust_v2.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "dedust_v2_vault" {
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
			name:     `dedust_v2_swap`,
			boc:      `b5ee9c7201010201003d000165ea06185d000001903a2b0e3a405f5e100800efb1c38b87907acdb11f2d65b65da200b72e4716e142d83aa8b6067f01e1b3c2040100090000000002`,
			expected: `{"data":{"creator":"EQDMqltTt5cr2xsM5K_aw_sInrQ8Ijyl4Ep4H1OVaZH6efgU","ton_fun_address":"EQDqcMtClISI7J12xKOuMmRGbnUzjnXdbNjChED0VXmvVaex","jetton_master_address":"EQBK-EcElH-nfBRKqFqhaFxRDFW296kVkzyO0DMX67vdmko-"}}`,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		fmt.Println(j)
		require.Equal(t, test.expected, j)
	}
}
