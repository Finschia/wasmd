package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	wasmvmtypes "github.com/Finschia/wasmvm/types"
)

func NewCustomCallablePointEvents(evts wasmvmtypes.Events, contractAddr sdk.AccAddress, callstack []byte) (sdk.Events, error) {
	return newCustomCallablePointEvents(evts, contractAddr, callstack)
}
