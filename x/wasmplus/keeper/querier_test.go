package keeper

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

func TestQueryInactiveContracts(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	keeper := keepers.WasmKeeper

	example1 := InstantiateHackatomExampleContract(t, ctx, keepers)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	example2 := InstantiateHackatomExampleContract(t, ctx, keepers)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// set inactive
	err := keeper.deactivateContract(ctx, example1.Contract)
	require.NoError(t, err)
	err = keeper.deactivateContract(ctx, example2.Contract)
	require.NoError(t, err)

	q := Querier(keeper)
	rq := types.QueryInactiveContractsRequest{}
	res, err := q.InactiveContracts(sdk.WrapSDKContext(ctx), &rq)
	require.NoError(t, err)
	expect := []string{example1.Contract.String(), example2.Contract.String()}
	for _, exp := range expect {
		assert.Contains(t, res.Addresses, exp)
	}
}

func TestQueryInactiveContract(t *testing.T) {
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	keeper := keepers.WasmKeeper

	example := InstantiateHackatomExampleContract(t, ctx, keepers)
	contractAddr := example.Contract
	q := Querier(keeper)
	rq := &types.QueryInactiveContractRequest{Address: example.Contract.String()}

	// confirm that Contract is active
	got, err := q.InactiveContract(sdk.WrapSDKContext(ctx), rq)
	require.NoError(t, err)
	require.False(t, got.Inactivated)

	// set inactive
	err = keeper.deactivateContract(ctx, example.Contract)
	require.NoError(t, err)

	specs := map[string]struct {
		srcQuery       *types.QueryInactiveContractRequest
		expInactivated bool
		expErr         error
	}{
		"query": {
			srcQuery:       &types.QueryInactiveContractRequest{Address: contractAddr.String()},
			expInactivated: true,
		},
		"query with unknown address": {
			srcQuery: &types.QueryInactiveContractRequest{Address: RandomBech32AccountAddress(t)},
			expErr:   wasmtypes.ErrNotFound,
		},
		"with empty request": {
			srcQuery: nil,
			expErr:   status.Error(codes.InvalidArgument, "empty request"),
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err = q.InactiveContract(sdk.WrapSDKContext(ctx), spec.srcQuery)

			if spec.expErr != nil {
				require.Equal(t, spec.expErr, err, "but got %+v", err)
				return
			}

			require.NoError(t, err)
			require.True(t, got.Inactivated)
		})
	}
}
