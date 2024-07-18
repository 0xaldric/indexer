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
			boc:      `te6ccgECDAEAAgYAAQjvEsPUAQKtgB7UeDTNXCzm+dmoHjIHSMNaygwID8UwqG9sW2JtZLIf9dA83/2ic7he3zO4gOLmbPYJL2SDdOJ3WEU63yPgUmIdAh4Z4Mm6sj/4+FAh4Z4Mm6sj/4+EAgMCAUgEBQIBSAgJAUK/ncNM7ouJtMUP+8f2MjMG+ymUtNgAq0zOb19FbqwWY0cGAUK/p5SaNoLgrMiYJBmzs13TIdPDxSg/eS56ZAD9/N5gAKQHAH+AHp+EJeM8TboyH9bs2/z9K/cCH4OvzaQoLYr/yPTr8a0A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAMEjCc5UABAH+ACqGteg9BMYS86JRRHCr+ZoKfQjR8Gv9iYSYT54f7vP6A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAMEjCc5UABAUK/ncNM7ouJtMUP+8f2MjMG+ymUtNgAq0zOb19FbqwWY0cKAUK/p5SaNoLgrMiYJBmzs13TIdPDxSg/eS56ZAD9/N5gAKQLAIGAHp+EJeM8TboyH9bs2/z9K/cCH4OvzaQoLYr/yPTr8a0A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAMEjCc5UAAEACBgAqhrXoPQTGEvOiUURwq/maCn0I0fBr/YmEmE+eH+7z+gN4Lazp2QAAAAAAAAGdlx5P6EAedAAAADBIwnOVAABA=`,
			expected: `{"data":{"user_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","pool_key":78931953226258007893740732125093179493168800338006405804546432085500874625808,"tokens_in":{"13462138750431230630641883345896090919134063452067867993909573781623785546567":{"token_wallet_address":"EQD0_CEvGeJt0ZD-t2bf5-lfuBD8HX5tIUFsV_5Hp1-NaNyZ","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"10000000000000"},"46850781108221188480805398713470464247352601057275293894364456324824800690340":{"token_wallet_address":"EQBVDWvQegmMJedEoojhV_M0FPoRo-DX-xMJMJ88P93n9OAk","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"10000000000000"}},"tokens_dict":{"13462138750431230630641883345896090919134063452067867993909573781623785546567":{"token_wallet_address":"EQD0_CEvGeJt0ZD-t2bf5-lfuBD8HX5tIUFsV_5Hp1-NaNyZ","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"10000000000000","collected_protocol_fee":"0"},"46850781108221188480805398713470464247352601057275293894364456324824800690340":{"token_wallet_address":"EQBVDWvQegmMJedEoojhV_M0FPoRo-DX-xMJMJ88P93n9OAk","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"10000000000000","collected_protocol_fee":"0"}},"amount_apt_out":"19999999999999999996400","total_supply_apt":"19999999999999999996400"}}`,
		},
		{
			name:     `swap_success`,
			boc:      `te6ccgECCQEAAZEAAQh91/yjAQPFgB7UeDTNXCzm+dmoHjIHSMNaygwID8UwqG9sW2JtZLIf8APajwaZq4Wc3zs1A8ZA6RhrWUGBAfimFQ3ti2xNrJZD/FnwNJk+oiTK9A0Rw+i2o08kqMB3+Pk9QOXKbSl59PCrAgMEAIWAGiU7vjmV6RNkpCRE8os8X88azNR3pDl+zCSh/acXql0QAepcOpHUVOz/03Y+N85nDfYYqC3Pd1YStUuxY4a3qvQOAAoicQITbwIBWAUGAUK/uhk9udhjF25odbeccpzwDGuvQOxl/RWZmHews+MHpFQHAUK/j1g63I1XhE+6CxFQuJHZIxw20TvX8wdIg1+SpV0Fn7kIAIOAGiU7vjmV6RNkpCRE8os8X88azNR3pDl+zCSh/acXql0A3gtrOnZAAAAAAAAAZ2XHk/oQB50AAAAODjX6kxrqYBAAhYAPUuHUjqKnZ/6bsfG+czhvsMVBbnu6sJWqXYscNb1XoGDeC2s6dkAAAAAAAABnZceT+hAHnQAAAA4HGv1JjItkIZA=`,
			expected: `{"data":{"user_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","to_address":"EQD2o8GmauFnN87NQPGQOkYa1lBgQH4phUN7YtsTayWQ_zvy","pool_key":10170062460433864532360247382786154450879953670526333890165448804269308066858,"token_info":{"token_in_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","token_out_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib"},"amount_info":{"amount_in":"10000","amount_out":"4975"},"tokens_dict":{"84174787030012514383756880442724328713373342661635207814144282419731859219540":{"token_wallet_address":"EQDRKd3xzK9ImyUhIieUWeL-eNZmo70hy_ZhJQ_tOL1S6Jn_","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"2000000000030000","collected_protocol_fee":"0"},"93784648445798056305784803018835104413884644461776470512879246004088539029433":{"token_wallet_address":"EQB6lw6kdRU7P_Tdj43zmcN9hioLc93VhK1S7Fjhreq9Awib","weight":500000000000000000,"scaling_factor":1000000000000000000000000000,"amount":"999999999985074","collected_protocol_fee":"12"}}}}`,
		},
	}

	for _, test := range testCases {

		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}
