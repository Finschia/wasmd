package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	wasmTypes "github.com/Finschia/wasmd/x/wasm/types"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint:staticcheck
	legacy.RegisterAminoMsg(cdc, &MsgStoreCodeAndInstantiateContract{}, "wasm/MsgStoreCodeAndInstantiateContract")

	cdc.RegisterConcrete(&DeactivateContractProposal{}, "wasm/DeactivateContractProposal", nil)
	cdc.RegisterConcrete(&ActivateContractProposal{}, "wasm/ActivateContractProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	wasmTypes.RegisterInterfaces(registry)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCodeAndInstantiateContract{},
	)
	registry.RegisterImplementations(
		(*v1beta1.Content)(nil),
		&DeactivateContractProposal{},
		&ActivateContractProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
