package keeper_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/wasmd/appplus"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

//go:embed testdata/reflect.wasm
var wasmContract []byte

func TestStoreAndInstantiateContract(t *testing.T) {
	wasmApp := appplus.Setup(false)
	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	var myAddress sdk.AccAddress = make([]byte, wasmtypes.ContractAddrLen)

	specs := map[string]struct {
		addr       string
		permission *wasmtypes.AccessConfig
		expErr     bool
	}{
		"address can instantiate a contract when permission is everybody": {
			addr:       myAddress.String(),
			permission: &wasmtypes.AllowEverybody,
			expErr:     false,
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
			events := rsp.Events
			assert.Equal(t, 3, len(events))
			assert.Equal(t, "store_code", events[0].Type)
			assert.Equal(t, 2, len(events[0].Attributes))
			assert.Equal(t, "code_checksum", string(events[0].Attributes[0].Key))
			assert.Equal(t, "code_id", string(events[0].Attributes[1].Key))
			assert.Equal(t, "1", string(events[0].Attributes[1].Value))
			assert.Equal(t, "message", events[1].Type)
			assert.Equal(t, 2, len(events[1].Attributes))
			assert.Equal(t, "module", string(events[1].Attributes[0].Key))
			assert.Equal(t, "wasm", string(events[1].Attributes[0].Value))
			assert.Equal(t, "sender", string(events[1].Attributes[1].Key))
			assert.Equal(t, "instantiate", events[2].Type)
			assert.Equal(t, "_contract_address", string(events[2].Attributes[0].Key))
			assert.Equal(t, storeAndInstantiateResponse.Address, string(events[2].Attributes[0].Value))
			assert.Equal(t, "code_id", string(events[2].Attributes[1].Key))
			assert.Equal(t, "1", string(events[2].Attributes[1].Value))

			require.NoError(t, err)
		})
	}
}
