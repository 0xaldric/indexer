package known_test

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/addr"
	"github.com/xssnick/tonutils-go/address"
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
	
	var testCases = []*struct {
		name string
		addr *address.Address
		code string
		data string
		expected []any
	} {
		{
			name: "get_ton_fun_factory_data",
			addr: addr.MustFromBase64("EQBN1RnFrOzz1L9uudyKw2Ng04RiZvqNYGn6IW9iP4BwmJnP").MustToTonutils(),
			code: `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`,
			data: `te6ccgECSAEAD5QABEWAD1Uy7olDGF9bGtFySIzrn0LZr3qp/xfYI7SRaXigkwZgGAECAwQATICxorwuxQAACALGivC7FAAASy0F4AgCxorwuxQAAFA3gb9QAgPoART/APSkE/S88sgLCQIBzQUGAgAoKQEBWAcBAUgIAAAAh4AcBCx7sqmHsUj0+Ti0e/FZBqFv5QPJEbM9LKkeG/LRgXADLcce6SZW3lk39MubJ55w5v8Lt+EBs60EapzNyrKHQ0AgAgFiCgsCAs8MMgIBWCQlAvc7aLt+zMixwCSXwPg0NMDMfpA+kAx+gAxcdch+gAx+gAwc6m0AAPTH9M/7UTQ+gAB+GH6AAH4YvoAAfhj+gAB+GT6AAH4ZfoAAfhm0w8B+Gf6AAH4aNMAAfhp1DDQ+kAB+Gr6QAH4a9QB+GzUMPhtghBzYtCcUjC64wAigDQ4C/vgo+Ev4THBUIBNUFAPIUAT6AljPFgHPFszJIsjLARL0APQAywDJ+QBwdMjLAsoHy//J0FJAxwXy4r76APpA+EnAAY64MDM0NHBagEABECVwQGVwIIIQD4p+pcjLHxbLP1AD+gIhzxYBzxYTywAB+gLLAMkCcrFDMNs82zHg1AEhDwNeghDtWTE1uo8WMvhJwAGOjBAkXwRwgEBysRLbPOBDROA0NPhKEscF4wJfA4QP8vAiFBUE/tAjwgDy4sLTHyGCENRwDEG6j2E2WzU1I/ABRBQB+gD6QDD4QvhBJVlSIKgCoKkEUgO5+EH4Q6FSMLyxjrNbcFmAQAEQJXBAZXAgghAPin6lyMsfFss/UAP6AiHPFgHPFhPLAAH6AssAyQJysUMw2zzjDtsx4DOCEE4L0CC64wIhEBESA7AyNFMDcXKxEts8cIBAcrES2zz4QSOh+GH4QiGg+GKCEE6IQpBBM8j4QfoC+EL6AsnIUAT6Alj6AgHPFszJcIMHcYAMyMsDywHLCFIwy//LYRLLH8zJcPsAIiIbAcAwMzMB+gD6AFMhoBW+8uMj+EJYoPhi+EUBoPhlAdMAAZJfBOMN+E34TMj4Ss8W+EvPFszMyfhJ+EfI+EH6AvhC+gL4Q/oC+ET6AvhF+gL4RvoCyw/4SPoCywDMye1U2zETAARsIQP81DDQ+kBEMBL6APoA+kAwA6oBUiCgggnJw4CgUkC5j1z4R1IggScQqYRTIKH4QvhEofhB+EJSIqgCoakEUhC2CFIQoWah+EH4QiJZUiCoAqCpBFIFufhC+EShUlC8sY6MEHhfCHCAQHKxEts84FBkoYIJycOAofgo+Ev4TOMNIhYXA+z6APoA+kAwA6oBUiCgggnJw4CgUkC5j1z4R1IggScQqYRTIKH4QvhEofhB+EJSIqgCoakEUhC2CFIQoWah+EH4QiJZUiCoAqCpBFIFufhC+EShUlC8sY6MEHhfCHCAQHKxEts84FBkoYIJycOAofgo+Ev4TOMNIhYXAsBZMSGCEEEuGda6jtIwghAelCOGuo7H+Er4R4BAcrES2zxw+Gf4TfhMyPhKzxb4S88WzMzJ+En4R8j4QfoC+EL6AvhD+gL4RPoC+EX6AvhG+gLLD/hI+gLLAMzJ7VTg4w0iIwL6cFQgE1QUA8hQBPoCWM8WAc8WzMkiyMsBEvQA9ADLAMn5AHB0yMsCygfL/8nQEDgQJXApAhAlcEBlcCCCEA+KfqXIyx8Wyz9QA/oCIc8WAc8WE8sAAfoCywDJAnKxQzDbPPhBIaD4YfhCJaH4YvhIWKD4aPhB+Ea++EL4RKEhGAEYEEVfBXCAQHKxEts8IgT0ghA7msoAu7CPJY8i7aLt+3H4afhN0NMHIcABjoMx2zyOiwHAAo6E2zzbMeAw4tjeIcIAjohSInBysRLbPJEx4oIQI7SXYUEzyPhB+gL4QvoCychQBPoCWPoCAc8WzMlwgwdxgAzIywPLAcsIUjDL/8thEssfzMlw+wAZGiIbAfL6QPpAMPhL+EwjWXBUIBNUFAPIUAT6AljPFgHPFszJIsjLARL0APQAywDJ+QBwdMjLAsoHy//J0PhB+EOh+EX4KPhL+ExwVCATVBQDyFAE+gJYzxYBzxbMySLIywES9AD0AMsAyfkAcHTIywLKB8v/ydBwJUQTVEdgHAGs+kAx+kD6QPhG+EOh+EX4S/pEccjLAxLLB8v/ydD4KPhL+ExwVCATVBQDyFAE+gJYzxYBzxbMySLIywES9AD0AMsAyfkAcHTIywLKB8v/ydBUYzBUZkoeAHL4TfhMyPhKzxb4S88WzMzJ+En4R8j4QfoC+EL6AvhD+gL4RPoC+EX6AvhG+gLLD/hI+gLLAMzJ7VQCmoIQ/Pnlj8jLH1jPFgH6AsmCEAjw0YBSMKCCCvrwgKBtcCCCEA+KfqXIyx/LP1AG+gJQBM8WFMsBEvQAghAI8NGA+gL0AMkScds8QTNxIR0BjoIQ/Pnlj8jLH1jPFgH6AsmCEAjw0YCCCvrwgKBtcCCCEA+KfqXIyx/LP1AG+gJQBM8WFMsBEvQAghAI8NGA+gL0AMkScds8IQKOyHH6AiX6AlAE+gLJghAI8NGAUlCgggr68ICgbW1wIIIQ1V5GhsjLH8s/UAn6AhjLAFAFzxZQA88WzBT0APQAyXHbPEVAQTAfIAAscYAQyMsFUATPFlAE+gISy2rMyQH7AAG4ghAI8NGAggr68ICgcFMAghBA4QjWyMsfywBQBc8WUAPPFnH6AlAH+gIl+gISywHLAcltcCCCEA+KfqXIyx/LP1AG+gJQA88WFMsB9ACCEAjw0YD6AhL0AMlx2zwhACxxgBjIywVQBM8WUAT6AhLLaszJAfsAAChwgBjIywVQA88WUAP6AstqyQH7AAB+MdMPMPhn+E34TMj4Ss8W+EvPFszMyfhJ+EfI+EH6AvhC+gL4Q/oC+ET6AvhF+gL4RvoCyw/4SPoCywDMye1UAgEgJicAu7sivtRND6AAH4YfoAAfhi+gAB+GP6AAH4ZPoAAfhl+gAB+GbTDwH4Z/oAAfho0wAB+GnUMND6QAH4avpAAfhr1AH4bNQw+G34QfhC+EP4RPhF+Eb4R/hI+En4SvhLgA2bTmHaiaH0AAPww/QAA/DF9AAD8Mf0AAPwyfQAA/DL9AAD8M2mHgPwz/QAA/DRpgAD8NOoYaH0gAPw1fSAA/DXqAPw2ahh8NvwjqQhAk4hUwlD8I3wg0KkIXksYfCN8INDvfCD8ISkQVAFQVIJAApbXWfaiaH0AAPww/QAA/DF9AAD8Mf0AAPwyfQAA/DL9AAD8M2mHgPwz/QAA/DRpgAD8NOoYaH0gAPw1fSAA/DXqAPw2ahh8NvwhfCCpEFQBUFSCQART/APSkE/S88sgLKgEU/wD0pBP0vPLICzsCAWIrLAICyy0uABug9gXaiaH0AfSB9IGoYQIBIC8wAIPSAINch7UTQ+gD6QPpA1DAE0x+CEBeNRRlSILqCEHvdl94TuhKx8uLF0z8x+gAwE6BQI8hQBPoCWM8WAc8WzMntVICAdQxMgIBWDM0AMMIMcAkl8E4AHQ0wMBcbCVE18D8BDg+kD6QDH6ADFx1yH6ADH6ADBzqbQAAtMfghAPin6lUiC6lTE0WfAN4IIQF41FGVIgupYxREQD8A7gNYIQWV8HvLqTWfAP4F8EhA/y8IAARPpEMHC68uFNgAfFQPTP/oA+kAh8AHtRND6APpA+kDUMFE2oVIqxwXy4sEowv/y4sJUNEJwVCATVBQDyFAE+gJYzxYBzxbMySLIywES9AD0AMsAySD5AHB0yMsCygfL/8nQBPpA9AQx+gAg10nCAPLixHeAGMjLBVAIzxZw+gIXy2sTzINQIBIDY3AK6CEBeNRRnIyx8Zyz9QB/oCIs8WUAbPFiX6AlADzxbJUAXMI5FykXHiUAioE6CCCOThwKoAggiYloCgoBS88uLFBMmAQPsAECPIUAT6AljPFgHPFszJ7VQD9ztRND6APpA+kDUMAjTP/oAUVGgBfpA+kBTW8cFVHNtcFQgE1QUA8hQBPoCWM8WAc8WzMkiyMsBEvQA9ADLAMn5AHB0yMsCygfL/8nQUA3HBRyx8uLDCvoAUaihggiYloCCCJiWgBK2CKGCCOThwKAYoSfjDyXXCwHDACOA4OToA2ztRND6APpA+kDUMAfTP/oA+kAwUVGhUknHBfLiwSfC//LiwoII5OHAqgAWoBa88uLDghB73ZfeyMsfFcs/UAP6AiLPFgHPFslxgBjIywUkzxZw+gLLaszJgED7AEATyFAE+gJYzxYBzxbMye1UgAHBSeaAYoYIQc2LQnMjLH1Iwyz9Y+gJQB88WUAfPFslxgBDIywUkzxZQBvoCFctqFMzJcfsAECQQIwAOEEkQODdfBAB2wgCwjiGCENUydttwgBDIywVQCM8WUAT6AhbLahLLHxLLP8ly+wCTNWwh4gPIUAT6AljPFgHPFszJ7VQCAWI8PQICzD4/AgN6YEZHAfXZBjgEkvgfAA6GmBgLjYSS+B8H0gfSAY/QAYuOuQ/QAY/QAYOdTaAAFpj+mf9qJofQB9IGpqGEAKqThdRxmampsbKJpjgvlwJIF9IH0AahgQaEAwa5D9ABgSiBooIXgGilAoGeQoAn0BLGeLZmZk9qpwQQg97svvKThdRAAJO78FCIBuCoQCaoKAeQoAn0BLGeLAOeLZmSRZGWAiXoAegBlgGSQfIA4OmRlgWUD5f/k6DvADGRlgqxniygCfQEJ5bWJZmZkuP2AQT0juA2NzcB+gD6QPgoVBIGcFQgE1QUA8hQBPoCWM8WAc8WzMkiyMsBEvQA9ADLAMn5AHB0yMsCygfL/8nQUAbHBfLgSqEDRUXIUAT6AljPFszMye1UAfpAMCDXCwHDAJFb4w3gghAsdrlzUnC64wI1NzcjwAPjAjUCwARBQkNEAD6CENUydttwgBDIywVQA88WIvoCEstqyx/LP8mAQvsAAf42XwOCCJiWgBWgFbzy4EsC+kDTADCVyCHPFsmRbeKCENFzVABwgBjIywVQBc8WJPoCFMtqE8sfFMs/I/pEMHC6jjP4KEQDcFQgE1QUA8hQBPoCWM8WAc8WzMkiyMsBEvQA9ADLAMn5AHB0yMsCygfL/8nQzxaWbCJwAcsB4vQARQA0M1A1xwXy4EkD+kAwWchQBPoCWM8WzMzJ7VQAQo4YUSTHBfLgSdQwQwDIUAT6AljPFszMye1U4F8FhA/y8AAKyYBA+wAAfa289qJofQB9IGpqGDYY/BQAuCoQCaoKAeQoAn0BLGeLAOeLZmSRZGWAiXoAegBlgGT8gDg6ZGWBZQPl/+ToQAAfrxb2omh9AH0gamoYP6qQQA==`,
			expected: []any {
				addr.MustFromBase64("EQB6qZd0ShjC-tjWi5JEZ1z6Fs171U_4vsEdpItLxQSYM5eb").MustToTonutils(),
				big.NewInt(0),
				big.NewInt(0),
			},
		},
	}

	for _, test := range testCases {
		ret := execGetMethod(t, i, test.addr, test.name, test.code, test.data)
		require.Equal(t, test.expected, ret)
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
		// {
		// 	name: `buy_jetton_success`,
		// 	boc: `te6cckEBAwEASAABCCO0l2EBAVtDuLh8BxvWIClwvMWAAIdVaP0xUNB/vAwsflbhFLo6AflbG1QXPURgp2Si3SMQAgAcUXoaJR94CvBHPi7d8HUYE6SL`,
		// 	expected: `{"data":{"amount_in":"999000000","amount_out":"7835257993739461","sender_address":"EQAEOqtH6YqGg_3gYWPytwil0dAPytjaoLnqIwU7JRbpGLnL","reserve_data":{"reserve_0":"101496017399","reserve_1":"788208267189678197"}}}`,
		// },
		// {
		// 	name: `sell_jetton_success`,
		// 	boc: `te6cckEBAwEASAABCE6IQpABAVtw4OZJQBUsZB3rpYmAGZVLanby5XtjYZyV+1h/YRPWh4RHlLwJTwPqcq0yP08wAgAcUXZhbKN4CwwdXlhOrTqzMtQV`,
		// 	expected: `{"data":{"amount_in":"3956474816582342","amount_out":"501982601","sender_address":"EQDMqltTt5cr2xsM5K_aw_sInrQ8Ijyl4Ep4H1OVaZH6efgU","reserve_data":{"reserve_0":"100497017399","reserve_1":"796043525183417658"}}}`,
		// },
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)
		require.Equal(t, test.expected, j)
	}
}
