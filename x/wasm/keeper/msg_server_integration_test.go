package keeper_test

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"

	"github.com/Finschia/wasmd/app"
	"github.com/Finschia/wasmd/x/wasm/keeper"
	"github.com/Finschia/wasmd/x/wasm/types"
)

//go:embed testdata/reflect.wasm
var wasmContract []byte

//go:embed testdata/hackatom.wasm
var hackatomContract []byte

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

func TestInstantiateContract2(t *testing.T) {
	wasmApp := app.Setup(false)
	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	var (
		myAddress sdk.AccAddress = make([]byte, types.ContractAddrLen)
	)

	specs := map[string]struct {
		addr       string
		permission *types.AccessConfig
		salt       string
		expErr     bool
	}{
		"address can instantiate a contract when permission is everybody": {
			addr:       myAddress.String(),
			permission: &types.AllowEverybody,
			salt:       "salt1",
			expErr:     false,
		},
		"address cannot instantiate a contract when permission is nobody": {
			addr:       myAddress.String(),
			permission: &types.AllowNobody,
			salt:       "salt2",
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
			msgInstantiate := &types.MsgInstantiateContract2{
				Sender: spec.addr,
				Admin:  myAddress.String(),
				CodeID: result.CodeID,
				Label:  "test",
				Msg:    []byte(`{}`),
				Funds:  sdk.Coins{},
				Salt:   []byte(spec.salt),
				FixMsg: true,
			}
			rsp, err = wasmApp.MsgServiceRouter().Handler(msgInstantiate)(xCtx, msgInstantiate)

			//then
			if spec.expErr {
				require.Error(t, err)
				return
			}

			var instantiateResponse types.MsgInstantiateContractResponse
			require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &instantiateResponse))

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
			assert.Equal(t, instantiateResponse.Address, string(events[1].Attributes[0].Value))
			assert.Equal(t, "code_id", string(events[1].Attributes[1].Key))
			assert.Equal(t, "1", string(events[1].Attributes[1].Value))

			require.NoError(t, err)
		})
	}
}

func TestMigrateContract(t *testing.T) {
	wasmApp := app.Setup(false)
	ctx := wasmApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	var (
		myAddress       sdk.AccAddress = make([]byte, types.ContractAddrLen)
		_, _, otherAddr                = testdata.KeyTestPubAddr()
	)

	specs := map[string]struct {
		addr   string
		expErr bool
	}{
		"admin can migrate a contract": {
			addr:   myAddress.String(),
			expErr: false,
		},
		"other address cannot migrate a contract": {
			addr:   otherAddr.String(),
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			xCtx, _ := ctx.CacheContext()
			// setup
			_, _, sender := testdata.KeyTestPubAddr()
			msg := types.MsgStoreCodeFixture(func(m *types.MsgStoreCode) {
				m.WASMByteCode = hackatomContract
				m.Sender = sender.String()
			})

			// store code
			rsp, err := wasmApp.MsgServiceRouter().Handler(msg)(xCtx, msg)
			require.NoError(t, err)
			var storeCodeResponse types.MsgStoreCodeResponse
			require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &storeCodeResponse))

			// instantiate contract
			initMsg := keeper.HackatomExampleInitMsg{
				Verifier:    sender,
				Beneficiary: myAddress,
			}
			initMsgBz, err := json.Marshal(initMsg)
			require.NoError(t, err)

			msgInstantiate := &types.MsgInstantiateContract{
				Sender: sender.String(),
				Admin:  myAddress.String(),
				CodeID: storeCodeResponse.CodeID,
				Label:  "test",
				Msg:    initMsgBz,
				Funds:  sdk.Coins{},
			}
			rsp, err = wasmApp.MsgServiceRouter().Handler(msgInstantiate)(xCtx, msgInstantiate)
			require.NoError(t, err)
			var instantiateResponse types.MsgInstantiateContractResponse
			require.NoError(t, wasmApp.AppCodec().Unmarshal(rsp.Data, &instantiateResponse))

			// when
			migMsg := struct {
				Verifier sdk.AccAddress `json:"verifier"`
			}{Verifier: myAddress}
			migMsgBz, err := json.Marshal(migMsg)
			require.NoError(t, err)
			msgMigrateContract := &types.MsgMigrateContract{
				Sender:   spec.addr,
				Msg:      migMsgBz,
				Contract: instantiateResponse.Address,
				CodeID:   storeCodeResponse.CodeID,
			}
			rsp, err = wasmApp.MsgServiceRouter().Handler(msgMigrateContract)(xCtx, msgMigrateContract)

			// then
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
			assert.Equal(t, "migrate", events[1].Type)
			assert.Equal(t, 2, len(events[1].Attributes))
			assert.Equal(t, "code_id", string(events[1].Attributes[0].Key))
			assert.Equal(t, "1", string(events[1].Attributes[0].Value))
			assert.Equal(t, "_contract_address", string(events[1].Attributes[1].Key))
			assert.Equal(t, instantiateResponse.Address, string(events[1].Attributes[1].Value))

			require.NoError(t, err)
		})
	}
}
