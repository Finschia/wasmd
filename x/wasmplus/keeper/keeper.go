package keeper

import (
	"context"
	"fmt"

	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/log"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type Keeper struct {
	wasmkeeper.Keeper
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
	metrics      *wasmkeeper.Metrics
	bank         wasmtypes.BankKeeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper wasmtypes.AccountKeeper,
	bankKeeper wasmtypes.BankKeeper,
	stakingKeeper wasmtypes.StakingKeeper,
	distKeeper wasmtypes.DistributionKeeper,
	ics4Wrapper wasmtypes.ICS4Wrapper,
	channelKeeper wasmtypes.ChannelKeeper,
	portKeeper wasmtypes.PortKeeper,
	capabilityKeeper wasmtypes.CapabilityKeeper,
	portSource wasmtypes.ICS20TransferPortSource,
	router wasmkeeper.MessageRouter,
	queryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	wasmConfig wasmtypes.WasmConfig,
	availableCapabilities string,
	authority string,
	opts ...wasmkeeper.Option,
) Keeper {
	bankKeeper, ok := bankKeeper.(bankkeeper.Keeper)
	if !ok {
		panic("bankKeeper should be bankPlusKeeper")
	}
	result := Keeper{
		cdc:          cdc,
		storeService: storeService,
		metrics:      wasmkeeper.NopMetrics(),
		bank:         bankKeeper,
	}
	result.Keeper = wasmkeeper.NewKeeper(
		cdc,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distKeeper,
		ics4Wrapper,
		channelKeeper,
		portKeeper,
		capabilityKeeper,
		portSource,
		router,
		queryRouter,
		homeDir,
		wasmConfig,
		availableCapabilities,
		authority,
		opts...,
	)
	return result
}

func WasmQuerier(k *Keeper) wasmtypes.QueryServer {
	return wasmkeeper.NewGrpcQuerier(k.cdc, k.storeService, k, k.QueryGasLimit())
}

func Querier(k *Keeper) types.QueryServer {
	return newGrpcQuerier(k.storeService, k)
}

func (Keeper) Logger(ctx sdk.Context) log.Logger {
	return ModuleLogger(ctx)
}

func ModuleLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IsInactiveContract(ctx context.Context, contractAddress sdk.AccAddress) bool {
	store := k.storeService.OpenKVStore(ctx)
	ok, err := store.Has(types.GetInactiveContractKey(contractAddress))
	if err != nil {
		panic(err)
	}
	return ok
}

func (k Keeper) IterateInactiveContracts(ctx context.Context, fn func(contractAddress sdk.AccAddress) (stop bool)) {
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.InactiveContractPrefix)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractAddress := sdk.AccAddress(iterator.Value())
		if stop := fn(contractAddress); stop {
			break
		}
	}
}

func (k Keeper) addInactiveContract(ctx context.Context, contractAddress sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetInactiveContractKey(contractAddress)

	err := store.Set(key, contractAddress)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) deleteInactiveContract(ctx context.Context, contractAddress sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetInactiveContractKey(contractAddress)
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
}

// activateContract delete the contract address from inactivateContract list if the contract is deactivated.
func (k Keeper) activateContract(ctx context.Context, contractAddress sdk.AccAddress) error {
	if !k.IsInactiveContract(ctx, contractAddress) {
		return errorsmod.Wrapf(wasmtypes.ErrNotFound, "no inactivate contract %s", contractAddress.String())
	}

	k.deleteInactiveContract(ctx, contractAddress)
	// todo: add bankplus function
	// k.bank.DeleteFromInactiveAddr(ctx, contractAddress)

	return nil
}

// deactivateContract add the contract address to inactivateContract list.
func (k Keeper) deactivateContract(ctx context.Context, contractAddress sdk.AccAddress) error {
	if k.IsInactiveContract(ctx, contractAddress) {
		return errorsmod.Wrapf(wasmtypes.ErrAccountExists, "already inactivate contract %s", contractAddress.String())
	}
	if !k.HasContractInfo(ctx, contractAddress) {
		return errorsmod.Wrapf(wasmtypes.ErrInvalid, "no contract %s", contractAddress.String())
	}

	k.addInactiveContract(ctx, contractAddress)
	// todo: add bankplus function
	// k.bank.AddToInactiveAddr(ctx, contractAddress)

	return nil
}
