// test.go
package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type AssetsmanangementChaincode struct{}

func (t *AssetsmanangementChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 4 {
		return shim.Error("please specify the correct account and balence")
	}
	var a = args[0]
	var abalence = args[1]
	var b = args[2]
	var bbalence = args[3]
	//检验输入的金额是否正确
	_, err := strconv.Atoi(abalence)
	if err != nil {
		return shim.Error("please specify the correct balence of a ")
	}
	_, err = strconv.Atoi(bbalence)
	if err != nil {
		return shim.Error("please specify the correct balence of b ")
	}
	//put the acount state of a b  into the ledger
	err = stub.PutState(a, []byte(abalence)) //将键值对写入账本
	if err != nil {
		return shim.Error("erro when put the a acount into ledger ")
	}
	err = stub.PutState(b, []byte(bbalence)) //将键值对写入账本
	if err != nil {
		return shim.Error("erro when put the b acount into ledger ")
	}
	return shim.Success([]byte("sucessful init"))
}

func (t *AssetsmanangementChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, args := stub.GetFunctionAndParameters()
	if fun == "set" { //存钱
		return t.set(stub, args)
	} else if fun == "get" { //取钱
		return t.get(stub, args)
	} else if fun == "transfer" {
		return transfer(stub, args)
	} else if fun == "query" {
		return query(stub, args)
	}
	return shim.Error("wrong operation")

}

func (t *AssetsmanangementChaincode) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("please specify the correct account and balence")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("information query fail")
	}
	if result == nil {
		return shim.Error("no this account")
	}
	balence, err := strconv.Atoi(string(result))
	if err != nil {
		return shim.Error("erro when dealwith the balence")
	}
	x, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("please enter the correct money")
	}
	balence = balence + x
	err = stub.PutState(args[0], []byte(strconv.Itoa(balence))) //将键值对写入账本
	if err != nil {
		return shim.Error("erro when put the a acount into ledger ")
	}
	return shim.Success([]byte("sucessful save money "))
}

func (t *AssetsmanangementChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("please enter the correct account and balence ")
	}
	result, err := stub.GetState(args[0]) //从账本中获取args[0]对应的value
	if err != nil {
		return shim.Error("information query fail")
	}
	if result == nil {
		return shim.Error("no this account")
	}
	balence, err := strconv.Atoi(string(result))
	if err != nil {
		return shim.Error("erro when dealwith the balence")
	}
	x, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("please enter the correct money")
	}
	if x > balence {
		return shim.Error("no enough balence")
	}
	balence = balence - x
	err = stub.PutState(args[0], []byte(strconv.Itoa(balence))) //将键值对写入账本
	if err != nil {
		return shim.Error("erro when put the a acount into ledger ")
	}
	return shim.Success([]byte("sucessful get money "))
}

func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("请输入正确的账户")
	}
	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("获取账户失败")
	}
	return shim.Success(result)
}

func transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("please enter correct account and amount")
	}
	source := args[0]
	target := args[1]
	amount := args[2]

	source1, err := stub.GetState(source) //获取原账户的余额(由key获取value)
	if err != nil {
		return shim.Error("获取原账户信息失败")
	}
	s1, err := strconv.Atoi(string(source1)) //将余额抓换为数字
	if err != nil {
		return shim.Error("处理失败")
	}

	target1, err := stub.GetState(target)
	if err != nil {
		return shim.Error("获取mubiao账户信息失败")
	}
	t1, err := strconv.Atoi(string(target1))
	if err != nil {
		return shim.Error("处理失败")
	}

	a, err := strconv.Atoi(amount)
	if err != nil {
		return shim.Error("处理失败")
	}

	if s1 < a {
		return shim.Error("no enough money on source account")
	}
	s1 = s1 - a
	t1 = t1 + a

	err = stub.PutState(source, []byte(strconv.Itoa(s1)))
	if err != nil {
		return shim.Error("保存处理后源数据失败")
	}

	err = stub.PutState(target, []byte(strconv.Itoa(t1)))
	if err != nil {
		return shim.Error("保存处理后目标数据失败")
	}
	return shim.Success([]byte("sucessful transfer"))
}
func main() {
	err := shim.Start(new(AssetsmanangementChaincode))
	if err != nil {
		fmt.Println("start error")
	}
}
