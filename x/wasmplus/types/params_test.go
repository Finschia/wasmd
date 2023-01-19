package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"

	wasmkeeper "github.com/line/wasmd/x/wasm/keeper"
	wasmtypes "github.com/line/wasmd/x/wasm/types"
)

func TestValidateParams(t *testing.T) {
	var (
		anyAddress     sdk.AccAddress = make([]byte, wasmtypes.ContractAddrLen)
		otherAddress   sdk.AccAddress = bytes.Repeat([]byte{1}, wasmtypes.ContractAddrLen)
		invalidAddress                = "invalid address"
	)

	specs := map[string]struct {
		src    Params
		expErr bool
	}{
		"all good with defaults": {
			src: DefaultParams(),
		},
		"all good with nobody": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeNobody,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
		},
		"all good with everybody": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowEverybody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
		},
		"all good with only address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessTypeOnlyAddress.With(anyAddress),
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
		},
		"all good with anyOf address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessTypeAnyOfAddresses.With(anyAddress),
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
		},
		"all good with anyOf addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessTypeAnyOfAddresses.With(anyAddress, otherAddress),
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
		},
		"reject empty type in instantiate permission": {
			src: Params{
				CodeUploadAccess: wasmtypes.AllowNobody,
				GasMultiplier:    wasmkeeper.DefaultGasMultiplier,
				InstanceCost:     wasmkeeper.DefaultInstanceCost,
				CompileCost:      wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject unknown type in instantiate": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowNobody,
				InstantiateDefaultPermission: 1111,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess invalid address in only address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeOnlyAddress, Address: invalidAddress},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject wrong field addresses in only address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeOnlyAddress, Address: anyAddress.String(), Addresses: []string{anyAddress.String()}},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Everybody with obsolete address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeEverybody, Address: anyAddress.String()},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject CodeUploadAccess Nobody with obsolete address": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeNobody, Address: anyAddress.String()},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty CodeUploadAccess": {
			src: Params{
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject undefined permission in CodeUploadAccess": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeUnspecified},
				InstantiateDefaultPermission: wasmtypes.AccessTypeOnlyAddress,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject empty addresses in any of addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeAnyOfAddresses, Addresses: []string{}},
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject addresses not set in any of addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeAnyOfAddresses},
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject invalid address in any of addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeAnyOfAddresses, Addresses: []string{invalidAddress}},
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject duplicate address in any of addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeAnyOfAddresses, Addresses: []string{anyAddress.String(), anyAddress.String()}},
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject wrong field address in any of  addresses": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeAnyOfAddresses, Address: anyAddress.String(), Addresses: []string{anyAddress.String()}},
				InstantiateDefaultPermission: wasmtypes.AccessTypeAnyOfAddresses,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject zero GasMultiplier": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowEverybody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
				GasMultiplier:                0,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject zero InstanceCost": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowEverybody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 0,
				CompileCost:                  wasmkeeper.DefaultCompileCost,
			},
			expErr: true,
		},
		"reject zero CompileCost": {
			src: Params{
				CodeUploadAccess:             wasmtypes.AllowEverybody,
				InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
				GasMultiplier:                wasmkeeper.DefaultGasMultiplier,
				InstanceCost:                 wasmkeeper.DefaultInstanceCost,
				CompileCost:                  0,
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
