# Steps to run/develop on the environment locally

# Warning:

>Currently the Executor/Prover does not run on ARM-powered Macs. For Windows users, WSL/WSL2 use is not recommended. 

## Overview

This documentation will help you running the following components:

- zkEVM Node Databases
- Explorer Databases
- L1 Network
- Prover
- zkEVM Node components
- Explorers

## Requirements

The current version of the environment requires `go`, `docker` and `docker-compose` to be previously installed, check the links below to understand how to install them:

- <https://go.dev/doc/install>
- <https://www.docker.com/get-started>
- <https://docs.docker.com/compose/install/>

The `zkevm-node` docker image must be built at least once and every time a change is made to the code.
If you haven't build the `zkevm-node` image yet, you must run:

```bash
make build-docker
```

## A look at how the binary works:

The `zkevm-node` allows certain commands to interact with smart contracts, run certain components, create encryption files and print out debug information.

To interact with the binary program we provide docker compose files, and a Makefile to spin up/down the different services and components, ensuring a smooth deployment locally and better interface in command line for developers.

## Controlling the environment

> All the data is stored inside of each docker container, this means once you remove the container, the data will be lost.

To run the environment:

The `test/` directory contains scripts and files for developing and debugging. The most commonly used conf files are as following, which can control the network behavior:
```bash
config/test.genesis.config.json
config/test.node.config.toml
config/test.prover.config.json
.env
docker-compose.yml
Makefile
```


```bash
cd test/
```

Then:

```bash
make run
```

To stop the environment:

```bash
make stop
```

To restart the environment:

```bash
make restart
```

## Sample data

The `make run` will execute the containers needed to run the environment but this will not execute anything else, so the L2 will be basically empty.

If you need sample data already deployed to the network, we have the following scripts:

**To add some examples of transactions and smart contracts:**

```bash
make deploy-sc
```

**To deploy a full a uniswap environment:**

```bash
make deploy-uniswap
```

**To grant the Matic smart contract a set amount of tokens, run:**

```bash
make run-approve-matic
```

## Accessing the environment

- **Databases**:
  - zkEVM Node *State* Database 
    - `Type:` Postgres DB
    - `User:` state_user
    - `Password:` state_password
    - `Database:` state-db
    - `Host:` localhost
    - `Port:` 5432
    - `Url:` <postgres://state_user:srare_password@localhost:5432/state-db>
  - zkEVM Node *Pool* Database 
    - `Type:` Postgres DB
    - `User:` pool_user
    - `Password:` pool_password
    - `Database:` pool_db
    - `Host:` localhost
    - `Port:` 5433
    - `Url:` <postgres://pool_user:pool_password@localhost:5433/pool_db>
  - zkEVM Node *JSON-RPC* Database 
    - `Type:` Postgres DB
    - `User:` rpc_user
    - `Password:` rpc_password
    - `Database:` rpc_db
    - `Host:` localhost
    - `Port:` 5434
    - `Url:` <postgres://rpc_user:rpc_password@localhost:5434/rpc_db>
  - Explorer L1 Database
    - `Type:` Postgres DB
    - `User:` l1_explorer_user
    - `Password:` l1_explorer_password
    - `Database:` l1_explorer_db
    - `Host:` localhost
    - `Port:` 5435
    - `Url:` <postgres://l1_explorer_user:l1_explorer_password@localhost:5435/l1_explorer_db>
  - Explorer L2 Database
    - `Type:` Postgres DB
    - `User:` l2_explorer_user
    - `Password:` l2_explorer_password
    - `Database:` l2_explorer_db
    - `Host:` localhost
    - `Port:` 5436
    - `Url:` <postgres://l2_explorer_user:l2_explorer_password@localhost:5436/l2_explorer_db>
- **Networks**:
  - L1 Network
    - `Type:` Geth
    - `Host:` localhost
    - `Port:` 8545
    - `Url:` <http://localhost:8545>
  - zkEVM Node
    - `Type:` JSON RPC
    - `Host:` localhost
    - `Port:` 8123
    - `Url:` <http://localhost:8123>
- **Explorers**:
  - Explorer L1
    - `Type:` Web
    - `Host:` localhost
    - `Port:` 4000
    - `Url:` <http://localhost:4000>
  - Explorer L2
    - `Type:` Web
    - `Host:` localhost
    - `Port:` 4001
    - `Url:` <http://localhost:4001>
- Prover
  - `Type:` Mock
  - `Host:` localhost
  - `Port:` Depending on the prover image, if it's mock or not: 
    - Prod prover: 50052 for Prover, 50061 for Merkle Tree, 50071 for Executor
    - Mock prover: 43061 for MT, 43071 for Executor
  - `Url:` <http://localhost:50001>
## Metamask

> Metamask requires the network to be running while configuring it, so make sure your network is running before starting.

To configure your Metamask to use your local environment, follow these steps:

1. Log in to your Metamask wallet
2. Click on your account picture and then on Settings
3. On the left menu, click on Networks
4. Click on `Add Network` button
5. Fill up the L2 network information
    1. `Network Name:` Polygon zkEVM - Local
    2. `New RPC URL:` <http://localhost:8123>
    3. `ChainID:` 102
    4. `Currency Symbol:` ETH
    5. `Block Explorer URL:` <http://localhost:4000>
6. Click on Save
7. Click on `Add Network` button
8. Fill up the L1 network information
    1. `Network Name:` Geth - Local
    2. `New RPC URL:` <http://localhost:8545>
    3. `ChainID:` 1002
    4. `Currency Symbol:` ETH
9. Click on Save

## L1 Addresses
1. can view config/test.genesis.config.json, for example 
    ```bash
    intel-10400f test git:(dev) ✗ jq .l1Config config/test.genesis.config.json
    {
      "chainId": 102,
      "polygonZkEVMAddress": "0x67d269191c92Caf3cD7723F116c85e6E9bf55933",
      "maticTokenAddress": "0x3Aa5ebB10DC797CAC828524e59A333d0A371443c",
      "polygonZkEVMGlobalExitRootAddress": "0x09635F643e140090A9A8Dcd712eD6285858ceBef"
    }
    ```

## L2 Addresses
1. Deployer/Sequencer Account, see https://github.com/b2network/b2-zkevm-contracts/blob/dev/docker/scripts/deploy_parameters_docker.json

## L2 GenesisAccounts
1. can view config/test.genesis.config.json, for example 
    ```bash
    intel-10400f test git:(dev) ✗ jq '.genesis[].address'  config/test.genesis.config.json
    "0x4b2700570f8426A24EA85e0324611E527BdD55B8"
    "0xf065BaE7C019ff5627E09ed48D4EeA317D211956"
    "0xf23919bb44BCa81aeAb4586BE71Ee3fd4E99B951"
    "0xff0EE8ea08cEf5cb4322777F5CC3E8A584B8A4A0"
    "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
    "0xa40d5f56745a118d0906a34e69aec8c0db1cb8fa"
    "0x0165878A594ca255338adfa4d48449f69242Eb8F"
    "0x20E7077d25fe79C5F6c2D3ae4905E96aA7C89c13"
    "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
    ```
## Start L1/L2 explorer
1. clone repo https://github.com/b2network/blockscout/
1. checkout dev branch
2. cd docker-compose
3. `bash helper.sh restartL1L2` will  start l1,l2 exploler
   1. access l1 explorer via port 30000
   2. access l2 explorer via port 30001