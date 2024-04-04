package types

import (
	errorsmod "cosmossdk.io/errors"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
)

func validateWasmCode(s []byte) error {
	if len(s) == 0 {
		return errorsmod.Wrap(wasmtypes.ErrEmpty, "is required")
	}
	if len(s) > wasmtypes.MaxWasmSize {
		return errorsmod.Wrapf(wasmtypes.ErrLimit, "cannot be longer than %d bytes", wasmtypes.MaxWasmSize)
	}
	return nil
}

func validateLabel(label string) error {
	if label == "" {
		return errorsmod.Wrap(wasmtypes.ErrEmpty, "is required")
	}
	if len(label) > wasmtypes.MaxLabelSize {
		return errorsmod.Wrapf(wasmtypes.ErrLimit, "cannot be longer than %d characters", wasmtypes.MaxWasmSize)
	}
	return nil
}
