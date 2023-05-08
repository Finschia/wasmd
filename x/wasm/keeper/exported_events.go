package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	wasmvmtypes "github.com/line/wasmvm/types"
)

func NewCustomCallablePointEvents(evts wasmvmtypes.Events, contractAddr sdk.AccAddress, callstack []byte) (sdk.Events, error) {
	return newCustomCallablePointEvents(evts, contractAddr, callstack)
}
