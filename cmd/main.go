package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/urfave/cli/v2"
)

const appName = "zkevm-node"

const (
	// AGGREGATOR is the aggregator component identifier
	AGGREGATOR = "aggregator"
	// SEQUENCER is the sequencer component identifier
	SEQUENCER = "sequencer"
	// RPC is the RPC component identifier
	RPC = "rpc"
	// SYNCHRONIZER is the synchronizer component identifier
	SYNCHRONIZER = "synchronizer"
	// ETHTXMANAGER is the service that manages the tx sent to L1
	ETHTXMANAGER = "eth-tx-manager"
	// L2GASPRICER is the l2 gas pricer component identifier
	L2GASPRICER = "l2gaspricer"
	// SEQUENCE_SENDER is the sequence sender component identifier
	SEQUENCE_SENDER = "sequence-sender"
)

const (
	// NODE_CONFIGFILE name to identify the node config-file
	NODE_CONFIGFILE = "node"
	// NETWORK_CONFIGFILE name to identify the netowk_custom (genesis) config-file
	NETWORK_CONFIGFILE = "custom_network"
)

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: true,
	}
	networkFlag = cli.StringFlag{
		Name:     config.FlagNetwork,
		Aliases:  []string{"net"},
		Usage:    "Load default network configuration. Supported values: [`mainnet`, `testnet`, `custom`]",
		Required: true,
	}
	customNetworkFlag = cli.StringFlag{
		Name:     config.FlagCustomNetwork,
		Aliases:  []string{"net-file"},
		Usage:    "Load the network configuration file if --network=custom",
		Required: false,
	}
	yesFlag = cli.BoolFlag{
		Name:     config.FlagYes,
		Aliases:  []string{"y"},
		Usage:    "Automatically accepts any confirmation to execute the command",
		Required: false,
	}
	componentsFlag = cli.StringSliceFlag{
		Name:     config.FlagComponents,
		Aliases:  []string{"co"},
		Usage:    "List of components to run",
		Required: false,
		Value:    cli.NewStringSlice(AGGREGATOR, SEQUENCER, RPC, SYNCHRONIZER, ETHTXMANAGER, L2GASPRICER, SEQUENCE_SENDER),
	}
	httpAPIFlag = cli.StringSliceFlag{
		Name:     config.FlagHTTPAPI,
		Aliases:  []string{"ha"},
		Usage:    fmt.Sprintf("List of JSON RPC apis to be exposed by the server: --http.api=%v,%v,%v,%v,%v,%v", jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIDebug, jsonrpc.APIZKEVM, jsonrpc.APITxPool, jsonrpc.APIWeb3),
		Required: false,
		Value:    cli.NewStringSlice(jsonrpc.APIEth, jsonrpc.APINet, jsonrpc.APIZKEVM, jsonrpc.APITxPool, jsonrpc.APIWeb3),
	}
	migrationsFlag = cli.BoolFlag{
		Name:     config.FlagMigrations,
		Aliases:  []string{"mig"},
		Usage:    "Blocks the migrations in stateDB to not run them",
		Required: false,
	}
	outputFileFlag = cli.StringFlag{
		Name:     config.FlagOutputFile,
		Usage:    "Indicate the output file",
		Required: true,
	}
	documentationFileTypeFlag = cli.StringFlag{
		Name:     config.FlagDocumentationFileType,
		Usage:    fmt.Sprintf("Indicate the type of file to generate json-schema: %v,%v ", NODE_CONFIGFILE, NETWORK_CONFIGFILE),
		Required: true,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = zkevm.Version
	flags := []cli.Flag{
		&configFileFlag,
		&yesFlag,
		&componentsFlag,
		&httpAPIFlag,
	}
	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Application version and build",
			Action:  versionCmd,
		},
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Run the zkevm-node",
			Action:  start,
			Flags:   append(flags, &networkFlag, &customNetworkFlag, &migrationsFlag),
		},
		{
			Name:    "approve",
			Aliases: []string{"ap"},
			Usage:   "Approve tokens to be spent by the smart contract",
			Action:  approveTokens,
			Flags: append(flags,
				&cli.StringFlag{
					Name:     config.FlagKeyStorePath,
					Aliases:  []string{""},
					Usage:    "the path of the key store file containing the private key of the account going to sign and approve the tokens",
					Required: true,
				},
				&cli.StringFlag{
					Name:     config.FlagPassword,
					Aliases:  []string{"pw"},
					Usage:    "the password do decrypt the key store file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     config.FlagAmount,
					Aliases:  []string{"am"},
					Usage:    "Amount that is gonna be approved",
					Required: false,
				},
				&cli.StringFlag{
					Name:     config.FlagMaxAmount,
					Aliases:  []string{"mam"},
					Usage:    "Maximum amount is gonna be approved",
					Required: false,
				},
				&networkFlag,
				&customNetworkFlag,
			),
		},
		{
			Name:    "encryptKey",
			Aliases: []string{},
			Usage:   "Encrypts the privatekey with a password and create a keystore file",
			Action:  encryptKey,
			Flags:   encryptKeyFlags,
		},
		{
			Name:    "dumpState",
			Aliases: []string{},
			Usage:   "Dumps the state in a JSON file, for debug purposes",
			Action:  dumpState,
			Flags:   dumpStateFlags,
		},
		{
			Name:   "generate-json-schema",
			Usage:  "Generate the json-schema for the configuration file, and store it on docs/schema.json",
			Action: genJSONSchema,
			Flags:  []cli.Flag{&outputFileFlag, &documentationFileTypeFlag},
		},
		{
			Name:    "snapshot",
			Aliases: []string{"snap"},
			Usage:   "Snapshot the state db",
			Action:  snapshot,
			Flags:   snapshotFlags,
		},
		{
			Name:    "restore",
			Aliases: []string{},
			Usage:   "Restore snapshot of the state db",
			Action:  restore,
			Flags:   restoreFlags,
		},
		{
			Name:    "update-genesis-root",
			Aliases: []string{"ugr"},
			Usage:   "Update the root value in genesis config file",
			Action:  updateGenesisRoot,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Path to genesis config file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "root",
					Aliases:  []string{"r"},
					Usage:    "New root value to set",
					Required: true,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type GenesisConfig struct {
	L1Config struct {
		ChainID                           int    `json:"chainId"`
		PolygonZkEVMAddress               string `json:"polygonZkEVMAddress"`
		MaticTokenAddress                 string `json:"maticTokenAddress"`
		PolygonZkEVMGlobalExitRootAddress string `json:"polygonZkEVMGlobalExitRootAddress"`
	} `json:"l1Config"`
	GenesisBlockNumber int    `json:"genesisBlockNumber"`
	Root               string `json:"root"`
	Genesis            struct {
		GenesisActions []struct {
			Address  string `json:"address"`
			Type     int    `json:"type"`
			Value    string `json:"value"`
			Bytecode string `json:"bytecode,omitempty"`
		} `json:"genesisActions"`
	} `json:"genesis"`
}

func updateGenesisRoot(ctx *cli.Context) error {
	configPath := ctx.String("config")
	newRoot := ctx.String("root")

	if configPath == "" {
		return fmt.Errorf("config path is required")
	}

	if newRoot == "" {
		return fmt.Errorf("new root value is required")
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析JSON
	var config GenesisConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// 备份原文件
	backupPath := configPath + ".backup"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}

	fmt.Printf("Backup created at: %s\n", backupPath)

	// 更新root值
	oldRoot := config.Root
	config.Root = newRoot

	// 重新序列化
	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write updated config: %v", err)
	}

	fmt.Printf("Successfully updated genesis root:\n")
	fmt.Printf("  Old root: %s\n", oldRoot)
	fmt.Printf("  New root: %s\n", newRoot)
	fmt.Printf("  Config file: %s\n", configPath)

	return nil
}
