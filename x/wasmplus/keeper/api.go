package keeper

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	wasmvm "github.com/Finschia/wasmvm"
)

type cosmwasmAPIImpl struct {
	keeper  *Keeper
	ctx     *sdk.Context
	wasmAPI wasmkeeper.CosmwasmAPIImpl
}

type cosmwasmAPIGeneratorImpl struct {
	keeper *Keeper
}

func (a cosmwasmAPIImpl) callCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	contractAddr, err := sdk.AccAddressFromBech32(contractAddrStr)

	if err != nil {
		return nil, 0, fmt.Errorf("specified callee address is invalid: %s", err)
	}

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("called contract cannot be executed")
	}

	return a.wasmAPI.CallCallablePoint(contractAddrStr, name, args, isReadonly, callstack, gasLimit)
}

func (a cosmwasmAPIImpl) validateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	contractAddr, err := sdk.AccAddressFromBech32(contractAddrStr)

	if err != nil {
		return nil, 0, fmt.Errorf("specified contract address is invalid: %s", err)
	}

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("try to validate a contract cannot be executed")
	}

	return a.wasmAPI.ValidateInterface(contractAddrStr, expectedInterface)
}

func (g cosmwasmAPIGeneratorImpl) Generate(ctx *sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		keeper:  g.keeper,
		ctx:     ctx,
		wasmAPI: wasmkeeper.NewCosmwasmAPIImpl(&g.keeper.Keeper, ctx),
	}
	return wasmvm.GoAPI{
		HumanAddress:      wasmkeeper.HumanAddress,
		CanonicalAddress:  wasmkeeper.CanonicalAddress,
		CallCallablePoint: x.callCallablePoint,
		ValidateInterface: x.validateInterface,
	}
}
