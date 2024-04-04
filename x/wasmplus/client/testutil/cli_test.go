//go:build cli_test
// +build cli_test

package testutil

import (
	"fmt"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/suite"

	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	wasmapp "github.com/Finschia/wasmd/app"
	wasmappparams "github.com/Finschia/wasmd/app/params"
	wasmplusapp "github.com/Finschia/wasmd/appplus"
	"github.com/Finschia/wasmd/x/wasmplus/types"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := DefaultConfig()
	cfg.NumValidators = 1
	suite.Run(t, NewIntegrationTestSuite(cfg))
}

func DefaultConfig() network.Config {
	encCfg := wasmplusapp.MakeEncodingConfig()

	return network.Config{
		Codec:             encCfg.Marshaler,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		TxConfig:          encCfg.TxConfig,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(encCfg),
		GenesisState:      wasmplusapp.ModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:     1 * time.Second,
		ChainID:           "chain-" + rand.NewRand().Str(6),
		NumValidators:     4,
		BondDenom:         sdk.DefaultBondDenom,
		MinGasPrices:      fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:     sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:     sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:      sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy:   storetypes.PruningOptionNothing,
		CleanupDir:        true,
		SigningAlgo:       string(hd.Secp256k1Type),
		KeyringOptions:    []keyring.Option{},
	}
}

func NewAppConstructor(encodingCfg wasmappparams.EncodingConfig) network.AppConstructor {
	return func(val network.Validator) servertypes.Application {
		return wasmplusapp.NewWasmApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true,
			map[int64]bool{}, val.Ctx.Config.RootDir, 0,
			encodingCfg,
			types.EnableAllProposals,
			wasmapp.EmptyBaseAppOptions{},
			nil,
		)
	}
}
