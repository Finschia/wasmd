package keeper

import (
	"fmt"
	
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	bankpluskeeper "github.com/line/lbm-sdk/x/bankplus/keeper"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	"github.com/line/ostracon/libs/log"

	wasmkeeper "github.com/line/wasmd/x/wasm/keeper"
	wasmtypes "github.com/line/wasmd/x/wasm/types"
	"github.com/line/wasmd/x/wasmplus/types"
)

type Keeper struct {
	wasmkeeper.Keeper
	cdc        codec.Codec
	storeKey   sdk.StoreKey
	metrics    *Metrics
	bank       bankpluskeeper.Keeper
	paramSpace paramtypes.Subspace
}

func NewKeeper(
	cdc codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper wasmtypes.AccountKeeper,
	bankKeeper wasmtypes.BankKeeper,
	stakingKeeper wasmtypes.StakingKeeper,
	distKeeper wasmtypes.DistributionKeeper,
	channelKeeper wasmtypes.ChannelKeeper,
	portKeeper wasmtypes.PortKeeper,
	capabilityKeeper wasmtypes.CapabilityKeeper,
	portSource wasmtypes.ICS20TransferPortSource,
	router wasmkeeper.MessageRouter,
	queryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	wasmConfig wasmtypes.WasmConfig,
	availableCapabilities string,
	customEncoders *wasmkeeper.MessageEncoders,
	customPlugins *wasmkeeper.QueryPlugins,
	opts ...wasmkeeper.Option,
) Keeper {
	bankPlusKeeper, ok := bankKeeper.(bankpluskeeper.Keeper)
	if !ok {
		panic("bankKeeper should be bankPlusKeeper")
	}
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	result := Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		metrics:    NopMetrics(),
		bank:       bankPlusKeeper,
		paramSpace: paramSpace,
	}
	result.Keeper = wasmkeeper.NewKeeper(
		cdc,
		storeKey,
		paramSpace,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		channelKeeper,
		portKeeper,
		capabilityKeeper,
		portSource,
		router,
		queryRouter,
		homeDir,
		wasmConfig,
		availableCapabilities,
		//customEncoders,
		//customPlugins,
		opts...,
	)
	return result
}

func WasmQuerier(k *Keeper) wasmtypes.QueryServer {
	return wasmkeeper.NewGrpcQuerier(k.cdc, k.storeKey, k, k.QueryGasLimit())
}

func Querier(k *Keeper) types.QueryServer {
	return NewQuerier(k.storeKey, k)
}

func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ModuleLogger(ctx)
}

func ModuleLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCustomParams returns the total set of wasm parameters.
func (k Keeper) getParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) setParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}
