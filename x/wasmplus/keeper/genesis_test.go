package keeper

import (
	"crypto/sha256"
	"os"
	"testing"
	"time"

	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/wasmd/x/wasm/keeper"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	wasmTypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

const (
	firstCodeID  = 1
	humanAddress = "link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8"

	AvailableCapabilities = "iterator,staking,stargate,cosmwasm_1_1"
)

func TestGenesisExportImport(t *testing.T) {
	wasmKeeper, srcCtx := setupKeeper(t)
	contractKeeper := NewPermissionedKeeper(*wasmkeeper.NewGovPermissionKeeper(wasmKeeper), wasmKeeper)

	wasmCode, err := os.ReadFile("../../wasm/keeper/testdata/hackatom.wasm")
	require.NoError(t, err)

	// store some test data
	f := fuzz.New().Funcs(wasmkeeper.ModelFuzzers...)

	wasmKeeper.SetParams(srcCtx, wasmTypes.DefaultParams())

	for i := 0; i < 5; i++ {
		var (
			codeInfo          wasmTypes.CodeInfo
			contract          wasmTypes.ContractInfo
			stateModels       []wasmTypes.Model
			history           []wasmTypes.ContractCodeHistoryEntry
			pinned            bool
			contractExtension bool
			verifier          sdk.AccAddress
			beneficiary       sdk.AccAddress
		)
		f.Fuzz(&codeInfo)
		f.Fuzz(&contract)
		f.Fuzz(&stateModels)
		f.NilChance(0).Fuzz(&history)
		f.Fuzz(&pinned)
		f.Fuzz(&contractExtension)
		f.Fuzz(&verifier)
		f.Fuzz(&beneficiary)

		creatorAddr, err := sdk.AccAddressFromBech32(codeInfo.Creator)
		require.NoError(t, err)
		codeID, _, err := contractKeeper.Create(srcCtx, creatorAddr, wasmCode, &codeInfo.InstantiateConfig)
		require.NoError(t, err)
		if pinned {
			err = contractKeeper.PinCode(srcCtx, codeID)
			require.NoError(t, err)
		}
		if contractExtension {
			anyTime := time.Now().UTC()
			var nestedType v1beta1.TextProposal
			f.NilChance(0).Fuzz(&nestedType)
			myExtension, err := v1beta1.NewProposal(&nestedType, 1, anyTime, anyTime)
			require.NoError(t, err)
			err = contract.SetExtension(&myExtension)
			require.NoError(t, err)
		}

		initMsgBz := HackatomExampleInitMsg{
			Verifier:    verifier,
			Beneficiary: beneficiary,
		}.GetBytes(t)

		_, _, err = contractKeeper.Instantiate(srcCtx, codeID, creatorAddr, creatorAddr, initMsgBz, "test", nil)
		require.NoError(t, err)
	}
	var wasmParams wasmTypes.Params
	f.NilChance(0).Fuzz(&wasmParams)
	wasmKeeper.SetParams(srcCtx, wasmParams)

	// add inactiveContractAddr
	var inactiveContractAddr []sdk.AccAddress
	wasmKeeper.IterateContractInfo(srcCtx, func(address sdk.AccAddress, info wasmTypes.ContractInfo) bool {
		err = contractKeeper.DeactivateContract(srcCtx, address)
		require.NoError(t, err)
		inactiveContractAddr = append(inactiveContractAddr, address)
		return false
	})

	// export
	exportedState := ExportGenesis(srcCtx, wasmKeeper)
	exportedGenesis, err := wasmKeeper.cdc.MarshalJSON(exportedState)
	require.NoError(t, err)

	// setup new instances
	dstKeeper, dstCtx := setupKeeper(t)

	// re-import
	var importState types.GenesisState
	err = dstKeeper.cdc.UnmarshalJSON(exportedGenesis, &importState)
	require.NoError(t, err)
	_, err = InitGenesis(dstCtx, dstKeeper, importState)
	require.NoError(t, err)

	// compare
	dstParams := dstKeeper.GetParams(dstCtx)
	require.Equal(t, wasmParams, dstParams)

	var destInactiveContractAddr []sdk.AccAddress
	dstKeeper.IterateInactiveContracts(dstCtx, func(contractAddress sdk.AccAddress) (stop bool) {
		destInactiveContractAddr = append(destInactiveContractAddr, contractAddress)
		return false
	})
	require.Equal(t, inactiveContractAddr, destInactiveContractAddr)
}

func TestGenesisInit(t *testing.T) {
	wasmCode, err := os.ReadFile("../../wasm/keeper/testdata/hackatom.wasm")
	require.NoError(t, err)

	myCodeInfo := wasmTypes.CodeInfoFixture(wasmTypes.WithSHA256CodeHash(wasmCode))
	specs := map[string]struct {
		src            types.GenesisState
		stakingMock    StakingKeeperMock
		msgHandlerMock MockMsgHandler
		expSuccess     bool
	}{
		"happy path: code info correct": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params: wasmTypes.DefaultParams(),
			},
			expSuccess: true,
		},
		"happy path: code ids can contain gaps": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}, {
					CodeID:    3,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 10},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params: wasmTypes.DefaultParams(),
			},
			expSuccess: true,
		},
		"happy path: code order does not matter": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    2,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}, {
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: nil,
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 3},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params: wasmTypes.DefaultParams(),
			},
			expSuccess: true,
		},
		"prevent code hash mismatch": {src: types.GenesisState{
			Codes: []wasmTypes.Code{{
				CodeID:    firstCodeID,
				CodeInfo:  wasmTypes.CodeInfoFixture(func(i *wasmTypes.CodeInfo) { i.CodeHash = make([]byte, sha256.Size) }),
				CodeBytes: wasmCode,
			}},
			Params: wasmTypes.DefaultParams(),
		}},
		"prevent duplicate codeIDs": {src: types.GenesisState{
			Codes: []wasmTypes.Code{
				{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				},
				{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				},
			},
			Params: wasmTypes.DefaultParams(),
		}},
		"codes with same checksum can be pinned": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{
					{
						CodeID:    firstCodeID,
						CodeInfo:  myCodeInfo,
						CodeBytes: wasmCode,
						Pinned:    true,
					},
					{
						CodeID:    2,
						CodeInfo:  myCodeInfo,
						CodeBytes: wasmCode,
						Pinned:    true,
					},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"happy path: code id in info and contract do match": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 2},
				},
				Params: wasmTypes.DefaultParams(),
			},
			expSuccess: true,
		},
		"happy path: code info with two contracts": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					}, {
						ContractAddress: keeper.BuildContractAddressClassic(1, 2).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 3},
				},
				Params: wasmTypes.DefaultParams(),
			},
			expSuccess: true,
		},
		"prevent contracts that points to non existing codeID": {
			src: types.GenesisState{
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"prevent duplicate contract address": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					}, {
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"prevent duplicate contract model keys": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
						ContractState: []wasmTypes.Model{
							{
								Key:   []byte{0x1},
								Value: []byte("foo"),
							},
							{
								Key:   []byte{0x1},
								Value: []byte("bar"),
							},
						},
					},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"prevent duplicate sequences": {
			src: types.GenesisState{
				Sequences: []wasmTypes.Sequence{
					{IDKey: []byte("foo"), Value: 1},
					{IDKey: []byte("foo"), Value: 9999},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"prevent code id seq init value == max codeID used": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    2,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 1},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"prevent contract id seq init value == count contracts": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params: wasmTypes.DefaultParams(),
			},
		},
		"happy path: inactiveContract": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Contracts: []wasmTypes.Contract{
					{
						ContractAddress: keeper.BuildContractAddressClassic(1, 1).String(),
						ContractInfo:    wasmTypes.ContractInfoFixture(func(c *wasmTypes.ContractInfo) { c.CodeID = 1 }, wasmTypes.OnlyGenesisFields),
					},
				},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 3},
				},
				Params:                    wasmTypes.DefaultParams(),
				InactiveContractAddresses: []string{keeper.BuildContractAddressClassic(1, 1).String()},
			},
			expSuccess: true,
		},
		"invalid path: inactiveContract - human address": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params:                    wasmTypes.DefaultParams(),
				InactiveContractAddresses: []string{humanAddress},
			},
		},
		"invalid path: inactiveContract - do not imported": {
			src: types.GenesisState{
				Codes: []wasmTypes.Code{{
					CodeID:    firstCodeID,
					CodeInfo:  myCodeInfo,
					CodeBytes: wasmCode,
				}},
				Sequences: []wasmTypes.Sequence{
					{IDKey: wasmTypes.KeySequenceCodeID, Value: 2},
					{IDKey: wasmTypes.KeySequenceInstanceID, Value: 1},
				},
				Params:                    wasmTypes.DefaultParams(),
				InactiveContractAddresses: []string{keeper.BuildContractAddressClassic(1, 1).String()},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			require.NoError(t, types.ValidateGenesis(spec.src))
			gotValidatorSet, gotErr := InitGenesis(ctx, keeper, spec.src)
			if !spec.expSuccess {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			spec.msgHandlerMock.verifyCalls(t)
			spec.stakingMock.verifyCalls(t)
			assert.Equal(t, spec.stakingMock.validatorUpdate, gotValidatorSet)
			for _, c := range spec.src.Codes {
				assert.Equal(t, c.Pinned, keeper.IsPinnedCode(ctx, c.CodeID))
			}
		})
	}
}

func setupKeeper(t *testing.T) (*Keeper, sdk.Context) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "wasm")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	keyWasm := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())
	ms.MountStoreWithDB(keyWasm, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, cmtproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, false, log.NewNopLogger())

	encodingConfig := MakeEncodingConfig(t)
	// register an example extension. must be protobuf
	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*wasmTypes.ContractInfoExtension)(nil),
		&v1beta1.Proposal{},
	)
	// also registering gov interfaces for nested Any type
	v1beta1.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	wasmConfig := wasmTypes.DefaultWasmConfig()

	srcKeeper := NewKeeper(
		encodingConfig.Codec,
		runtime.NewKVStoreService(keyWasm),
		authkeeper.AccountKeeper{},
		&bankkeeper.BaseKeeper{},
		stakingkeeper.Keeper{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		tempDir,
		wasmConfig,
		AvailableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	return &srcKeeper, ctx
}

type StakingKeeperMock struct {
	err             error
	validatorUpdate []abci.ValidatorUpdate
	expCalls        int
	gotCalls        int
}

func (s *StakingKeeperMock) ApplyAndReturnValidatorSetUpdates(_ sdk.Context) ([]abci.ValidatorUpdate, error) {
	s.gotCalls++
	return s.validatorUpdate, s.err
}

func (s *StakingKeeperMock) verifyCalls(t *testing.T) {
	assert.Equal(t, s.expCalls, s.gotCalls, "number calls")
}

type MockMsgHandler struct {
	result   *sdk.Result
	err      error
	expCalls int
	gotCalls int
	expMsg   sdk.Msg
	gotMsg   sdk.Msg
}

func (m *MockMsgHandler) Handle(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	m.gotCalls++
	m.gotMsg = msg
	return m.result, m.err
}

func (m *MockMsgHandler) verifyCalls(t *testing.T) {
	assert.Equal(t, m.expMsg, m.gotMsg, "message param")
	assert.Equal(t, m.expCalls, m.gotCalls, "number calls")
}
