package keeper

import (
	"context"

	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	wasmtypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

type queryKeeper interface {
	IterateInactiveContracts(ctx context.Context, fn func(contractAddress sdk.AccAddress) bool)
	IsInactiveContract(ctx context.Context, contractAddress sdk.AccAddress) bool
	HasContractInfo(ctx context.Context, contractAddress sdk.AccAddress) bool
}

var _ types.QueryServer = &grpcQuerier{}

type grpcQuerier struct {
	keeper       queryKeeper
	storeService corestoretypes.KVStoreService
}

// newGrpcQuerier constructor
func newGrpcQuerier(storeService corestoretypes.KVStoreService, keeper queryKeeper) *grpcQuerier {
	return &grpcQuerier{storeService: storeService, keeper: keeper}
}

func (q grpcQuerier) InactiveContracts(c context.Context, req *types.QueryInactiveContractsRequest) (*types.QueryInactiveContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	addresses := make([]string, 0)
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(q.storeService.OpenKVStore(ctx)), types.InactiveContractPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			contractAddress := sdk.AccAddress(value)
			addresses = append(addresses, contractAddress.String())
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryInactiveContractsResponse{
		Addresses:  addresses,
		Pagination: pageRes,
	}, nil
}

func (q grpcQuerier) InactiveContract(c context.Context, req *types.QueryInactiveContractRequest) (*types.QueryInactiveContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	contractAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	if !q.keeper.HasContractInfo(ctx, contractAddr) {
		return nil, wasmtypes.ErrNotFound
	}

	inactivated := q.keeper.IsInactiveContract(ctx, contractAddr)
	return &types.QueryInactiveContractResponse{
		Inactivated: inactivated,
	}, nil
}
