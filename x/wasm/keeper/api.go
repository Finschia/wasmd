package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	wasmvm "github.com/line/wasmvm"
	wasmvmtypes "github.com/line/wasmvm/types"

	types "github.com/line/wasmd/x/wasm/types"
)

type cosmwasmAPIImpl struct {
	keeper *Keeper
	ctx    *sdk.Context
}

const (
	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1
)

var (
	costJSONDeserialization = wasmvmtypes.UFraction{
		Numerator:   DefaultDeserializationCostPerByte * types.DefaultGasMultiplier,
		Denominator: 1,
	}
)

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	gasMultiplier := a.keeper.getGasMultiplier(*a.ctx)
	gas := gasMultiplier.ToWasmVMGas(5)
	if err := sdk.VerifyAddressFormat(canon); err != nil {
		return "", gas, err
	}

	return sdk.AccAddress(canon).String(), gas, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	gasMultiplier := a.keeper.getGasMultiplier(*a.ctx)
	return bz, gasMultiplier.ToWasmVMGas(4), err
}

// returns result, gas used, error
func (a cosmwasmAPIImpl) callCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)
	contractInfo, codeInfo, prefixStore, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("called contract cannot be executed")
	}

	env := types.NewEnv(*a.ctx, contractAddr)
	wasmStore := types.NewWasmStore(prefixStore)
	gasMultiplier := a.keeper.getGasMultiplier(*a.ctx)
	querier := NewQueryHandler(*a.ctx, a.keeper.wasmVMQueryHandler, contractAddr, gasMultiplier)
	gasMeter := a.keeper.gasMeter(*a.ctx)
	api := a.keeper.cosmwasmAPI(*a.ctx)

	instantiateCost := gasMultiplier.ToWasmVMGas(a.keeper.instantiateContractCosts(a.keeper.gasRegister, *a.ctx, a.keeper.IsPinnedCode(*a.ctx, contractInfo.CodeID), len(args)))
	if gasLimit < instantiateCost {
		return nil, 0, fmt.Errorf("lack of gas for calling callable point")
	}
	wasmGasLimit := gasLimit - instantiateCost

	result, events, attrs, gas, err := a.keeper.wasmVM.CallCallablePoint(name, codeInfo.CodeHash, isReadonly, callstack, env, args, wasmStore, api, querier, gasMeter, wasmGasLimit, costJSONDeserialization)
	gas += instantiateCost

	if err != nil {
		return nil, gas, err
	}

	if !isReadonly {
		// issue events and attrs
		if len(attrs) != 0 {
			eventsByAttr, err := newCallablePointEvent(attrs, contractAddr, callstack)
			if err != nil {
				return nil, gas, err
			}
			a.ctx.EventManager().EmitEvents(eventsByAttr)
		}

		if len(events) != 0 {
			customEvents, err := newCustomCallablePointEvents(events, contractAddr, callstack)
			if err != nil {
				return nil, gas, err
			}
			a.ctx.EventManager().EmitEvents(customEvents)
		}
	}

	return result, gas, err
}

// returns result, gas used, error
func (a cosmwasmAPIImpl) validateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)

	if a.keeper.IsInactiveContract(*a.ctx, contractAddr) {
		return nil, 0, fmt.Errorf("try to validate a contract cannot be executed")
	}

	_, codeInfo, _, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	result, err := a.keeper.wasmVM.ValidateDynamicLinkInterface(codeInfo.CodeHash, expectedInterface)

	return result, 0, err
}

func (k Keeper) cosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		keeper: &k,
		ctx:    &ctx,
	}
	return wasmvm.GoAPI{
		HumanAddress:      x.humanAddress,
		CanonicalAddress:  x.canonicalAddress,
		CallCallablePoint: x.callCallablePoint,
		ValidateInterface: x.validateInterface,
	}
}
