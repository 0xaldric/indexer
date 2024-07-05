package known_test

import (
	"encoding/json"
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
			name:     `dedust_v2_native_vault_swap`,
			boc:      `te6ccgECBAEAARUAAW/qBhhdAAAAC7aqBAJDuaygCAGCS6iT9M8XO5NnkRWXlIJQ+Q/SKfNidLJizQ55msajrlFUHYpjVAECjgAAAACAHPZ3UMgBLE+iKjQMWNLQ66In7oA5/Kc9w71rTTvTHDDQAqXyHUubwW9DTF/1P+4zQ3UO3VamYm1HG2yWpAqZbyJHAgMAjYAK225Dv75o8Qi0P+MQQO9exLtRbkG50bqQuDNePAuKDPACpfIdS5vBb0NMX/U/7jNDdQ7dVqZibUcbbJakCplvIkTmJaAgAIeACttuQ7++aPEItD/jEEDvXsS7UW5BudG6kLgzXjwLigzwAqXyHUubwW9DTF/1P+4zQ3UO3VamYm1HG2yWpAqZbyJEIA==`,
			expected: `{"query_id":50309235714,"amount":"1000000000","pool_addr":"EQDBJdRJ-meLncmzyIrLykEofIfpFPmxOlkxZoc8zWNR14J0","swap_kind":false,"swap_limit":"91299030581","next":null,"swap_params":{"deadline":0,"recipient_addr":"EQDns7qGQAlifRFRoGLGloddET90Ac_lOe4d61pp3pjhhoqF","referral_addr":"EQCpfIdS5vBb0NMX_U_7jNDdQ7dVqZibUcbbJakCplvIkUxK","fulfill_payload":"te6cckEBAQEASQAAjYAK225Dv75o8Qi0P+MQQO9exLtRbkG50bqQuDNePAuKDPACpfIdS5vBb0NMX/U/7jNDdQ7dVqZibUcbbJakCplvIkTmJaAgI2C0LA==","reject_payload":"te6cckEBAQEARgAAh4AK225Dv75o8Qi0P+MQQO9exLtRbkG50bqQuDNePAuKDPACpfIdS5vBb0NMX/U/7jNDdQ7dVqZibUcbbJakCplvIkQgvGlNKQ=="}}`,
		},
		{
			name: `dedust_v2_native_vault_add_lp`,
			boc: `te6cckEBAgEATAABadVeRoZe2qlyA22hJWBAnEn3AAAIBYidTKWoElCzjPtInJlHW6ysthxRL6yBRYo39m4bEO/xAQAjVhFrCdG2BAnEn3AAUHhH1wzYwS7rzg==`,
			expected: `{"body_data":"te6cckEBAgEASAABYV7aqXIDbaElYECcSfcAAAgFiJ1MpagSULOM+0icmUdbrKy2HFEvrIFFijf2bhsQ7/EBACNWEWsJ0bYECcSfcABQeEfXDNiC1CeB"}`,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}
