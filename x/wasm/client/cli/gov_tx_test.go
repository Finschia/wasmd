package cli

import (
	"os"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/wasmd/x/wasm/keeper/testdata"
	"github.com/Finschia/wasmd/x/wasm/types"
)

func TestParseAccessConfigUpdates(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("link", "linkpub")
	config.SetBech32PrefixForValidator("linkvaloper", "linkvaloperpub")
	config.SetBech32PrefixForConsensusNode("linkvalcons", "linkvalconspub")
	config.SetCoinType(438)
	config.Seal()

	specs := map[string]struct {
		src    []string
		exp    []types.AccessConfigUpdate
		expErr bool
	}{
		"nobody": {
			src: []string{"1:nobody"},
			exp: []types.AccessConfigUpdate{{
				CodeID:                1,
				InstantiatePermission: types.AccessConfig{Permission: types.AccessTypeNobody},
			}},
		},
		"everybody": {
			src: []string{"1:everybody"},
			exp: []types.AccessConfigUpdate{{
				CodeID:                1,
				InstantiatePermission: types.AccessConfig{Permission: types.AccessTypeEverybody},
			}},
		},
		"any of addresses - single": {
			src: []string{"1:link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395"},
			exp: []types.AccessConfigUpdate{
				{
					CodeID: 1,
					InstantiatePermission: types.AccessConfig{
						Permission: types.AccessTypeAnyOfAddresses,
						Addresses:  []string{"link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395"},
					},
				},
			},
		},
		"any of addresses - multiple": {
			src: []string{"1:link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395,link1hmayw7vv0p3gzeh3jzwmw9xj8fy8a3kmpqgjrysljdnecqkps02qrq5rvm"},
			exp: []types.AccessConfigUpdate{
				{
					CodeID: 1,
					InstantiatePermission: types.AccessConfig{
						Permission: types.AccessTypeAnyOfAddresses,
						Addresses:  []string{"link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395", "link1hmayw7vv0p3gzeh3jzwmw9xj8fy8a3kmpqgjrysljdnecqkps02qrq5rvm"},
					},
				},
			},
		},
		"multiple code ids with different permissions": {
			src: []string{"1:link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395,link1hmayw7vv0p3gzeh3jzwmw9xj8fy8a3kmpqgjrysljdnecqkps02qrq5rvm", "2:nobody"},
			exp: []types.AccessConfigUpdate{
				{
					CodeID: 1,
					InstantiatePermission: types.AccessConfig{
						Permission: types.AccessTypeAnyOfAddresses,
						Addresses:  []string{"link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395", "link1hmayw7vv0p3gzeh3jzwmw9xj8fy8a3kmpqgjrysljdnecqkps02qrq5rvm"},
					},
				}, {
					CodeID: 2,
					InstantiatePermission: types.AccessConfig{
						Permission: types.AccessTypeNobody,
					},
				},
			},
		},
		"any of addresses - empty list": {
			src:    []string{"1:"},
			expErr: true,
		},
		"any of addresses - invalid address": {
			src:    []string{"1:foo"},
			expErr: true,
		},
		"any of addresses - duplicate address": {
			src:    []string{"1:link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395,link12kr02kew9fl73rqekalavuu0xaxcgwr6nkk395"},
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			got, gotErr := parseAccessConfigUpdates(spec.src)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, spec.exp, got)
		})
	}
}

func TestParseCodeInfoFlags(t *testing.T) {
	correctSource := "https://github.com/CosmWasm/wasmd/blob/main/x/wasm/keeper/testdata/hackatom.wasm"
	correctBuilderRef := "cosmwasm/workspace-optimizer:0.12.9"

	wasmBin, err := os.ReadFile("../../keeper/testdata/hackatom.wasm.gzip")
	require.NoError(t, err)

	checksumStr := strings.ToUpper(testdata.ChecksumHackatom)

	specs := map[string]struct {
		args   []string
		expErr bool
	}{
		"source missing": {
			args:   []string{"--builder=" + correctBuilderRef, "--code-hash=" + checksumStr},
			expErr: true,
		},
		"builder missing": {
			args:   []string{"--code-source-url=" + correctSource, "--code-hash=" + checksumStr},
			expErr: true,
		},
		"code hash missing": {
			args:   []string{"--code-source-url=" + correctSource, "--builder=" + correctBuilderRef},
			expErr: true,
		},
		"source format wrong": {
			args:   []string{"--code-source-url=" + "format_wrong", "--builder=" + correctBuilderRef, "--code-hash=" + checksumStr},
			expErr: true,
		},
		"builder format wrong": {
			args:   []string{"--code-source-url=" + correctSource, "--builder=" + "format//", "--code-hash=" + checksumStr},
			expErr: true,
		},
		"code hash wrong": {
			args:   []string{"--code-source-url=" + correctSource, "--builder=" + correctBuilderRef, "--code-hash=" + "AA"},
			expErr: true,
		},
		"happy path, none set": {
			args:   []string{},
			expErr: false,
		},
		"happy path all set": {
			args:   []string{"--code-source-url=" + correctSource, "--builder=" + correctBuilderRef, "--code-hash=" + checksumStr},
			expErr: false,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			flags := ProposalStoreAndInstantiateContractCmd().Flags()
			require.NoError(t, flags.Parse(spec.args))
			_, _, _, gotErr := parseVerificationFlags(wasmBin, flags)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
		})
	}
}
