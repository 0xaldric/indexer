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
		if i.Name == "multi_token_pool" {
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
			boc:      `te6ccgECDAEAAgEAAQjvEsPUAQKbgB7UeDTNXCzm+dmoHjIHSMNaygwID8UwqG9sW2JtZLIf4s+BpMn1ESZXoGiOH0W1GnklRgO/x8nqBy5TaUvPp4VdgSt4q7vJv59/UmtAAgMCAVgEBQIBWAgJAUK/uhk9udhjF25odbeccpzwDGuvQOxl/RWZmHews+MHpFQGAUK/j1g63I1XhE+6CxFQuJHZIxw20TvX8wdIg1+SpV0Fn7kHAIGAGiU7vjmV6RNkpCRE8os8X88azNR3pDl+zCSh/acXql0A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAODjX6kxoAAQCBgA9S4dSOoqdn/pux8b5zOG+wxUFue7qwlapdixw1vVegYN4Lazp2QAAAAAAAAGdlx5P6EAedAAAADgca/UmNAAEBQr+6GT252GMXbmh1t5xynPAMa69A7GX9FZmYd7Cz4wekVAoBQr+PWDrcjVeET7oLEVC4kdkjHDbRO9fzB0iDX5KlXQWfuQsAg4AaJTu+OZXpE2SkJETyizxfzxrM1HekOX7MJKH9pxeqXQDeC2s6dkAAAAAAAABnZceT+hAHnQAAAA4ONfqTGgAAEACDgA9S4dSOoqdn/pux8b5zOG+wxUFue7qwlapdixw1vVegYN4Lazp2QAAAAAAAAGdlx5P6EAedAAAADgca/UmNAAAQ`,
			expected: `{"data":{"user_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","pool_key":10170062460433864532360247382786154450879953670526333890165448804269308066858,"tokens_in":{"84174787030012514383756880442724328713373342661635207814144282419731859219540":{"token_wallet_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"2000000000000000"},"93784648445798056305784803018835104413884644461776470512879246004088539029433":{"token_wallet_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"1000000000000000"}},"tokens_dict":{"84174787030012514383756880442724328713373342661635207814144282419731859219540":{"token_wallet_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"2000000000000000","collected_protocol_fee":"0"},"93784648445798056305784803018835104413884644461776470512879246004088539029433":{"token_wallet_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"1000000000000000","collected_protocol_fee":"0"}},"amount_bpt_out":"2828427124746190095033558"}}`,
		},
		// {
		// 	name:     `provide_liquidity_success`,
		// 	boc:      ``,
		// 	expected: ``,
		// }, 
		// {
		// 	name:     `remove_liquidity_success`,
		// 	boc:      ``,
		// 	expected: ``,
		// }, 
		{
			name:     `swap_success`,
			boc:      `te6ccgECCQEAAZEAAQh91/yjAQPFgB7UeDTNXCzm+dmoHjIHSMNaygwID8UwqG9sW2JtZLIf8APajwaZq4Wc3zs1A8ZA6RhrWUGBAfimFQ3ti2xNrJZD/FnwNJk+oiTK9A0Rw+i2o08kqMB3+Pk9QOXKbSl59PCrAgMEAIWAGiU7vjmV6RNkpCRE8os8X88azNR3pDl+zCSh/acXql0QAepcOpHUVOz/03Y+N85nDfYYqC3Pd1YStUuxY4a3qvQOAAoicQITbwIBWAUGAUK/uhk9udhjF25odbeccpzwDGuvQOxl/RWZmHews+MHpFQHAUK/j1g63I1XhE+6CxFQuJHZIxw20TvX8wdIg1+SpV0Fn7kIAIOAGiU7vjmV6RNkpCRE8os8X88azNR3pDl+zCSh/acXql0A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAODjX6kxrqYBAAhYAPUuHUjqKnZ/6bsfG+czhvsMVBbnu6sJWqXYscNb1XoGDeC2s6dkAAAAAAAABnZceT+hAHnQAAAA4HGv1JjItkIZA=`,
			expected: `{"data":{"user_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","to_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","pool_key":10170062460433864532360247382786154450879953670526333890165448804269308066858,"token_info":{"token_in_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","token_out_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib"},"amount_info":{"amount_in":"10000","amount_out":"4975"},"tokens_dict":{"84174787030012514383756880442724328713373342661635207814144282419731859219540":{"token_wallet_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"2000000000030000","collected_protocol_fee":"0"},"93784648445798056305784803018835104413884644461776470512879246004088539029433":{"token_wallet_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"999999999985074","collected_protocol_fee":"12"}}}}`,
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		// fmt.Println(j)
		require.Equal(t, test.expected, j)
	}
}