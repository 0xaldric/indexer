package known_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
)

func TestOperationDesc_MultiTokenPool(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("multi_token_pool.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "multi-token-pool" {
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
			name:     `create_pool_success`,
			boc:      ``,
			expected: ``,
		},
		{
			name:     `provide_liquidity_success`,
			boc:      ``,
			expected: ``,
		}, {
			name:     `remove_liquidity_success`,
			boc:      ``,
			expected: ``,
		}, {
			name:     `swap_success`,
			boc:      ``,
			expected: ``,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}