package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/wasmd/x/wasm/types"
	wasmvm "github.com/line/wasmvm"
	wasmvmtypes "github.com/line/wasmvm/types"
)

const (
	// DefaultGasCostHumanAddress is how moch SDK gas we charge to convert to a human address format
	DefaultGasCostHumanAddress = 5
	// DefaultGasCostCanonicalAddress is how moch SDK gas we charge to convert to a canonical address format
	DefaultGasCostCanonicalAddress = 4

	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1
)

var (
	costHumanize            = DefaultGasCostHumanAddress * DefaultGasMultiplier
	costCanonical           = DefaultGasCostCanonicalAddress * DefaultGasMultiplier
	costJSONDeserialization = wasmvmtypes.UFraction{
		Numerator:   DefaultDeserializationCostPerByte * DefaultGasMultiplier,
		Denominator: 1,
	}
)

type CosmwasmAPIImpl struct {
	keeper *Keeper
	ctx    *sdk.Context
}

type cosmwasmAPIGeneratorImpl struct {
	keeper *Keeper
}

type CosmwasmAPIGenerator interface {
	Generate(ctx *sdk.Context) wasmvm.GoAPI
}

func humanAddress(canon []byte) (string, uint64, error) {
	if err := sdk.VerifyAddressFormat(canon); err != nil {
		return "", costHumanize, err
	}
	return sdk.AccAddress(canon).String(), costHumanize, nil
}

func canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	return bz, costCanonical, err
}

// callCallablePoint is a wrapper function of `wasmvm`
// returns result, gas used, error
func (a CosmwasmAPIImpl) callCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)
	contractInfo, codeInfo, prefixStore, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	env := types.NewEnv(*a.ctx, contractAddr)
	wasmStore := types.NewWasmStore(prefixStore)
	gasRegister := a.keeper.GetGasRegister()
	querier := NewQueryHandler(*a.ctx, a.keeper.wasmVMQueryHandler, contractAddr, gasRegister)
	gasMeter := a.keeper.gasMeter(*a.ctx)
	api := a.keeper.cosmwasmAPIGenerator.Generate(a.ctx)

	instantiateCost := gasRegister.ToWasmVMGas(gasRegister.InstantiateContractCosts(a.keeper.IsPinnedCode(*a.ctx, contractInfo.CodeID), len(args)))
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

// validateInterface is a wrapper function of `wasmvm`
// returns result, gas used, error
func (a CosmwasmAPIImpl) validateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)

	_, codeInfo, _, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return nil, 0, err
	}

	result, err := a.keeper.wasmVM.ValidateDynamicLinkInterface(codeInfo.CodeHash, expectedInterface)

	return result, 0, err
}

func (c cosmwasmAPIGeneratorImpl) Generate(ctx *sdk.Context) wasmvm.GoAPI {
	x := CosmwasmAPIImpl{
		keeper: c.keeper,
		ctx:    ctx,
	}
	return wasmvm.GoAPI{
		HumanAddress:      humanAddress,
		CanonicalAddress:  canonicalAddress,
		CallCallablePoint: x.callCallablePoint,
		ValidateInterface: x.validateInterface,
	}
}
