package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	wasmcli "github.com/Finschia/wasmd/x/wasm/client/cli"
	"github.com/Finschia/wasmd/x/wasm/client/cli/os"
	"github.com/Finschia/wasmd/x/wasm/ioutils"
	wasmTypes "github.com/Finschia/wasmd/x/wasm/types"
	"github.com/Finschia/wasmd/x/wasmplus/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagAmount                    = "amount"
	flagLabel                     = "label"
	flagAdmin                     = "admin"
	flagInstantiateByEverybody    = "instantiate-everybody"
	flagInstantiateNobody         = "instantiate-nobody"
	flagInstantiateByAddress      = "instantiate-only-address"
	flagInstantiateByAnyOfAddress = "instantiate-anyof-addresses"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Wasm transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		wasmcli.StoreCodeCmd(),
		wasmcli.InstantiateContractCmd(),
		wasmcli.InstantiateContract2Cmd(),
		StoreCodeAndInstantiateContractCmd(),
		wasmcli.ExecuteContractCmd(),
		wasmcli.MigrateContractCmd(),
		wasmcli.UpdateContractAdminCmd(),
		wasmcli.ClearContractAdminCmd(),
	)
	return txCmd
}

// StoreCodeAndInstantiateContractCmd will upload code and instantiate a contract using it
func StoreCodeAndInstantiateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store-instantiate [wasm file] [json_encoded_init_args] --label [text] --admin [address,optional] --amount [coins,optional]",
		Short: "Upload a wasm binary and instantiate a wasm contract from the code",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg, err := parseStoreCodeAndInstantiateContractArgs(args[0], args[1], clientCtx.GetFromAddress(), cmd.Flags())
			if err != nil {
				return err
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagInstantiateByEverybody, "", "Everybody can instantiate a contract from the code, optional")
	cmd.Flags().String(flagInstantiateNobody, "", "Nobody except the governance process can instantiate a contract from the code, optional")
	cmd.Flags().String(flagInstantiateByAddress, "", fmt.Sprintf("Removed: use %s instead", flagInstantiateByAnyOfAddress))
	cmd.Flags().StringSlice(flagInstantiateByAnyOfAddress, []string{}, "Any of the addresses can instantiate a contract from the code, optional")
	cmd.Flags().String(flagAmount, "", "Coins to send to the contract during instantiation")
	cmd.Flags().String(flagLabel, "", "A human-readable name for this contract in lists")
	cmd.Flags().String(flagAdmin, "", "Address of an admin")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseStoreCodeAndInstantiateContractArgs(file string, initMsg string, sender sdk.AccAddress, flags *flag.FlagSet) (*types.MsgStoreCodeAndInstantiateContract, error) {
	wasm, err := os.ReadFileWithSizeLimit(file, int64(wasmTypes.MaxWasmSize))
	if err != nil {
		return nil, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return nil, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return nil, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	var perm *wasmTypes.AccessConfig
	onlyAddrStr, err := flags.GetString(flagInstantiateByAddress)
	if err != nil {
		return nil, fmt.Errorf("instantiate by address: %s", err)
	}
	if onlyAddrStr != "" {
		return nil, fmt.Errorf("not supported anymore. Use: %s", flagInstantiateByAnyOfAddress)
	}
	everybodyStr, err := flags.GetString(flagInstantiateByEverybody)
	if err != nil {
		return nil, fmt.Errorf("instantiate by everybody: %s", err)
	}
	if everybodyStr != "" {
		ok, err := strconv.ParseBool(everybodyStr)
		if err != nil {
			return nil, fmt.Errorf("boolean value expected for instantiate by everybody: %s", err)
		}
		if ok {
			perm = &wasmTypes.AllowEverybody
		}
	}

	amountStr, err := flags.GetString(flagAmount)
	if err != nil {
		return nil, fmt.Errorf("amount: %s", err)
	}
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return nil, fmt.Errorf("amount: %s", err)
	}
	label, err := flags.GetString(flagLabel)
	if err != nil {
		return nil, fmt.Errorf("label: %s", err)
	}
	if label == "" {
		return nil, errors.New("label is required on all contracts")
	}
	adminStr, err := flags.GetString(flagAdmin)
	if err != nil {
		return nil, fmt.Errorf("admin: %s", err)
	}

	msg := types.MsgStoreCodeAndInstantiateContract{
		Sender:                sender.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: perm,
		Label:                 label,
		Funds:                 amount,
		Msg:                   []byte(initMsg),
		Admin:                 adminStr,
	}
	return &msg, nil
}
