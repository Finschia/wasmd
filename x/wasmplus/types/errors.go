package types

import (
	errorsmod "cosmossdk.io/errors"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
)

// Codes for wasm contract errors
var (
	// ErrInactiveContract error if the contract set inactive
	ErrInactiveContract = errorsmod.Register(wasmtypes.DefaultCodespace, 101, "inactive contract")
)
