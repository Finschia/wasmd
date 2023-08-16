package keeper_test

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/wasmd/app"
	"github.com/Finschia/wasmd/x/wasm/types"
)

//go:embed testdata/reflect.wasm
var wasmContract []byte

func TestStoreCode(t *testing.T) {
	wasmApp := app.Setup(false)
	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{})
	_, _, sender := testdata.KeyTestPubAddr()
	msg := types.MsgStoreCodeFixture(func(m *types.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender.String()
	})
	expHash := sha256.Sum256(wasmContract)

	// when
	rsp, err := wasmApp.MsgServiceRouter().Handler(msg)(ctx, msg)

	// check event
	require.Equal(t, 2, len(rsp.Events))
	assert.Equal(t, "message", rsp.Events[0].Type)
	assert.Equal(t, 2, len(rsp.Events[0].Attributes))
	assert.Equal(t, "module", string(rsp.Events[0].Attributes[0].Key))
	assert.Equal(t, "wasm", string(rsp.Events[0].Attributes[0].Value))
	assert.Equal(t, "sender", string(rsp.Events[0].Attributes[1].Key))
	assert.Equal(t, sender.String(), string(rsp.Events[0].Attributes[1].Value))
	assert.Equal(t, "store_code", rsp.Events[1].Type)
	assert.Equal(t, 2, len(rsp.Events[1].Attributes))
	assert.Equal(t, "code_checksum", string(rsp.Events[1].Attributes[0].Key))
	assert.Equal(t, hex.EncodeToString(expHash[:]), string(rsp.Events[1].Attributes[0].Value))
	assert.Equal(t, "code_id", string(rsp.Events[1].Attributes[1].Key))
	assert.Equal(t, "1", string(rsp.Events[1].Attributes[1].Value))

	// then
	require.NoError(t, err)
	var result types.MsgStoreCodeResponse
	require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &result))
	assert.Equal(t, uint64(1), result.CodeID)
	assert.Equal(t, expHash[:], result.Checksum)
	// and
	info := wasmApp.WasmKeeper.GetCodeInfo(ctx, 1)
	assert.NotNil(t, info)
	assert.Equal(t, expHash[:], info.CodeHash)
	assert.Equal(t, sender.String(), info.Creator)
	assert.Equal(t, types.DefaultParams().InstantiateDefaultPermission.With(sender), info.InstantiateConfig)
}

func TestInstantiateContract(t *testing.T) {
	wasmApp := app.Setup(false)
	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	var (
		myAddress sdk.AccAddress = make([]byte, types.ContractAddrLen)
	)

	specs := map[string]struct {
		addr       string
		permission *types.AccessConfig
		expErr     bool
	}{
		"address can instantiate a contract when permission is everybody": {
			addr:       myAddress.String(),
			permission: &types.AllowEverybody,
			expErr:     false,
		},
		"address cannot instantiate a contract when permission is nobody": {
			addr:       myAddress.String(),
			permission: &types.AllowNobody,
			expErr:     true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			xCtx, _ := ctx.CacheContext()
			// setup
			_, _, sender := testdata.KeyTestPubAddr()
			msg := types.MsgStoreCodeFixture(func(m *types.MsgStoreCode) {
				m.WASMByteCode = wasmContract
				m.Sender = sender.String()
				m.InstantiatePermission = spec.permission
			})

			// store code
			rsp, err := wasmApp.MsgServiceRouter().Handler(msg)(xCtx, msg)
			require.NoError(t, err)
			var result types.MsgStoreCodeResponse
			require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &result))

			// when
			msgInstantiate := &types.MsgInstantiateContract{
				Sender: spec.addr,
				Admin:  myAddress.String(),
				CodeID: result.CodeID,
				Label:  "test",
				Msg:    []byte(`{}`),
				Funds:  sdk.Coins{},
			}
			rsp, err = wasmApp.MsgServiceRouter().Handler(msgInstantiate)(xCtx, msgInstantiate)

			//then
			if spec.expErr {
				require.Error(t, err)
				return
			}

			// check event
			events := rsp.Events
			assert.Equal(t, 2, len(events))
			assert.Equal(t, "message", events[0].Type)
			assert.Equal(t, 2, len(events[0].Attributes))
			assert.Equal(t, "module", string(events[0].Attributes[0].Key))
			assert.Equal(t, "wasm", string(events[0].Attributes[0].Value))
			assert.Equal(t, "sender", string(events[0].Attributes[1].Key))
			assert.Equal(t, myAddress.String(), string(events[0].Attributes[1].Value))
			assert.Equal(t, "instantiate", events[1].Type)
			assert.Equal(t, 2, len(events[1].Attributes))
			assert.Equal(t, "_contract_address", string(events[1].Attributes[0].Key))
			assert.Contains(t, string(rsp.Data), string(events[1].Attributes[0].Value))
			assert.Equal(t, "code_id", string(events[1].Attributes[1].Key))
			assert.Equal(t, "1", string(events[1].Attributes[1].Value))

			require.NoError(t, err)
		})
	}
}
