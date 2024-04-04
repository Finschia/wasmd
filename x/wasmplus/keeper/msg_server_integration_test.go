package keeper_test

import (
	_ "embed"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/wasmd/appplus"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

//go:embed testdata/reflect.wasm
var wasmContract []byte

func TestStoreAndInstantiateContract(t *testing.T) {
	wasmApp := appplus.Setup(t)
	ctx := wasmApp.BaseApp.NewContextLegacy(false, tmproto.Header{Time: time.Now()})

	var myAddress sdk.AccAddress = make([]byte, wasmtypes.ContractAddrLen)

	specs := map[string]struct {
		addr       string
		permission *wasmtypes.AccessConfig
		expEvents  []abci.Event
		expErr     bool
	}{
		"address can instantiate a contract when permission is everybody": {
			addr:       myAddress.String(),
			permission: &wasmtypes.AllowEverybody,
			expEvents: []abci.Event{
				{
					Type: "store_code",
					Attributes: []abci.EventAttribute{
						{
							Key:   "code_checksum",
							Value: "2843664c3b6c1de8bdeca672267c508aeb79bb947c87f75d8053f971d8658c89",
							Index: false,
						}, {
							Key:   "code_id",
							Value: "1",
							Index: false,
						},
					},
				},
				{
					Type: "instantiate",
					Attributes: []abci.EventAttribute{
						{
							Key:   "_contract_address",
							Value: "link14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sgf2vn8",
							Index: false,
						}, {
							Key:   "code_id",
							Value: "1",
							Index: false,
						},
					},
				},
			},
			expErr: false,
		},
		"address cannot instantiate a contract when permission is nobody": {
			addr:       myAddress.String(),
			permission: &wasmtypes.AllowNobody,
			expErr:     true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			xCtx, _ := ctx.CacheContext()
			// when
			msg := &types.MsgStoreCodeAndInstantiateContract{
				Sender:                spec.addr,
				WASMByteCode:          wasmContract,
				InstantiatePermission: spec.permission,
				Admin:                 myAddress.String(),
				Label:                 "test",
				Msg:                   []byte(`{}`),
				Funds:                 sdk.Coins{},
			}
			rsp, err := wasmApp.MsgServiceRouter().Handler(msg)(xCtx, msg)

			// then
			if spec.expErr {
				require.Error(t, err)
				return
			}

			var storeAndInstantiateResponse types.MsgStoreCodeAndInstantiateContractResponse
			require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &storeAndInstantiateResponse))

			// check event
			assert.Equal(t, spec.expEvents, rsp.Events)

			require.NoError(t, err)
		})
	}
}
