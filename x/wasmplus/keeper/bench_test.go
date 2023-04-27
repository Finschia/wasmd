package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/wasmd/x/wasm/types"
)

func BenchmarkAPI(b *testing.B) {
	wasmConfig := types.WasmConfig{MemoryCacheSize: 0}
	ctx, keepers := createTestInput(b, false, AvailableCapabilities, wasmConfig, dbm.NewMemDB())
	example := InstantiateHackatomExampleContract(b, ctx, keepers)
	api := keepers.WasmKeeper.GetCosmwasmAPIGenerator().Generate(&ctx)
	addrStr := example.Contract.String()
	addrBytes, err := sdk.AccAddressFromBech32(example.Contract.String())
	require.NoError(b, err)

	b.Run("CanonicalAddress", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := api.CanonicalAddress(addrStr)
			require.NoError(b, err)
		}
	})

	b.Run("HumanAddress", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, err := api.HumanAddress(addrBytes)
			require.NoError(b, err)
		}
	})
}
