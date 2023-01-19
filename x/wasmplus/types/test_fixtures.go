package types

import (
	"math/rand"

	wasmtypes "github.com/line/wasmd/x/wasm/types"
)

func GenesisFixture(mutators ...func(*GenesisState)) GenesisState {
	const (
		numCodes     = 2
		numContracts = 2
		numSequences = 2
		numMsg       = 3
	)

	fixture := GenesisState{
		Params:                    DefaultParams(),
		Codes:                     make([]wasmtypes.Code, numCodes),
		Contracts:                 make([]wasmtypes.Contract, numContracts),
		Sequences:                 make([]wasmtypes.Sequence, numSequences),
		InactiveContractAddresses: make([]string, numContracts),
	}
	for i := 0; i < numCodes; i++ {
		fixture.Codes[i] = wasmtypes.CodeFixture()
	}
	for i := 0; i < numContracts; i++ {
		contract := wasmtypes.ContractFixture()
		fixture.Contracts[i] = contract
		fixture.InactiveContractAddresses[i] = contract.ContractAddress
	}
	for i := 0; i < numSequences; i++ {
		fixture.Sequences[i] = wasmtypes.Sequence{
			IDKey: randBytes(5),
			Value: uint64(i),
		}
	}
	fixture.GenMsgs = []wasmtypes.GenesisState_GenMsgs{
		{Sum: &wasmtypes.GenesisState_GenMsgs_StoreCode{StoreCode: wasmtypes.MsgStoreCodeFixture()}},
		{Sum: &wasmtypes.GenesisState_GenMsgs_InstantiateContract{InstantiateContract: wasmtypes.MsgInstantiateContractFixture()}},
		{Sum: &wasmtypes.GenesisState_GenMsgs_ExecuteContract{ExecuteContract: wasmtypes.MsgExecuteContractFixture()}},
	}
	for _, m := range mutators {
		m(&fixture)
	}
	return fixture
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	rand.Read(r)
	return r
}
