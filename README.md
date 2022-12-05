# A backend dapp demo

## Prepare

```shell
solc --combined-json abi,bin,userdoc,devdoc --optimize -o . --overwrite ERC20.sol

abigen --combined-json combined.json --pkg main --out erc20.go
```

Config account key, and the `rpcUrl` in the config.json file.

> If you want to use the `eth_subscribe` feature, then the `rpcUrl` must be a websocket url or an IPC file path.
> The http schema do not support subscription.
