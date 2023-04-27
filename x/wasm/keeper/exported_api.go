package keeper

import sdk "github.com/line/lbm-sdk/types"

var (
	HumanAddress            = humanAddress
	CanonicalAddress        = canonicalAddress
	CostHumanize            = DefaultGasCostHumanAddress * DefaultGasMultiplier
	CostCanonical           = DefaultGasCostCanonicalAddress * DefaultGasMultiplier
	CostJSONDeserialization = costJSONDeserialization
)

type CosmwasmAPIImpl = cosmwasmAPIImpl

func NewCosmwasmAPIImpl(k *Keeper, ctx *sdk.Context) cosmwasmAPIImpl {
	return cosmwasmAPIImpl{keeper: k, ctx: ctx}
}

func (a cosmwasmAPIImpl) CallCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	return a.callCallablePoint(contractAddrStr, name, args, isReadonly, callstack, gasLimit)
}

func (a cosmwasmAPIImpl) ValidateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	return a.validateInterface(contractAddrStr, expectedInterface)
}
