package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
)

type ViewKeeper interface {
	IterateInactiveContracts(ctx context.Context, fn func(contractAddress sdk.AccAddress) bool)
	IsInactiveContract(ctx context.Context, contractAddress sdk.AccAddress) bool
}

type ContractOpsKeeper interface {
	wasmtypes.ContractOpsKeeper

	// DeactivateContract add the contract address to inactive contract list.
	DeactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error

	// ActivateContract remove the contract address from inactive contract list.
	ActivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error
}
