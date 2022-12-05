package main

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// EIP20MetaData contains all meta data concerning the EIP20 contract.
var EIP20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_initialAmount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_tokenName\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimalUnits\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_tokenSymbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"remaining\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516108de3803806108de8339818101604052608081101561003357600080fd5b81516020830180516040519294929383019291908464010000000082111561005a57600080fd5b90830190602082018581111561006f57600080fd5b825164010000000081118282018810171561008957600080fd5b82525081516020918201929091019080838360005b838110156100b657818101518382015260200161009e565b50505050905090810190601f1680156100e35780820380516001836020036101000a031916815260200191505b5060408181526020830151920180519294919391928464010000000082111561010b57600080fd5b90830190602082018581111561012057600080fd5b825164010000000081118282018810171561013a57600080fd5b82525081516020918201929091019080838360005b8381101561016757818101518382015260200161014f565b50505050905090810190601f1680156101945780820380516001836020036101000a031916815260200191505b5060409081523360009081526001602090815291812089905588905586516101c594506003935090870191506101f2565b506004805460ff191660ff841617905580516101e89060059060208401906101f2565b5050505050610285565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061023357805160ff1916838001178555610260565b82800160010185558215610260579182015b82811115610260578251825591602001919060010190610245565b5061026c929150610270565b5090565b5b8082111561026c5760008155600101610271565b61064a806102946000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063313ce56711610071578063313ce567146101e15780635c658165146101ff57806370a082311461022d57806395d89b4114610253578063a9059cbb1461025b578063dd62ed3e14610287576100a9565b806306fdde03146100ae578063095ea7b31461012b57806318160ddd1461016b57806323b872dd1461018557806327e235e3146101bb575b600080fd5b6100b66102b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f05781810151838201526020016100d8565b50505050905090810190601f16801561011d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101576004803603604081101561014157600080fd5b506001600160a01b038135169060200135610343565b604080519115158252519081900360200190f35b6101736103a9565b60408051918252519081900360200190f35b6101576004803603606081101561019b57600080fd5b506001600160a01b038135811691602081013590911690604001356103af565b610173600480360360208110156101d157600080fd5b50356001600160a01b03166104b1565b6101e96104c3565b6040805160ff9092168252519081900360200190f35b6101736004803603604081101561021557600080fd5b506001600160a01b03813581169160200135166104cc565b6101736004803603602081101561024357600080fd5b50356001600160a01b03166104e9565b6100b6610504565b6101576004803603604081101561027157600080fd5b506001600160a01b03813516906020013561055f565b6101736004803603604081101561029d57600080fd5b506001600160a01b03813581169160200135166105e9565b6003805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561033b5780601f106103105761010080835404028352916020019161033b565b820191906000526020600020905b81548152906001019060200180831161031e57829003601f168201915b505050505081565b3360008181526002602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60005481565b6001600160a01b0383166000818152600260209081526040808320338452825280832054938352600190915281205490919083118015906103f05750828110155b6103f957600080fd5b6001600160a01b038085166000908152600160205260408082208054870190559187168152208054849003905560001981101561045b576001600160a01b03851660009081526002602090815260408083203384529091529020805484900390555b836001600160a01b0316856001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a3506001949350505050565b60016020526000908152604090205481565b60045460ff1681565b600260209081526000928352604080842090915290825290205481565b6001600160a01b031660009081526001602052604090205490565b6005805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561033b5780601f106103105761010080835404028352916020019161033b565b3360009081526001602052604081205482111561057b57600080fd5b336000818152600160209081526040808320805487900390556001600160a01b03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a350600192915050565b6001600160a01b0391821660009081526002602090815260408083209390941682529190915220549056fea2646970667358221220d735179f5013280e01a8920cdb567618f0c4f908512c41b2ad3d140c9468150864736f6c634300060c0033",
}

// EIP20ABI is the input ABI used to generate the binding from.
// Deprecated: Use EIP20MetaData.ABI instead.
var EIP20ABI = EIP20MetaData.ABI

// EIP20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EIP20MetaData.Bin instead.
var EIP20Bin = EIP20MetaData.Bin

// DeployEIP20 deploys a new Ethereum contract, binding an instance of EIP20 to it.
func DeployEIP20(auth *bind.TransactOpts, backend bind.ContractBackend, _initialAmount *big.Int, _tokenName string, _decimalUnits uint8, _tokenSymbol string) (common.Address, *types.Transaction, *EIP20, error) {
	parsed, err := EIP20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EIP20Bin), backend, _initialAmount, _tokenName, _decimalUnits, _tokenSymbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EIP20{EIP20Caller: EIP20Caller{contract: contract}, EIP20Transactor: EIP20Transactor{contract: contract}, EIP20Filterer: EIP20Filterer{contract: contract}}, nil
}

// EIP20 is an auto generated Go binding around an Ethereum contract.
type EIP20 struct {
	EIP20Caller     // Read-only binding to the contract
	EIP20Transactor // Write-only binding to the contract
	EIP20Filterer   // Log filterer for contract events
}

// EIP20Caller is an auto generated read-only Go binding around an Ethereum contract.
type EIP20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP20Session struct {
	Contract     *EIP20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP20CallerSession struct {
	Contract *EIP20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EIP20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP20TransactorSession struct {
	Contract     *EIP20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP20Raw is an auto generated low-level Go binding around an Ethereum contract.
type EIP20Raw struct {
	Contract *EIP20 // Generic contract binding to access the raw methods on
}

// EIP20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP20CallerRaw struct {
	Contract *EIP20Caller // Generic read-only contract binding to access the raw methods on
}

// EIP20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP20TransactorRaw struct {
	Contract *EIP20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP20 creates a new instance of EIP20, bound to a specific deployed contract.
func NewEIP20(address common.Address, backend bind.ContractBackend) (*EIP20, error) {
	contract, err := bindEIP20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP20{EIP20Caller: EIP20Caller{contract: contract}, EIP20Transactor: EIP20Transactor{contract: contract}, EIP20Filterer: EIP20Filterer{contract: contract}}, nil
}

// NewEIP20Caller creates a new read-only instance of EIP20, bound to a specific deployed contract.
func NewEIP20Caller(address common.Address, caller bind.ContractCaller) (*EIP20Caller, error) {
	contract, err := bindEIP20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20Caller{contract: contract}, nil
}

// NewEIP20Transactor creates a new write-only instance of EIP20, bound to a specific deployed contract.
func NewEIP20Transactor(address common.Address, transactor bind.ContractTransactor) (*EIP20Transactor, error) {
	contract, err := bindEIP20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP20Transactor{contract: contract}, nil
}

// NewEIP20Filterer creates a new log filterer instance of EIP20, bound to a specific deployed contract.
func NewEIP20Filterer(address common.Address, filterer bind.ContractFilterer) (*EIP20Filterer, error) {
	contract, err := bindEIP20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP20Filterer{contract: contract}, nil
}

// bindEIP20 binds a generic wrapper to an already deployed contract.
func bindEIP20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EIP20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20 *EIP20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP20.Contract.EIP20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20 *EIP20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20.Contract.EIP20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20 *EIP20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20.Contract.EIP20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP20 *EIP20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP20 *EIP20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP20 *EIP20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256 remaining)
func (_EIP20 *EIP20Caller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256 remaining)
func (_EIP20 *EIP20Session) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowance(&_EIP20.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256 remaining)
func (_EIP20 *EIP20CallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowance(&_EIP20.CallOpts, _owner, _spender)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) view returns(uint256)
func (_EIP20 *EIP20Caller) Allowed(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "allowed", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) view returns(uint256)
func (_EIP20 *EIP20Session) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowed(&_EIP20.CallOpts, arg0, arg1)
}

// Allowed is a free data retrieval call binding the contract method 0x5c658165.
//
// Solidity: function allowed(address , address ) view returns(uint256)
func (_EIP20 *EIP20CallerSession) Allowed(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Allowed(&_EIP20.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_EIP20 *EIP20Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_EIP20 *EIP20Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20.Contract.BalanceOf(&_EIP20.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_EIP20 *EIP20CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _EIP20.Contract.BalanceOf(&_EIP20.CallOpts, _owner)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_EIP20 *EIP20Caller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_EIP20 *EIP20Session) Balances(arg0 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Balances(&_EIP20.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_EIP20 *EIP20CallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _EIP20.Contract.Balances(&_EIP20.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EIP20 *EIP20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EIP20 *EIP20Session) Decimals() (uint8, error) {
	return _EIP20.Contract.Decimals(&_EIP20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_EIP20 *EIP20CallerSession) Decimals() (uint8, error) {
	return _EIP20.Contract.Decimals(&_EIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EIP20 *EIP20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EIP20 *EIP20Session) Name() (string, error) {
	return _EIP20.Contract.Name(&_EIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_EIP20 *EIP20CallerSession) Name() (string, error) {
	return _EIP20.Contract.Name(&_EIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EIP20 *EIP20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EIP20 *EIP20Session) Symbol() (string, error) {
	return _EIP20.Contract.Symbol(&_EIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_EIP20 *EIP20CallerSession) Symbol() (string, error) {
	return _EIP20.Contract.Symbol(&_EIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EIP20 *EIP20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EIP20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EIP20 *EIP20Session) TotalSupply() (*big.Int, error) {
	return _EIP20.Contract.TotalSupply(&_EIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_EIP20 *EIP20CallerSession) TotalSupply() (*big.Int, error) {
	return _EIP20.Contract.TotalSupply(&_EIP20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Transactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Session) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Approve(&_EIP20.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool success)
func (_EIP20 *EIP20TransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Approve(&_EIP20.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Session) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Transfer(&_EIP20.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20TransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.Transfer(&_EIP20.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20Session) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.TransferFrom(&_EIP20.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool success)
func (_EIP20 *EIP20TransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _EIP20.Contract.TransferFrom(&_EIP20.TransactOpts, _from, _to, _value)
}

// EIP20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the EIP20 contract.
type EIP20ApprovalIterator struct {
	Event *EIP20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EIP20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EIP20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EIP20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20Approval represents a Approval event raised by the EIP20 contract.
type EIP20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_EIP20 *EIP20Filterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*EIP20ApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &EIP20ApprovalIterator{contract: _EIP20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_EIP20 *EIP20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *EIP20Approval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _EIP20.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20Approval)
				if err := _EIP20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_EIP20 *EIP20Filterer) ParseApproval(log types.Log) (*EIP20Approval, error) {
	event := new(EIP20Approval)
	if err := _EIP20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EIP20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the EIP20 contract.
type EIP20TransferIterator struct {
	Event *EIP20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EIP20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EIP20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EIP20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EIP20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EIP20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EIP20Transfer represents a Transfer event raised by the EIP20 contract.
type EIP20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_EIP20 *EIP20Filterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*EIP20TransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &EIP20TransferIterator{contract: _EIP20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_EIP20 *EIP20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EIP20Transfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _EIP20.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EIP20Transfer)
				if err := _EIP20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_EIP20 *EIP20Filterer) ParseTransfer(log types.Log) (*EIP20Transfer, error) {
	event := new(EIP20Transfer)
	if err := _EIP20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
