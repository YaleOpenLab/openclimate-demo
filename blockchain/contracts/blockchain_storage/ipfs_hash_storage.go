// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IpfsRootABI is the input ABI used to generate the binding from.
const IpfsRootABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"timeStamp\",\"type\":\"uint256\"},{\"name\":\"rootHash\",\"type\":\"bytes32\"}],\"name\":\"insertRoot\",\"outputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"timeStamp\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"rootHash\",\"type\":\"bytes32\"}],\"name\":\"LogNewRoot\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"name\":\"timeStamp\",\"type\":\"uint256\"}],\"name\":\"checkTimeStamp\",\"outputs\":[{\"name\":\"isIndeed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"timeStamp\",\"type\":\"uint256\"}],\"name\":\"getRoot\",\"outputs\":[{\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRootAtIndex\",\"outputs\":[{\"name\":\"timeStamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRootCount\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// IpfsRoot is an auto generated Go binding around an Ethereum contract.
type IpfsRoot struct {
	IpfsRootCaller     // Read-only binding to the contract
	IpfsRootTransactor // Write-only binding to the contract
	IpfsRootFilterer   // Log filterer for contract events
}

// IpfsRootCaller is an auto generated read-only Go binding around an Ethereum contract.
type IpfsRootCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IpfsRootTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IpfsRootTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IpfsRootFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IpfsRootFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IpfsRootSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IpfsRootSession struct {
	Contract     *IpfsRoot         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IpfsRootCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IpfsRootCallerSession struct {
	Contract *IpfsRootCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IpfsRootTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IpfsRootTransactorSession struct {
	Contract     *IpfsRootTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IpfsRootRaw is an auto generated low-level Go binding around an Ethereum contract.
type IpfsRootRaw struct {
	Contract *IpfsRoot // Generic contract binding to access the raw methods on
}

// IpfsRootCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IpfsRootCallerRaw struct {
	Contract *IpfsRootCaller // Generic read-only contract binding to access the raw methods on
}

// IpfsRootTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IpfsRootTransactorRaw struct {
	Contract *IpfsRootTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIpfsRoot creates a new instance of IpfsRoot, bound to a specific deployed contract.
func NewIpfsRoot(address common.Address, backend bind.ContractBackend) (*IpfsRoot, error) {
	contract, err := bindIpfsRoot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IpfsRoot{IpfsRootCaller: IpfsRootCaller{contract: contract}, IpfsRootTransactor: IpfsRootTransactor{contract: contract}, IpfsRootFilterer: IpfsRootFilterer{contract: contract}}, nil
}

// NewIpfsRootCaller creates a new read-only instance of IpfsRoot, bound to a specific deployed contract.
func NewIpfsRootCaller(address common.Address, caller bind.ContractCaller) (*IpfsRootCaller, error) {
	contract, err := bindIpfsRoot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IpfsRootCaller{contract: contract}, nil
}

// NewIpfsRootTransactor creates a new write-only instance of IpfsRoot, bound to a specific deployed contract.
func NewIpfsRootTransactor(address common.Address, transactor bind.ContractTransactor) (*IpfsRootTransactor, error) {
	contract, err := bindIpfsRoot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IpfsRootTransactor{contract: contract}, nil
}

// NewIpfsRootFilterer creates a new log filterer instance of IpfsRoot, bound to a specific deployed contract.
func NewIpfsRootFilterer(address common.Address, filterer bind.ContractFilterer) (*IpfsRootFilterer, error) {
	contract, err := bindIpfsRoot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IpfsRootFilterer{contract: contract}, nil
}

// bindIpfsRoot binds a generic wrapper to an already deployed contract.
func bindIpfsRoot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IpfsRootABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IpfsRoot *IpfsRootRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IpfsRoot.Contract.IpfsRootCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IpfsRoot *IpfsRootRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IpfsRoot.Contract.IpfsRootTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IpfsRoot *IpfsRootRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IpfsRoot.Contract.IpfsRootTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IpfsRoot *IpfsRootCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IpfsRoot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IpfsRoot *IpfsRootTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IpfsRoot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IpfsRoot *IpfsRootTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IpfsRoot.Contract.contract.Transact(opts, method, params...)
}

// CheckTimeStamp is a free data retrieval call binding the contract method 0xa89e7246.
//
// Solidity: function checkTimeStamp(uint256 timeStamp) constant returns(bool isIndeed)
func (_IpfsRoot *IpfsRootCaller) CheckTimeStamp(opts *bind.CallOpts, timeStamp *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _IpfsRoot.contract.Call(opts, out, "checkTimeStamp", timeStamp)
	return *ret0, err
}

// CheckTimeStamp is a free data retrieval call binding the contract method 0xa89e7246.
//
// Solidity: function checkTimeStamp(uint256 timeStamp) constant returns(bool isIndeed)
func (_IpfsRoot *IpfsRootSession) CheckTimeStamp(timeStamp *big.Int) (bool, error) {
	return _IpfsRoot.Contract.CheckTimeStamp(&_IpfsRoot.CallOpts, timeStamp)
}

// CheckTimeStamp is a free data retrieval call binding the contract method 0xa89e7246.
//
// Solidity: function checkTimeStamp(uint256 timeStamp) constant returns(bool isIndeed)
func (_IpfsRoot *IpfsRootCallerSession) CheckTimeStamp(timeStamp *big.Int) (bool, error) {
	return _IpfsRoot.Contract.CheckTimeStamp(&_IpfsRoot.CallOpts, timeStamp)
}

// GetRoot is a free data retrieval call binding the contract method 0x9b24b3b0.
//
// Solidity: function getRoot(uint256 timeStamp) constant returns(bytes32 rootHash, uint256 index)
func (_IpfsRoot *IpfsRootCaller) GetRoot(opts *bind.CallOpts, timeStamp *big.Int) (struct {
	RootHash [32]byte
	Index    *big.Int
}, error) {
	ret := new(struct {
		RootHash [32]byte
		Index    *big.Int
	})
	out := ret
	err := _IpfsRoot.contract.Call(opts, out, "getRoot", timeStamp)
	return *ret, err
}

// GetRoot is a free data retrieval call binding the contract method 0x9b24b3b0.
//
// Solidity: function getRoot(uint256 timeStamp) constant returns(bytes32 rootHash, uint256 index)
func (_IpfsRoot *IpfsRootSession) GetRoot(timeStamp *big.Int) (struct {
	RootHash [32]byte
	Index    *big.Int
}, error) {
	return _IpfsRoot.Contract.GetRoot(&_IpfsRoot.CallOpts, timeStamp)
}

// GetRoot is a free data retrieval call binding the contract method 0x9b24b3b0.
//
// Solidity: function getRoot(uint256 timeStamp) constant returns(bytes32 rootHash, uint256 index)
func (_IpfsRoot *IpfsRootCallerSession) GetRoot(timeStamp *big.Int) (struct {
	RootHash [32]byte
	Index    *big.Int
}, error) {
	return _IpfsRoot.Contract.GetRoot(&_IpfsRoot.CallOpts, timeStamp)
}

// GetRootAtIndex is a free data retrieval call binding the contract method 0x6b5f0373.
//
// Solidity: function getRootAtIndex(uint256 index) constant returns(uint256 timeStamp)
func (_IpfsRoot *IpfsRootCaller) GetRootAtIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IpfsRoot.contract.Call(opts, out, "getRootAtIndex", index)
	return *ret0, err
}

// GetRootAtIndex is a free data retrieval call binding the contract method 0x6b5f0373.
//
// Solidity: function getRootAtIndex(uint256 index) constant returns(uint256 timeStamp)
func (_IpfsRoot *IpfsRootSession) GetRootAtIndex(index *big.Int) (*big.Int, error) {
	return _IpfsRoot.Contract.GetRootAtIndex(&_IpfsRoot.CallOpts, index)
}

// GetRootAtIndex is a free data retrieval call binding the contract method 0x6b5f0373.
//
// Solidity: function getRootAtIndex(uint256 index) constant returns(uint256 timeStamp)
func (_IpfsRoot *IpfsRootCallerSession) GetRootAtIndex(index *big.Int) (*big.Int, error) {
	return _IpfsRoot.Contract.GetRootAtIndex(&_IpfsRoot.CallOpts, index)
}

// GetRootCount is a free data retrieval call binding the contract method 0xfa12779c.
//
// Solidity: function getRootCount() constant returns(uint256 count)
func (_IpfsRoot *IpfsRootCaller) GetRootCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _IpfsRoot.contract.Call(opts, out, "getRootCount")
	return *ret0, err
}

// GetRootCount is a free data retrieval call binding the contract method 0xfa12779c.
//
// Solidity: function getRootCount() constant returns(uint256 count)
func (_IpfsRoot *IpfsRootSession) GetRootCount() (*big.Int, error) {
	return _IpfsRoot.Contract.GetRootCount(&_IpfsRoot.CallOpts)
}

// GetRootCount is a free data retrieval call binding the contract method 0xfa12779c.
//
// Solidity: function getRootCount() constant returns(uint256 count)
func (_IpfsRoot *IpfsRootCallerSession) GetRootCount() (*big.Int, error) {
	return _IpfsRoot.Contract.GetRootCount(&_IpfsRoot.CallOpts)
}

// InsertRoot is a paid mutator transaction binding the contract method 0xf3f15c29.
//
// Solidity: function insertRoot(uint256 timeStamp, bytes32 rootHash) returns(uint256 index)
func (_IpfsRoot *IpfsRootTransactor) InsertRoot(opts *bind.TransactOpts, timeStamp *big.Int, rootHash [32]byte) (*types.Transaction, error) {
	return _IpfsRoot.contract.Transact(opts, "insertRoot", timeStamp, rootHash)
}

// InsertRoot is a paid mutator transaction binding the contract method 0xf3f15c29.
//
// Solidity: function insertRoot(uint256 timeStamp, bytes32 rootHash) returns(uint256 index)
func (_IpfsRoot *IpfsRootSession) InsertRoot(timeStamp *big.Int, rootHash [32]byte) (*types.Transaction, error) {
	return _IpfsRoot.Contract.InsertRoot(&_IpfsRoot.TransactOpts, timeStamp, rootHash)
}

// InsertRoot is a paid mutator transaction binding the contract method 0xf3f15c29.
//
// Solidity: function insertRoot(uint256 timeStamp, bytes32 rootHash) returns(uint256 index)
func (_IpfsRoot *IpfsRootTransactorSession) InsertRoot(timeStamp *big.Int, rootHash [32]byte) (*types.Transaction, error) {
	return _IpfsRoot.Contract.InsertRoot(&_IpfsRoot.TransactOpts, timeStamp, rootHash)
}

// IpfsRootLogNewRootIterator is returned from FilterLogNewRoot and is used to iterate over the raw logs and unpacked data for LogNewRoot events raised by the IpfsRoot contract.
type IpfsRootLogNewRootIterator struct {
	Event *IpfsRootLogNewRoot // Event containing the contract specifics and raw log

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
func (it *IpfsRootLogNewRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IpfsRootLogNewRoot)
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
		it.Event = new(IpfsRootLogNewRoot)
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
func (it *IpfsRootLogNewRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IpfsRootLogNewRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IpfsRootLogNewRoot represents a LogNewRoot event raised by the IpfsRoot contract.
type IpfsRootLogNewRoot struct {
	TimeStamp *big.Int
	Index     *big.Int
	RootHash  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogNewRoot is a free log retrieval operation binding the contract event 0xb843516edf882622b5cf95b5485456d18fb888581f9a154cfe021aeed228ee86.
//
// Solidity: event LogNewRoot(uint256 indexed timeStamp, uint256 index, bytes32 rootHash)
func (_IpfsRoot *IpfsRootFilterer) FilterLogNewRoot(opts *bind.FilterOpts, timeStamp []*big.Int) (*IpfsRootLogNewRootIterator, error) {

	var timeStampRule []interface{}
	for _, timeStampItem := range timeStamp {
		timeStampRule = append(timeStampRule, timeStampItem)
	}

	logs, sub, err := _IpfsRoot.contract.FilterLogs(opts, "LogNewRoot", timeStampRule)
	if err != nil {
		return nil, err
	}
	return &IpfsRootLogNewRootIterator{contract: _IpfsRoot.contract, event: "LogNewRoot", logs: logs, sub: sub}, nil
}

// WatchLogNewRoot is a free log subscription operation binding the contract event 0xb843516edf882622b5cf95b5485456d18fb888581f9a154cfe021aeed228ee86.
//
// Solidity: event LogNewRoot(uint256 indexed timeStamp, uint256 index, bytes32 rootHash)
func (_IpfsRoot *IpfsRootFilterer) WatchLogNewRoot(opts *bind.WatchOpts, sink chan<- *IpfsRootLogNewRoot, timeStamp []*big.Int) (event.Subscription, error) {

	var timeStampRule []interface{}
	for _, timeStampItem := range timeStamp {
		timeStampRule = append(timeStampRule, timeStampItem)
	}

	logs, sub, err := _IpfsRoot.contract.WatchLogs(opts, "LogNewRoot", timeStampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IpfsRootLogNewRoot)
				if err := _IpfsRoot.contract.UnpackLog(event, "LogNewRoot", log); err != nil {
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

// ParseLogNewRoot is a log parse operation binding the contract event 0xb843516edf882622b5cf95b5485456d18fb888581f9a154cfe021aeed228ee86.
//
// Solidity: event LogNewRoot(uint256 indexed timeStamp, uint256 index, bytes32 rootHash)
func (_IpfsRoot *IpfsRootFilterer) ParseLogNewRoot(log types.Log) (*IpfsRootLogNewRoot, error) {
	event := new(IpfsRootLogNewRoot)
	if err := _IpfsRoot.contract.UnpackLog(event, "LogNewRoot", log); err != nil {
		return nil, err
	}
	return event, nil
}
