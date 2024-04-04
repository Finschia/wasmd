package wasmplus

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/Finschia/wasmd/app/params"
	"github.com/Finschia/wasmd/x/wasm/exported"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	v2 "github.com/Finschia/wasmd/x/wasm/migrations/v2"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/keeper"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

type mockSubspace struct {
	ps v2.Params
}

func newMockSubspace(ps v2.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*v2.Params) = ms.ps
}

type testData struct {
	module           AppModule
	ctx              sdk.Context
	acctKeeper       authkeeper.AccountKeeper
	keeper           keeper.Keeper
	bankKeeper       bankkeeper.Keeper
	stakingKeeper    *stakingkeeper.Keeper
	faucet           *wasmkeeper.TestFaucet
	grpcQueryRouter  *baseapp.GRPCQueryRouter
	msgServiceRouter *baseapp.MsgServiceRouter
	encConf          params.EncodingConfig
}

func setupTest(t *testing.T) testData {
	t.Helper()
	DefaultParams := v2.Params{
		CodeUploadAccess:             v2.AccessConfig{Permission: v2.AccessTypeEverybody},
		InstantiateDefaultPermission: v2.AccessTypeEverybody,
	}

	ctx, keepers := keeper.CreateTestInput(t, false, "iterator,staking,stargate,cosmwasm_1_1")
	encConf := keeper.MakeEncodingConfig(t)
	queryRouter := baseapp.NewGRPCQueryRouter()
	serviceRouter := baseapp.NewMsgServiceRouter()
	queryRouter.SetInterfaceRegistry(encConf.InterfaceRegistry)
	serviceRouter.SetInterfaceRegistry(encConf.InterfaceRegistry)
	data := testData{
		module:           NewAppModule(encConf.Codec, keepers.WasmKeeper, keepers.StakingKeeper, keepers.AccountKeeper, keepers.BankKeeper, nil, newMockSubspace(DefaultParams)),
		ctx:              ctx,
		acctKeeper:       keepers.AccountKeeper,
		keeper:           *keepers.WasmKeeper,
		bankKeeper:       keepers.BankKeeper,
		stakingKeeper:    keepers.StakingKeeper,
		faucet:           keepers.Faucet,
		grpcQueryRouter:  queryRouter,
		msgServiceRouter: serviceRouter,
		encConf:          encConf,
	}
	data.module.RegisterServices(module.NewConfigurator(encConf.Codec, serviceRouter, queryRouter))
	return data
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mustLoad(path string) []byte {
	bz, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bz
}

var (
	_, _, addrAcc1 = keyPubAddr()
	addr1          = addrAcc1.String()
	testContract   = mustLoad("../wasm/keeper/testdata/hackatom.wasm")
	oldContract    = mustLoad("../wasm/testdata/escrow_0.7.wasm")
)

type initMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}

type emptyMsg struct{}

type state struct {
	Verifier    string `json:"verifier"`
	Beneficiary string `json:"beneficiary"`
	Funder      string `json:"funder"`
}

// ensures this returns a valid codeID and bech32 address and returns it
func parseStoreAndInitResponse(t *testing.T, data []byte) (uint64, string) {
	var res types.MsgStoreCodeAndInstantiateContractResponse
	require.NoError(t, res.Unmarshal(data))
	require.NotEmpty(t, res.CodeID)
	require.NotEmpty(t, res.Address)
	addr := res.Address
	codeID := res.CodeID
	// ensure this is a valid sdk address
	_, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)
	return codeID, addr
}

// ensure store code returns the expected response
func assertStoreCodeResponse(t *testing.T, data []byte, expected uint64) {
	var pStoreResp wasmtypes.MsgStoreCodeResponse
	require.NoError(t, pStoreResp.Unmarshal(data))
	require.Equal(t, pStoreResp.CodeID, expected)
}

type prettyEvent struct {
	Type string
	Attr []sdk.Attribute
}

func prettyEvents(evts []abci.Event) string {
	res := make([]prettyEvent, len(evts))
	for i, e := range evts {
		res[i] = prettyEvent{
			Type: e.Type,
			Attr: prettyAttrs(e.Attributes),
		}
	}
	bz, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func prettyAttrs(attrs []abci.EventAttribute) []sdk.Attribute {
	pretty := make([]sdk.Attribute, len(attrs))
	for i, a := range attrs {
		pretty[i] = prettyAttr(a)
	}
	return pretty
}

func prettyAttr(attr abci.EventAttribute) sdk.Attribute {
	return sdk.NewAttribute(string(attr.Key), string(attr.Value))
}

func assertAttribute(t *testing.T, key string, value string, attr abci.EventAttribute) {
	t.Helper()
	assert.Equal(t, key, string(attr.Key), prettyAttr(attr))
	assert.Equal(t, value, string(attr.Value), prettyAttr(attr))
}

func assertCodeList(t *testing.T, q *baseapp.GRPCQueryRouter, ctx sdk.Context, expectedNum int) {
	t.Helper()
	path := "/cosmwasm.wasm.v1.Query/Codes"
	resp, sdkerr := q.Route(path)(ctx, &abci.RequestQuery{Path: path})
	require.NoError(t, sdkerr)
	require.True(t, resp.IsOK())

	bz := resp.Value
	if len(bz) == 0 {
		require.Equal(t, expectedNum, 0)
		return
	}

	var res []wasmtypes.CodeInfo
	err := json.Unmarshal(bz, &res)
	require.NoError(t, err)

	assert.Equal(t, expectedNum, len(res))
}

func assertCodeBytes(t *testing.T, q *baseapp.GRPCQueryRouter, ctx sdk.Context, codeID uint64, expectedBytes []byte, marshaler codec.Codec) {
	t.Helper()
	bz, err := marshaler.Marshal(&wasmtypes.QueryCodeRequest{CodeId: codeID})
	require.NoError(t, err)

	path := "/cosmwasm.wasm.v1.Query/Code"
	resp, err := q.Route(path)(ctx, &abci.RequestQuery{Path: path, Data: bz})
	if len(expectedBytes) == 0 {
		require.Equal(t, wasmtypes.ErrNoSuchCodeFn(codeID).Wrapf("code id %d", codeID).Error(), err.Error())
		return
	}
	require.NoError(t, err)
	require.True(t, resp.IsOK())
	bz = resp.Value

	var rsp wasmtypes.QueryCodeResponse
	require.NoError(t, marshaler.Unmarshal(bz, &rsp))
	assert.Equal(t, expectedBytes, rsp.Data)
}

func assertContractList(t *testing.T, q *baseapp.GRPCQueryRouter, ctx sdk.Context, codeID uint64, expContractAddrs []string, marshaler codec.Codec) { //nolint:unparam
	t.Helper()
	bz, err := marshaler.Marshal(&wasmtypes.QueryContractsByCodeRequest{CodeId: codeID})
	require.NoError(t, err)

	path := "/cosmwasm.wasm.v1.Query/ContractsByCode"
	resp, sdkerr := q.Route(path)(ctx, &abci.RequestQuery{Path: path, Data: bz})
	if len(expContractAddrs) == 0 {
		assert.ErrorIs(t, err, wasmtypes.ErrNotFound)
		return
	}
	require.NoError(t, sdkerr)
	require.True(t, resp.IsOK())
	bz = resp.Value

	var rsp wasmtypes.QueryContractsByCodeResponse
	require.NoError(t, marshaler.Unmarshal(bz, &rsp))

	hasAddrs := make([]string, len(rsp.Contracts))
	for i, r := range rsp.Contracts { //nolint:gosimple
		hasAddrs[i] = r
	}
	assert.Equal(t, expContractAddrs, hasAddrs)
}

func assertContractInfo(t *testing.T, q *baseapp.GRPCQueryRouter, ctx sdk.Context, contractBech32Addr string, codeID uint64, creator sdk.AccAddress, marshaler codec.Codec) { //nolint:unparam
	t.Helper()
	bz, err := marshaler.Marshal(&wasmtypes.QueryContractInfoRequest{Address: contractBech32Addr})
	require.NoError(t, err)

	path := "/cosmwasm.wasm.v1.Query/ContractInfo"
	resp, sdkerr := q.Route(path)(ctx, &abci.RequestQuery{Path: path, Data: bz})
	require.NoError(t, sdkerr)
	require.True(t, resp.IsOK())
	bz = resp.Value

	var rsp wasmtypes.QueryContractInfoResponse
	require.NoError(t, marshaler.Unmarshal(bz, &rsp))

	assert.Equal(t, codeID, rsp.CodeID)
	assert.Equal(t, creator.String(), rsp.Creator)
}

func assertContractState(t *testing.T, q *baseapp.GRPCQueryRouter, ctx sdk.Context, contractBech32Addr string, expected state, marshaler codec.Codec) {
	t.Helper()
	bz, err := marshaler.Marshal(&wasmtypes.QueryRawContractStateRequest{Address: contractBech32Addr, QueryData: []byte("config")})
	require.NoError(t, err)

	path := "/cosmwasm.wasm.v1.Query/RawContractState"
	resp, sdkerr := q.Route(path)(ctx, &abci.RequestQuery{Path: path, Data: bz})
	require.NoError(t, sdkerr)
	require.True(t, resp.IsOK())
	bz = resp.Value

	var rsp wasmtypes.QueryRawContractStateResponse
	require.NoError(t, marshaler.Unmarshal(bz, &rsp))
	expectedBz, err := json.Marshal(expected)
	require.NoError(t, err)
	assert.Equal(t, expectedBz, rsp.Data)
}

func TestHandleStoreAndInstantiate(t *testing.T) {
	data := setupTest(t)
	creator := data.faucet.NewFundedRandomAccount(data.ctx, sdk.NewInt64Coin("denom", 100000))

	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()

	initPayload := initMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initPayload)
	require.NoError(t, err)

	// create with no balance is legal
	msg := &types.MsgStoreCodeAndInstantiateContract{
		Sender:       creator.String(),
		WASMByteCode: testContract,
		Msg:          initMsgBz,
		Label:        "contract for test",
		Funds:        nil,
	}
	h := data.msgServiceRouter.Handler(msg)
	q := data.grpcQueryRouter
	res, err := h(data.ctx, msg)
	require.NoError(t, err)
	codeID, contractBech32Addr := parseStoreAndInitResponse(t, res.Data)

	require.Equal(t, uint64(1), codeID)
	require.Equal(t, "link14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sgf2vn8", contractBech32Addr)
	// this should be standard x/wasm init event, nothing from contract
	require.Equal(t, 3, len(res.Events), prettyEvents(res.Events))
	assert.Equal(t, "store_code", res.Events[0].Type)
	assertAttribute(t, "code_id", "1", res.Events[0].Attributes[1])
	assert.Equal(t, "instantiate", res.Events[1].Type)
	assertAttribute(t, "_contract_address", contractBech32Addr, res.Events[1].Attributes[0])
	assertAttribute(t, "code_id", "1", res.Events[1].Attributes[1])
	assert.Equal(t, "wasm", res.Events[2].Type)
	assertAttribute(t, "_contract_address", contractBech32Addr, res.Events[2].Attributes[0])

	assertCodeList(t, q, data.ctx, 1)
	assertCodeBytes(t, q, data.ctx, 1, testContract, data.encConf.Codec)

	assertContractList(t, q, data.ctx, 1, []string{contractBech32Addr}, data.encConf.Codec)
	assertContractInfo(t, q, data.ctx, contractBech32Addr, 1, creator, data.encConf.Codec)
	assertContractState(t, q, data.ctx, contractBech32Addr, state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	}, data.encConf.Codec)
}

func TestErrorsCreateAndInstantiate(t *testing.T) {
	// init messages
	_, _, bob := keyPubAddr()
	_, _, fred := keyPubAddr()
	initMsg := initMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	validInitMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	invalidInitMsgBz, err := json.Marshal(emptyMsg{})

	expectedContractBech32Addr := "link14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sgf2vn8"

	// test cases
	cases := map[string]struct {
		msg           sdk.Msg
		isValid       bool
		expectedCodes int
		expectedBytes []byte
	}{
		"empty": {
			msg:           &types.MsgStoreCodeAndInstantiateContract{},
			isValid:       false,
			expectedCodes: 0,
			expectedBytes: nil,
		},
		"valid one": {
			msg: &types.MsgStoreCodeAndInstantiateContract{
				Sender:       addr1,
				WASMByteCode: testContract,
				Msg:          validInitMsgBz,
				Label:        "foo",
				Funds:        nil,
			},
			isValid:       true,
			expectedCodes: 1,
			expectedBytes: testContract,
		},
		"invalid wasm": {
			msg: &types.MsgStoreCodeAndInstantiateContract{
				Sender:       addr1,
				WASMByteCode: []byte("foobar"),
				Msg:          validInitMsgBz,
				Label:        "foo",
				Funds:        nil,
			},
			isValid:       false,
			expectedCodes: 0,
			expectedBytes: nil,
		},
		"old wasm (0.7)": {
			msg: &types.MsgStoreCodeAndInstantiateContract{
				Sender:       addr1,
				WASMByteCode: oldContract,
				Msg:          validInitMsgBz,
				Label:        "foo",
				Funds:        nil,
			},
			isValid:       false,
			expectedCodes: 0,
			expectedBytes: nil,
		},
		"invalid init message": {
			msg: &types.MsgStoreCodeAndInstantiateContract{
				Sender:       addr1,
				WASMByteCode: testContract,
				Msg:          invalidInitMsgBz,
				Label:        "foo",
				Funds:        nil,
			},
			isValid:       false,
			expectedCodes: 1,
			expectedBytes: testContract,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			data := setupTest(t)

			h := data.msgServiceRouter.Handler(tc.msg)
			q := data.grpcQueryRouter

			// asserting response
			res, err := h(data.ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				codeID, contractBech32Addr := parseStoreAndInitResponse(t, res.Data)
				require.Equal(t, uint64(1), codeID)
				require.Equal(t, expectedContractBech32Addr, contractBech32Addr)

			} else {
				require.Error(t, err, "%#v", res)
			}

			// asserting code state
			assertCodeList(t, q, data.ctx, tc.expectedCodes)
			assertCodeBytes(t, q, data.ctx, 1, tc.expectedBytes, data.encConf.Codec)

			// asserting contract state
			if tc.isValid {
				assertContractList(t, q, data.ctx, 1, []string{expectedContractBech32Addr}, data.encConf.Codec)
				assertContractInfo(t, q, data.ctx, expectedContractBech32Addr, 1, addrAcc1, data.encConf.Codec)
				assertContractState(t, q, data.ctx, expectedContractBech32Addr, state{
					Verifier:    fred.String(),
					Beneficiary: bob.String(),
					Funder:      addrAcc1.String(),
				}, data.encConf.Codec)
			} else {
				assertContractList(t, q, data.ctx, 0, []string{}, data.encConf.Codec)
			}
		})
	}
}

func TestHandleNonPlusWasmCreate(t *testing.T) {
	data := setupTest(t)
	creator := data.faucet.NewFundedRandomAccount(data.ctx, sdk.NewInt64Coin("denom", 100000))

	msg := &wasmtypes.MsgStoreCode{
		Sender:       creator.String(),
		WASMByteCode: testContract,
	}

	h := data.msgServiceRouter.Handler(msg)
	q := data.grpcQueryRouter

	res, err := h(data.ctx, msg)
	require.NoError(t, err)
	assertStoreCodeResponse(t, res.Data, 1)

	require.Equal(t, 1, len(res.Events), prettyEvents(res.Events))
	assert.Equal(t, "store_code", res.Events[0].Type)
	assertAttribute(t, "code_id", "1", res.Events[0].Attributes[1])

	assertCodeList(t, q, data.ctx, 1)
	assertCodeBytes(t, q, data.ctx, 1, testContract, data.encConf.Codec)
}

func TestErrorHandleNonPlusWasmCreate(t *testing.T) {
	data := setupTest(t)
	creator := data.faucet.NewFundedRandomAccount(data.ctx, sdk.NewInt64Coin("denom", 100000))

	msg := &wasmtypes.MsgStoreCode{
		Sender:       creator.String(),
		WASMByteCode: []byte("invalid WASM contract"),
	}

	h := data.msgServiceRouter.Handler(msg)

	_, err := h(data.ctx, msg)
	require.ErrorContains(t, err, "Wasm validation")
}
