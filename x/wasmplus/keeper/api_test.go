package keeper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	wasmvm "github.com/Finschia/wasmvm"

	"github.com/Finschia/wasmd/x/wasm/keeper"
	wasmkeeper "github.com/Finschia/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/Finschia/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAPI(t *testing.T) wasmvm.GoAPI {
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	return keepers.WasmKeeper.GetCosmwasmAPIGenerator().Generate(&ctx)
}

func TestAPIHumanAddress(t *testing.T) {
	// prepare API
	api := newAPI(t)

	t.Run("valid address", func(t *testing.T) {
		// address for alice in testnet
		addr := "link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705"
		bz, err := sdk.AccAddressFromBech32(addr)
		require.NoError(t, err)
		result, gas, err := api.HumanAddress(bz)
		require.NoError(t, err)
		assert.Equal(t, addr, result)
		assert.Equal(t, wasmkeeper.CostHumanize, gas)
	})

	t.Run("invalid address", func(t *testing.T) {
		_, gas, err := api.HumanAddress([]byte("invalid_address"))
		require.Error(t, err)
		assert.Equal(t, wasmkeeper.CostHumanize, gas)
	})
}

func TestAPICanonicalAddress(t *testing.T) {
	// prepare API
	api := newAPI(t)

	t.Run("valid address", func(t *testing.T) {
		addr := "link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705"
		expected, err := sdk.AccAddressFromBech32(addr)
		require.NoError(t, err)
		result, gas, err := api.CanonicalAddress(addr)
		require.NoError(t, err)
		assert.Equal(t, expected.Bytes(), result)
		assert.Equal(t, wasmkeeper.CostCanonical, gas)
	})

	t.Run("invalid address", func(t *testing.T) {
		_, gas, err := api.CanonicalAddress("invalid_address")
		assert.Error(t, err)
		assert.Equal(t, wasmkeeper.CostCanonical, gas)
	})
}

func TestCallCallablePoint(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	em := sdk.NewEventManager()
	ctx = ctx.WithEventManager(em)

	// instantiate an events contract
	numberWasm, err := ioutil.ReadFile("../testdata/events.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedRandomAccount(ctx, deposit...)
	codeID, _, err := keepers.ContractKeeper.Create(ctx, creator, numberWasm, nil)
	require.NoError(t, err)
	initMsg := []byte(`{}`)
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, nil, initMsg, "events", nil)
	require.NoError(t, err)
	callstack := []sdk.AccAddress{RandomAccountAddress(t), RandomAccountAddress(t)}
	callstackBin, err := json.Marshal(callstack)
	require.NoError(t, err)
	var gasLimit uint64 = keepers.WasmKeeper.GetGasRegister().ToWasmVMGas(400_000)

	// prepare API
	api := keepers.WasmKeeper.GetCosmwasmAPIGenerator().Generate(&ctx)

	// prepare arg for succeed
	eventsIn := wasmvmtypes.Events{
		wasmvmtypes.Event{
			Type: "ty1",
			Attributes: wasmvmtypes.EventAttributes{
				wasmvmtypes.EventAttribute{
					Key:   "alice",
					Value: "101010",
				},
				wasmvmtypes.EventAttribute{
					Key:   "bob",
					Value: "42",
				},
			},
		},
		wasmvmtypes.Event{
			Type: "ty2",
			Attributes: wasmvmtypes.EventAttributes{
				wasmvmtypes.EventAttribute{
					Key:   "ALICE",
					Value: "42",
				},
				wasmvmtypes.EventAttribute{
					Key:   "BOB",
					Value: "101010",
				},
			},
		},
	}
	eventsInBin, err := eventsIn.MarshalJSON()
	require.NoError(t, err)

	t.Run("succeed", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		res, _, err := api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, false, callstackBin, gasLimit)
		require.NoError(t, err)
		assert.Equal(t, []byte(`null`), res)

		eventsExpected, err := keeper.NewCustomCallablePointEvents(eventsIn, contractAddr, callstackBin)
		require.NoError(t, err)
		for _, e := range eventsExpected {
			assert.Contains(t, em.Events(), e)
		}
	})

	t.Run("fail with no arg", func(t *testing.T) {
		argsEv := [][]byte{}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, false, callstackBin, gasLimit)
		require.Error(t, err)
		assert.ErrorContains(t, err, "RuntimeError")
		assert.ErrorContains(t, err, "Parameters of type [I32] did not match signature [I32, I32] -> []")
	})

	t.Run("fail with invalid name", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "invalid"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, false, callstackBin, gasLimit)

		// fail to get permission
		require.Error(t, err)
		assert.ErrorContains(t, err, "Error during calling dynamic linked callable point")
		require.ErrorContains(t, err, "callee function properties has not key:invalid")
	})

	t.Run("fail with invalid address", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(RandomAccountAddress(t).String(), nameBin, argsEvBin, false, callstackBin, gasLimit)

		require.Error(t, err)
		assert.ErrorContains(t, err, "contract: not found")
	})

	t.Run("fail with lack of write permission", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, true, callstackBin, gasLimit)

		require.Error(t, err)
		assert.ErrorContains(t, err, "Error during calling dynamic linked callable point")
		assert.ErrorContains(t, err, "a read-write callable point is called in read-only context.")
	})

	t.Run("fail with re-entrancing", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)

		// callstack with re-entrancy
		callstack = append(callstack, contractAddr)
		callstackBin, err := json.Marshal(callstack)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, true, callstackBin, gasLimit)

		require.Error(t, err)
		assert.ErrorContains(t, err, "Error calling the VM")
		assert.ErrorContains(t, err, "A contract can only be called once per one call stack.")
	})

	t.Run("fail with inactive contract", func(t *testing.T) {
		// add contract to inactive
		keepers.WasmKeeper.addInactiveContract(ctx, contractAddr)

		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		_, _, err = api.CallCallablePoint(contractAddr.String(), nameBin, argsEvBin, false, callstackBin, gasLimit)
		assert.ErrorContains(t, err, "called contract cannot be executed")

		// reset inactive contracts
		keepers.WasmKeeper.deleteInactiveContract(ctx, contractAddr)
	})

	t.Run("fail with invalid callee address", func(t *testing.T) {
		argsEv := [][]byte{eventsInBin}
		argsEvBin, err := json.Marshal(argsEv)
		require.NoError(t, err)
		name := "add_events_dyn"
		nameBin, err := json.Marshal(name)
		require.NoError(t, err)
		invalidAddr := "invalidAddr"
		_, _, err = api.CallCallablePoint(invalidAddr, nameBin, argsEvBin, false, callstackBin, gasLimit)
		assert.ErrorContains(t, err, "specified callee address is invalid")
	})
}

func TestValidateDynamicLinkInterface(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	em := sdk.NewEventManager()
	ctx = ctx.WithEventManager(em)

	// instantiate an events contract
	numberWasm, err := ioutil.ReadFile("../testdata/events.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedRandomAccount(ctx, deposit...)
	codeID, _, err := keepers.ContractKeeper.Create(ctx, creator, numberWasm, nil)
	require.NoError(t, err)
	initMsg := []byte(`{}`)
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, nil, initMsg, "events", nil)
	require.NoError(t, err)

	// prepare API
	api := keepers.WasmKeeper.GetCosmwasmAPIGenerator().Generate(&ctx)

	t.Run("succeed valid", func(t *testing.T) {
		validInterface := []byte(`[{"name":"add_event_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_events_dyn","ty":{"params":["I32","I32"],"results":[]}},{"name":"add_attribute_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_attributes_dyn","ty":{"params":["I32","I32"],"results":[]}}]`)
		res, _, err := api.ValidateInterface(contractAddr.String(), validInterface)

		require.NoError(t, err)
		assert.Equal(t, []byte(`null`), res)
	})

	t.Run("succeed invalid", func(t *testing.T) {
		invalidInterface := []byte(`[{"name":"add_event","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_events","ty":{"params":["I32","I32"],"results":[]}},{"name":"add_attribute","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_attributes","ty":{"params":["I32","I32"],"results":[]}}]`)
		res, _, err := api.ValidateInterface(contractAddr.String(), invalidInterface)

		require.NoError(t, err)
		assert.Contains(t, string(res), `following functions are not implemented`)
		assert.Contains(t, string(res), `add_event`)
		assert.Contains(t, string(res), `add_events`)
		assert.Contains(t, string(res), `add_attribute`)
		assert.Contains(t, string(res), `add_attributes`)
	})

	t.Run("fail with invalid address", func(t *testing.T) {
		validInterface := []byte(`[{"name":"add_event_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_events_dyn","ty":{"params":["I32","I32"],"results":[]}},{"name":"add_attribute_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_attributes_dyn","ty":{"params":["I32","I32"],"results":[]}}]`)
		_, _, err := api.ValidateInterface(RandomAccountAddress(t).String(), validInterface)
		require.Error(t, err)
		assert.ErrorContains(t, err, "contract: not found")
	})

	t.Run("fail with inactive contract", func(t *testing.T) {
		// add contract to inactive
		keepers.WasmKeeper.addInactiveContract(ctx, contractAddr)

		validInterface := []byte(`[{"name":"add_event_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_events_dyn","ty":{"params":["I32","I32"],"results":[]}},{"name":"add_attribute_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_attributes_dyn","ty":{"params":["I32","I32"],"results":[]}}]`)
		_, _, err = api.ValidateInterface(contractAddr.String(), validInterface)

		assert.ErrorContains(t, err, "try to validate a contract cannot be executed")

		// reset inactive contracts
		keepers.WasmKeeper.deleteInactiveContract(ctx, contractAddr)
	})

	t.Run("fail with invalid contract address", func(t *testing.T) {
		validInterface := []byte(`[{"name":"add_event_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_events_dyn","ty":{"params":["I32","I32"],"results":[]}},{"name":"add_attribute_dyn","ty":{"params":["I32","I32","I32"],"results":[]}},{"name":"add_attributes_dyn","ty":{"params":["I32","I32"],"results":[]}}]`)
		invalidAddr := "invalidAddr"
		_, _, err = api.ValidateInterface(invalidAddr, validInterface)

		assert.ErrorContains(t, err, "specified contract address is invalid")
	})
}
