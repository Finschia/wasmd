package lbmtypes

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	paramtypes "github.com/line/lbm-sdk/x/params/types"

	wasmtypes "github.com/line/wasmd/x/wasm/types"
)

var ParamStoreKeyGasMultiplier = []byte("gasMultiplier")
var ParamStoreKeyInstanceCost = []byte("instanceCost")
var ParamStoreKeyCompileCost = []byte("compileCost")

func DefaultParams() Params {
	return Params{
		CodeUploadAccess:             wasmtypes.AllowEverybody,
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
		// todo should set wasmkepper's default after solving import cycle problem.
		GasMultiplier: 0, // keeper.DefaultGasMultiplier,
		InstanceCost:  0, // keeper.DefaultInstanceCost,
		CompileCost:   0, // keeper.DefaultCompileCost,
	}
}

func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyUploadAccess, &p.CodeUploadAccess, validateAccessConfig),
		paramtypes.NewParamSetPair(wasmtypes.ParamStoreKeyInstantiateAccess, &p.InstantiateDefaultPermission, validateAccessType),
		paramtypes.NewParamSetPair(ParamStoreKeyGasMultiplier, &p.GasMultiplier, validateGasMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyInstanceCost, &p.InstanceCost, validateInstanceCost),
		paramtypes.NewParamSetPair(ParamStoreKeyCompileCost, &p.CompileCost, validateCompileCost),
	}
}

// ValidateBasic performs basic validation on wasm parameters
func (p Params) ValidateBasic() error {
	if err := validateAccessType(p.InstantiateDefaultPermission); err != nil {
		return errors.Wrap(err, "instantiate default permission")
	}
	if err := validateAccessConfig(p.CodeUploadAccess); err != nil {
		return errors.Wrap(err, "upload access")
	}
	return nil
}

func validateAccessType(i interface{}) error {
	a, ok := i.(wasmtypes.AccessType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if a == wasmtypes.AccessTypeUnspecified {
		return sdkerrors.Wrap(wasmtypes.ErrEmpty, "type")
	}
	for _, v := range wasmtypes.AllAccessTypes {
		if v == a {
			return nil
		}
	}
	return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "unknown type: %q", a)
}

func validateAccessConfig(i interface{}) error {
	v, ok := i.(wasmtypes.AccessConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.ValidateBasic()
}

func validateGasMultiplier(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater than 0")
	}
	return nil
}

func validateInstanceCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater than 0")
	}
	return nil
}

func validateCompileCost(i interface{}) error {
	a, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "type: %T", i)
	}
	if a == 0 {
		return sdkerrors.Wrap(wasmtypes.ErrInvalid, "must be greater than 0")
	}
	return nil
}
