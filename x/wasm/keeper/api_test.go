package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	wasmtype "github.com/line/wasmd/x/wasm/types"
	wasmvm "github.com/line/wasmvm"

	wasmvmtypes "github.com/line/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAPI(t *testing.T) wasmvm.GoAPI {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)
	return keepers.WasmKeeper.cosmwasmAPI(ctx)
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
		assert.Equal(t, wasmtype.DefaultGasMultiplier*5, gas)
	})

	t.Run("invalid address", func(t *testing.T) {
		_, gas, err := api.HumanAddress([]byte("invalid_address"))
		require.Error(t, err)
		assert.Equal(t, wasmtype.DefaultGasMultiplier*5, gas)
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
		assert.Equal(t, wasmtype.DefaultGasMultiplier*4, gas)
	})

	t.Run("invalid address", func(t *testing.T) {
		_, gas, err := api.CanonicalAddress("invalid_address")
		assert.Error(t, err)
		assert.Equal(t, wasmtype.DefaultGasMultiplier*4, gas)
	})
}

func TestAPIGetContractEnv(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	// instantiate a number contract
	numberWasm, err := ioutil.ReadFile("../testdata/number.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit...)
	em := sdk.NewEventManager()
	codeID, err := keepers.ContractKeeper.Create(ctx.WithEventManager(em), creator, numberWasm, nil)
	require.NoError(t, err)
	value := 42
	initMsg := []byte(fmt.Sprintf(`{"value":%d}`, value))
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx.WithEventManager(em), codeID, creator, nil, initMsg, "number", nil)
	require.NoError(t, err)
	msgLen := 101010

	// prepare API
	api := keepers.WasmKeeper.cosmwasmAPI(ctx)

	t.Run("succeed", func(t *testing.T) {
		// omitted value is MultipliedGasMeter. It is not tested here.
		env, _, store, querier, _, _, instantiateCost, gas, err := api.GetContractEnv(contractAddr.String(), uint64(msgLen))

		require.NoError(t, err)

		assert.Equal(t, uint64(ctx.BlockHeight()), env.Block.Height)
		assert.Equal(t, uint64(ctx.BlockTime().UnixNano()), env.Block.Time)
		assert.Equal(t, ctx.ChainID(), env.Block.ChainID)
		assert.Equal(t, contractAddr.String(), env.Contract.Address)

		// "number" comes from https://github.com/line/cosmwasm/blob/d08b5a59115cc3d28f375b7425b902bfd1dac6a4/contracts/number/src/contract.rs#L9
		assert.Equal(t, []byte{uint8(0), uint8(0), uint8(0), uint8(value)}, store.Get([]byte("number")))

		queryMsg := []byte(`{"number":{}}`)
		query := wasmvmtypes.QueryRequest{
			Wasm: &wasmvmtypes.WasmQuery{
				Smart: &wasmvmtypes.SmartQuery{
					ContractAddr: contractAddr.String(),
					Msg:          queryMsg,
				},
			},
		}
		queryResult, err := querier.Query(query, 10_000_000_000_000)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"value":42}`), queryResult)

		expectedInstantiateCost := keepers.WasmKeeper.instantiateContractCosts(keepers.WasmKeeper.gasRegister, ctx, false, msgLen)
		assert.Equal(t, wasmtype.DefaultGasMultiplier*expectedInstantiateCost, instantiateCost)

		assert.Equal(t, wasmtype.DefaultGasMultiplier*11, gas)
	})

	t.Run("non-existed contract", func(t *testing.T) {
		nonExistedContractAddr := "link1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqyu0w3p"
		require.NotEqual(t, nonExistedContractAddr, contractAddr)
		_, _, _, _, _, _, _, _, err := api.GetContractEnv(nonExistedContractAddr, uint64(msgLen))
		require.Error(t, err)
	})
}

func TestCallCallablePoint(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	// instantiate an events contract
	numberWasm, err := ioutil.ReadFile("../testdata/events.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit...)
	em := sdk.NewEventManager()
	codeID, err := keepers.ContractKeeper.Create(ctx.WithEventManager(em), creator, numberWasm, nil)
	require.NoError(t, err)
	initMsg := []byte(`{}`)
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx.WithEventManager(em), codeID, creator, nil, initMsg, "events", nil)
	require.NoError(t, err)
	callstack := []sdk.AccAddress{RandomAccountAddress(t), RandomAccountAddress(t)}
	callstackBin, err := json.Marshal(callstack)
	require.NoError(t, err)
	var gasLimit uint64 = keepers.WasmKeeper.getGasMultiplier(ctx).ToWasmVMGas(400_000)

	// prepare API
	api := keepers.WasmKeeper.cosmwasmAPI(ctx.WithEventManager(em))

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

		eventsExpected, err := newCustomCallablePointEvents(eventsIn, contractAddr, callstackBin)
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
}

func TestValidateDynamicLinkInterface(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	// instantiate an events contract
	numberWasm, err := ioutil.ReadFile("../testdata/events.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit...)
	em := sdk.NewEventManager()
	codeID, err := keepers.ContractKeeper.Create(ctx.WithEventManager(em), creator, numberWasm, nil)
	require.NoError(t, err)
	initMsg := []byte(`{}`)
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx.WithEventManager(em), codeID, creator, nil, initMsg, "events", nil)
	require.NoError(t, err)

	// prepare API
	api := keepers.WasmKeeper.cosmwasmAPI(ctx.WithEventManager(em))

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
}
