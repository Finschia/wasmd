package keeper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	"github.com/Finschia/wasmd/x/wasmplus/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmappparams "github.com/Finschia/wasmd/app/params"
	"github.com/Finschia/wasmd/x/wasm/keeper/wasmtesting"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	wasmplustypes "github.com/Finschia/wasmd/x/wasmplus/types"
)

var moduleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
	}),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
)

func MakeTestCodec(t testing.TB) codec.Codec {
	return MakeEncodingConfig(t).Codec
}

func MakeEncodingConfig(_ testing.TB) wasmappparams.EncodingConfig {
	encodingConfig := wasmappparams.MakeEncodingConfig()
	amino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	moduleBasics.RegisterLegacyAminoCodec(amino)
	moduleBasics.RegisterInterfaces(interfaceRegistry)
	// add wasmplusd types
	wasmplustypes.RegisterInterfaces(interfaceRegistry)
	wasmplustypes.RegisterLegacyAminoCodec(amino)

	return encodingConfig
}

var TestingStakeParams = stakingtypes.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	HistoricalEntries: 10,
	BondDenom:         "stake",
}

type TestKeepers struct {
	AccountKeeper  authkeeper.AccountKeeper
	StakingKeeper  *stakingkeeper.Keeper
	DistKeeper     distributionkeeper.Keeper
	BankKeeper     bankkeeper.Keeper
	GovKeeper      *govkeeper.Keeper
	ContractKeeper types.ContractOpsKeeper
	WasmKeeper     *Keeper
	IBCKeeper      *ibckeeper.Keeper
	Router         wasmkeeper.MessageRouter
	EncodingConfig wasmappparams.EncodingConfig
	Faucet         *wasmkeeper.TestFaucet
	MultiStore     storetypes.CommitMultiStore
}

// CreateDefaultTestInput common settings for CreateTestInput
func CreateDefaultTestInput(t testing.TB) (sdk.Context, TestKeepers) {
	return CreateTestInput(t, false, "staking")
}

// CreateTestInput encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func CreateTestInput(
	t testing.TB,
	isCheckTx bool,
	availableCapabilities string,
	opts ...wasmkeeper.Option,
) (sdk.Context, TestKeepers) {
	// Load default wasm config
	return createTestInput(t, isCheckTx, availableCapabilities, wasmtypes.DefaultWasmConfig(), dbm.NewMemDB(), opts...)
}

// encoders can be nil to accept the defaults, or set it to override some of the message handlers (like default)
func createTestInput(
	t testing.TB,
	isCheckTx bool,
	availableCapabilities string,
	wasmConfig wasmtypes.WasmConfig,
	db dbm.DB,
	opts ...wasmkeeper.Option,
) (sdk.Context, TestKeepers) {
	tempDir := t.TempDir()

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distributiontypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		types.StoreKey,
	)
	logger := log.NewTestLogger(t)

	ms := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	for _, v := range keys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeIAVL, db)
	}
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeTransient, db)
	}

	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeMemory, db)
	}

	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.NewNopLogger())
	ctx = wasmtypes.WithTXCounter(ctx, 0)

	encodingConfig := MakeEncodingConfig(t)
	appCodec, legacyAmino := encodingConfig.Codec, encodingConfig.Amino

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	for _, m := range []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		minttypes.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		ibcexported.ModuleName,
		govtypes.ModuleName,
		types.ModuleName,
	} {
		paramsKeeper.Subspace(m)
	}
	subspace := func(m string) paramstypes.Subspace {
		r, ok := paramsKeeper.GetSubspace(m)
		require.True(t, ok)
		return r
	}
	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName:     nil,
		distributiontypes.ModuleName:   nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		types.ModuleName:               {authtypes.Burner},
	}
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	require.NoError(t, accountKeeper.Params.Set(ctx, authtypes.DefaultParams()))

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(banktypes.ModuleName).String(),
		logger,
	)
	require.NoError(t, bankKeeper.SetParams(ctx, banktypes.DefaultParams()))

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(stakingtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	stakingtypes.DefaultParams()
	require.NoError(t, stakingKeeper.SetParams(ctx, TestingStakeParams))

	distKeeper := distributionkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distributiontypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(distributiontypes.ModuleName).String(),
	)
	require.NoError(t, distKeeper.Params.Set(ctx, distributiontypes.DefaultParams()))
	require.NoError(t, distKeeper.FeePool.Set(ctx, distributiontypes.InitialFeePool()))
	stakingKeeper.SetHooks(distKeeper.Hooks())

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		tempDir,
		nil,
		authtypes.NewModuleAddress(upgradetypes.ModuleName).String(),
	)

	faucet := wasmkeeper.NewTestFaucet(t, ctx, bankKeeper, minttypes.ModuleName, sdk.NewCoin("stake", sdkmath.NewInt(100_000_000_000)))

	// set some funds ot pay out validatores, based on code from:
	// https://github.com/Finschia/finschia-sdk/blob/95b22d3a685f7eb531198e0023ef06873835e632/x/distribution/keeper/keeper_test.go#L49-L56
	distrAcc := distKeeper.GetDistributionAccount(ctx)
	faucet.Fund(ctx, distrAcc.GetAddress(), sdk.NewCoin("stake", sdkmath.NewInt(2000000)))
	accountKeeper.SetModuleAccount(ctx, distrAcc)

	capabilityKeeper := capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedWasmKeeper := capabilityKeeper.ScopeToModule(types.ModuleName)

	ibcKeeper := ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		subspace(ibcexported.ModuleName),
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	querier := baseapp.NewGRPCQueryRouter()
	querier.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())

	keeper := NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[types.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distributionkeeper.NewQuerier(distKeeper),
		ibcKeeper.ChannelKeeper,
		ibcKeeper.ChannelKeeper,
		ibcKeeper.PortKeeper,
		scopedWasmKeeper,
		wasmtesting.MockIBCTransferKeeper{},
		msgRouter,
		querier,
		tempDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		opts...,
	)
	require.NoError(t, keeper.SetParams(ctx, wasmtypes.DefaultParams()))
	// add wasm handler so we can loop-back (contracts calling contracts)
	contractKeeper := NewPermissionedKeeper(*wasmkeeper.NewDefaultPermissionKeeper(&keeper), keeper)

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		msgRouter,
		govtypes.DefaultConfig(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	require.NoError(t, govKeeper.Params.Set(ctx, govv1.DefaultParams()))

	am := module.NewManager( // minimal module set that we use for message/ query tests
		bank.NewAppModule(appCodec, bankKeeper, accountKeeper, subspace(banktypes.ModuleName)),
		staking.NewAppModule(appCodec, stakingKeeper, accountKeeper, bankKeeper, subspace(stakingtypes.ModuleName)),
		distribution.NewAppModule(appCodec, distKeeper, accountKeeper, bankKeeper, stakingKeeper, subspace(distributiontypes.ModuleName)),
		gov.NewAppModule(appCodec, govKeeper, accountKeeper, bankKeeper, subspace(govtypes.ModuleName)),
	)
	am.RegisterServices(module.NewConfigurator(appCodec, msgRouter, querier))
	// wasmplus service
	types.RegisterMsgServer(msgRouter, NewMsgServerImpl(wasmkeeper.NewDefaultPermissionKeeper(keeper)))
	types.RegisterQueryServer(querier, Querier(&keeper))
	// wasm service
	wasmtypes.RegisterMsgServer(msgRouter, wasmkeeper.NewMsgServerImpl(&keeper.Keeper))
	wasmtypes.RegisterQueryServer(querier, WasmQuerier(&keeper))

	keepers := TestKeepers{
		AccountKeeper:  accountKeeper,
		StakingKeeper:  stakingKeeper,
		DistKeeper:     distKeeper,
		ContractKeeper: contractKeeper,
		WasmKeeper:     &keeper,
		BankKeeper:     bankKeeper,
		GovKeeper:      govKeeper,
		IBCKeeper:      ibcKeeper,
		Router:         msgRouter,
		EncodingConfig: encodingConfig,
		Faucet:         faucet,
		MultiStore:     ms,
	}
	return ctx, keepers
}

// TestHandler returns a wasm handler for tests (to avoid circular imports)
func TestHandler(k types.ContractOpsKeeper) wasmkeeper.MessageRouter {
	return wasmkeeper.MessageRouterFunc(func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *wasmtypes.MsgStoreCode:
			return handleStoreCode(ctx, k, msg)
		case *wasmtypes.MsgInstantiateContract:
			return handleInstantiate(ctx, k, msg)
		case *wasmtypes.MsgExecuteContract:
			return handleExecute(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized wasm message type: %T", msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	})
}

func handleStoreCode(ctx sdk.Context, k types.ContractOpsKeeper, msg *wasmtypes.MsgStoreCode) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "sender")
	}
	codeID, _, err := k.Create(ctx, senderAddr, msg.WASMByteCode, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", codeID)),
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleInstantiate(ctx sdk.Context, k types.ContractOpsKeeper, msg *wasmtypes.MsgInstantiateContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, errorsmod.Wrap(err, "admin")
		}
	}

	contractAddr, _, err := k.Instantiate(ctx, msg.CodeID, senderAddr, adminAddr, msg.Msg, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   contractAddr,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleExecute(ctx sdk.Context, k types.ContractOpsKeeper, msg *wasmtypes.MsgExecuteContract) (*sdk.Result, error) {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, errorsmod.Wrap(err, "admin")
	}
	data, err := k.Execute(ctx, contractAddr, senderAddr, msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   data,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func RandomAccountAddress(_ testing.TB) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

// DeterministicAccountAddress  creates a test address with v repeated to valid address size
func DeterministicAccountAddress(_ testing.TB, v byte) sdk.AccAddress {
	return bytes.Repeat([]byte{v}, address.Len)
}

func RandomBech32AccountAddress(t testing.TB) string {
	return RandomAccountAddress(t).String()
}

type ExampleContract struct {
	InitialAmount sdk.Coins
	Creator       crypto.PrivKey
	CreatorAddr   sdk.AccAddress
	CodeID        uint64
	Checksum      []byte
}

func StoreHackatomExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "../../wasm/keeper/testdata/hackatom.wasm")
}

func StoreBurnerExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "../../wasm/keeper/testdata/burner.wasm")
}

func StoreIBCReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "../../wasm/keeper/testdata/ibc_reflect.wasm")
}

func StoreReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleContract {
	return StoreExampleContract(t, ctx, keepers, "../../wasm/keeper/testdata/reflect.wasm")
}

func StoreExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers, wasmFile string) ExampleContract {
	anyAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 1000))
	creator, _, creatorAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, creatorAddr, anyAmount)

	wasmCode, err := os.ReadFile(wasmFile)
	require.NoError(t, err)

	codeID, _, err := keepers.ContractKeeper.Create(ctx, creatorAddr, wasmCode, nil)
	require.NoError(t, err)
	hash := keepers.WasmKeeper.GetCodeInfo(ctx, codeID).CodeHash
	return ExampleContract{anyAmount, creator, creatorAddr, codeID, hash}
}

type ExampleContractInstance struct {
	ExampleContract
	Contract sdk.AccAddress
}

type HackatomExampleInstance struct {
	ExampleContract
	Contract        sdk.AccAddress
	Verifier        crypto.PrivKey
	VerifierAddr    sdk.AccAddress
	Beneficiary     crypto.PrivKey
	BeneficiaryAddr sdk.AccAddress
	Label           string
	Deposit         sdk.Coins
}

// InstantiateHackatomExampleContract load and instantiate the "./testdata/hackatom.wasm" contract
func InstantiateHackatomExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) HackatomExampleInstance {
	contract := StoreHackatomExampleContract(t, ctx, keepers)

	verifier, _, verifierAddr := keyPubAddr()
	fundAccounts(t, ctx, keepers.AccountKeeper, keepers.BankKeeper, verifierAddr, contract.InitialAmount)

	beneficiary, _, beneficiaryAddr := keyPubAddr()
	initMsgBz := HackatomExampleInitMsg{
		Verifier:    verifierAddr,
		Beneficiary: beneficiaryAddr,
	}.GetBytes(t)
	initialAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 100))

	adminAddr := contract.CreatorAddr
	label := "demo contract to query"
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, contract.CodeID, contract.CreatorAddr, adminAddr, initMsgBz, label, initialAmount)
	require.NoError(t, err)
	return HackatomExampleInstance{
		ExampleContract: contract,
		Contract:        contractAddr,
		Verifier:        verifier,
		VerifierAddr:    verifierAddr,
		Beneficiary:     beneficiary,
		BeneficiaryAddr: beneficiaryAddr,
		Label:           label,
		Deposit:         initialAmount,
	}
}

type ExampleInstance struct {
	ExampleContract
	Contract sdk.AccAddress
	Label    string
	Deposit  sdk.Coins
}

// InstantiateReflectExampleContract load and instantiate the "./testdata/reflect.wasm" contract
func InstantiateReflectExampleContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) ExampleInstance {
	example := StoreReflectContract(t, ctx, keepers)
	initialAmount := sdk.NewCoins(sdk.NewInt64Coin("denom", 100))
	label := "demo contract to query"
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, example.CodeID, example.CreatorAddr, nil, []byte("{}"), label, initialAmount)

	require.NoError(t, err)
	return ExampleInstance{
		ExampleContract: example,
		Contract:        contractAddr,
		Label:           label,
		Deposit:         initialAmount,
	}
}

type HackatomExampleInitMsg struct {
	Verifier    sdk.AccAddress `json:"verifier"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
}

func (m HackatomExampleInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type IBCReflectExampleInstance struct {
	Contract      sdk.AccAddress
	Admin         sdk.AccAddress
	CodeID        uint64
	ReflectCodeID uint64
}

// InstantiateIBCReflectContract load and instantiate the "./testdata/ibc_reflect.wasm" contract
func InstantiateIBCReflectContract(t testing.TB, ctx sdk.Context, keepers TestKeepers) IBCReflectExampleInstance {
	reflectID := StoreReflectContract(t, ctx, keepers).CodeID
	ibcReflectID := StoreIBCReflectContract(t, ctx, keepers).CodeID

	initMsgBz := IBCReflectInitMsg{
		ReflectCodeID: reflectID,
	}.GetBytes(t)
	adminAddr := RandomAccountAddress(t)

	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, ibcReflectID, adminAddr, adminAddr, initMsgBz, "ibc-reflect-factory", nil)
	require.NoError(t, err)
	return IBCReflectExampleInstance{
		Admin:         adminAddr,
		Contract:      contractAddr,
		CodeID:        ibcReflectID,
		ReflectCodeID: reflectID,
	}
}

type IBCReflectInitMsg struct {
	ReflectCodeID uint64 `json:"reflect_code_id"`
}

func (m IBCReflectInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

type BurnerExampleInitMsg struct {
	Payout sdk.AccAddress `json:"payout"`
}

func (m BurnerExampleInitMsg) GetBytes(t testing.TB) []byte {
	initMsgBz, err := json.Marshal(m)
	require.NoError(t, err)
	return initMsgBz
}

func fundAccounts(t testing.TB, ctx sdk.Context, am authkeeper.AccountKeeper, bank bankkeeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	acc := am.NewAccountWithAddress(ctx, addr)
	am.SetAccount(ctx, acc)
	wasmkeeper.NewTestFaucet(t, ctx, bank, minttypes.ModuleName, coins...).Fund(ctx, addr, coins...)
}

var keyCounter uint64

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
//
//nolint:unparam
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func GenerateSortedBech32Address(addrs ...[]byte) []string {
	sort.Slice(addrs, func(i, j int) bool {
		return string(addrs[i]) < string(addrs[j])
	})

	expAddrs := make([]string, len(addrs))
	for i, b := range addrs {
		expAddrs[i] = sdk.AccAddress(b).String()
	}

	return expAddrs
}
