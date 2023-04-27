package keeper

import (
	"github.com/line/lbm-sdk/store/prefix"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/wasmd/x/wasm/types"
	wasmvm "github.com/line/wasmvm"
)

type PlusKeeper interface {
	CosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI
}

func (k Keeper) ContractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.ContractInfo, types.CodeInfo, prefix.Store, error) {
	return k.contractInstance(ctx, contractAddress)
}

func (k Keeper) GasMeter(ctx sdk.Context) MultipliedGasMeter {
	return k.gasMeter(ctx)
}

func (k Keeper) GetGasRegister() GasRegister {
	return k.gasRegister
}

func (k Keeper) GetWasmVM() types.WasmerEngine {
	return k.wasmVM
}

func (k Keeper) GetWasmVMQueryHandler() WasmVMQueryHandler {
	return k.wasmVMQueryHandler
}

func (k Keeper) GetCosmwasmAPIGenerator() CosmwasmAPIGenerator {
	return k.cosmwasmAPIGenerator
}

func (k *Keeper) SetCosmwasmAPIGenerator(generator CosmwasmAPIGenerator) {
	k.cosmwasmAPIGenerator = generator
}
