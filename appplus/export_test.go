package appplus

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
)

func TestZeroHeightGenesis(t *testing.T) {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	db := dbm.NewMemDB()
	gapp := NewWasmApp(log.NewNopLogger(), db, nil, true, simtestutil.NewAppOptionsWithFlagHome(dir), nil)

	genesisState := gapp.DefaultGenesis()
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	// Initialize the chain
	_, err = gapp.InitChain(
		&abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	require.NoError(t, err)
	gapp.Commit()

	jailAllowedAddress := []string{"linkvaloper12kr02kew9fl73rqekalavuu0xaxcgwr6pz5vt8"}
	_, err = gapp.ExportAppStateAndValidators(true, jailAllowedAddress, nil)
	require.NoError(t, err)
}
