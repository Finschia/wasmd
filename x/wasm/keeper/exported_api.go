package keeper

import sdk "github.com/line/lbm-sdk/types"

var (
	HumanAddress     = humanAddress
	CanonicalAddress = canonicalAddress
	CostHumanize     = DefaultGasCostHumanAddress * DefaultGasMultiplier
	CostCanonical    = DefaultGasCostCanonicalAddress * DefaultGasMultiplier
)

func NewCosmwasmAPIImpl(k *Keeper, ctx *sdk.Context) CosmwasmAPIImpl {
	return CosmwasmAPIImpl{keeper: k, ctx: ctx}
}

func (a CosmwasmAPIImpl) CallCallablePoint(contractAddrStr string, name []byte, args []byte, isReadonly bool, callstack []byte, gasLimit uint64) ([]byte, uint64, error) {
	return a.callCallablePoint(contractAddrStr, name, args, isReadonly, callstack, gasLimit)
}

func (a CosmwasmAPIImpl) ValidateInterface(contractAddrStr string, expectedInterface []byte) ([]byte, uint64, error) {
	return a.validateInterface(contractAddrStr, expectedInterface)
}
