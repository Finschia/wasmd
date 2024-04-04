package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
)

const (
	ProposalTypeDeactivateContract wasmtypes.ProposalType = "DeactivateContract"
	ProposalTypeActivateContract   wasmtypes.ProposalType = "ActivateContract"
)

var EnableAllProposals = append([]wasmtypes.ProposalType{
	ProposalTypeDeactivateContract,
	ProposalTypeActivateContract,
}, wasmtypes.EnableAllProposals...)

func init() {
	v1beta1.RegisterProposalType(string(ProposalTypeDeactivateContract))
	v1beta1.RegisterProposalType(string(ProposalTypeActivateContract))
}

func (p DeactivateContractProposal) GetTitle() string { return p.Title }

func (p DeactivateContractProposal) GetDescription() string { return p.Description }

func (p DeactivateContractProposal) ProposalRoute() string { return wasmtypes.RouterKey }

func (p DeactivateContractProposal) ProposalType() string {
	return string(ProposalTypeDeactivateContract)
}

func (p DeactivateContractProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "contract")
	}

	return nil
}

func (p DeactivateContractProposal) String() string {
	return fmt.Sprintf(`Deactivate Contract Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
`, p.Title, p.Description, p.Contract)
}

func (p ActivateContractProposal) GetTitle() string { return p.Title }

func (p ActivateContractProposal) GetDescription() string { return p.Description }

func (p ActivateContractProposal) ProposalRoute() string { return wasmtypes.RouterKey }

func (p ActivateContractProposal) ProposalType() string {
	return string(ProposalTypeActivateContract)
}

func (p ActivateContractProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "contract")
	}

	return nil
}

func (p ActivateContractProposal) String() string {
	return fmt.Sprintf(`Activate Contract Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
`, p.Title, p.Description, p.Contract)
}
