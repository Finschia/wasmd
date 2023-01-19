package types

import (
	"testing"
	
	"github.com/stretchr/testify/require"
)

const humanAddress = "link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8"

func TestValidateGenesisState(t *testing.T) {
	specs := map[string]struct {
		srcMutator func(state *GenesisState)
		expError   bool
	}{
		"all good": {
			srcMutator: func(s *GenesisState) {},
		},
		"params invalid": {
			srcMutator: func(s *GenesisState) {
				s.Params = Params{}
			},
			expError: true,
		},
		"codeinfo invalid": {
			srcMutator: func(s *GenesisState) {
				s.Codes[0].CodeInfo.CodeHash = nil
			},
			expError: true,
		},
		"contract invalid": {
			srcMutator: func(s *GenesisState) {
				s.Contracts[0].ContractAddress = "invalid"
			},
			expError: true,
		},
		"sequence invalid": {
			srcMutator: func(s *GenesisState) {
				s.Sequences[0].IDKey = nil
			},
			expError: true,
		},
		"genesis store code message invalid": {
			srcMutator: func(s *GenesisState) {
				s.GenMsgs[0].GetStoreCode().WASMByteCode = nil
			},
			expError: true,
		},
		"genesis instantiate contract message invalid": {
			srcMutator: func(s *GenesisState) {
				s.GenMsgs[1].GetInstantiateContract().CodeID = 0
			},
			expError: true,
		},
		"genesis execute contract message invalid": {
			srcMutator: func(s *GenesisState) {
				s.GenMsgs[2].GetExecuteContract().Sender = "invalid"
			},
			expError: true,
		},
		"genesis invalid message type": {
			srcMutator: func(s *GenesisState) {
				s.GenMsgs[0].Sum = nil
			},
			expError: true,
		},
		"inactivate contract address invalid": {
			srcMutator: func(s *GenesisState) {
				s.InactiveContractAddresses[0] = "invalid_address"
			},
			expError: true,
		},
	}

	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			state := GenesisFixture(spec.srcMutator)
			got := state.ValidateBasic()
			if spec.expError {
				require.Error(t, got)
				return
			}
			require.NoError(t, got)
		})
	}
}
