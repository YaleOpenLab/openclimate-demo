package blockchain

import (
	"fmt"
	web3go "github.com/bcl-chain/web3.go/mobile"
	"log"
)

//  YOCL (GRC20) token is deployed on Kovan testnet. Initial supply is 1000000000.
// https://kovan.etherscan.io/tx/0x998aadf4602f79b8f8e4e84def8b84191237d3b6abe73e715546303c072259b5
// To interact with the token we are using Web3go library. Connection through infura endpoint.

const (
	contractAddress = "0xb19ac159b87a4491b3b8bef4554b59da2bf42555"
	ownerAddress    = "0xfE1827f2F1C366c04d458b3c07B8Bd207D42eab4"
	rpcUrl          = "https://kovan.infura.io/v3/def7370cf49d49d791b9df949986b9a0"
)

func CheckTokenBalance() {
	//New Ethereum Client
	client, err := web3go.NewEthereumClient(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	address, _ := web3go.NewAddressFromHex(ownerAddress)
	// 3. get ether balance from the latest block
	balance, err := client.GetBalanceAt(web3go.NewContext(), address, -1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	//Access the Token
	contractAddress, _ := web3go.NewAddressFromHex(contractAddress)
	YToken, err := web3go.NewERC20(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	/*
		//Check if the address is a Token owner
		isMinter, _ := YToken.IsMinter(address)
		if isMinter {
			fmt.Println("True")
		}else{
			fmt.Println("False")
		}
	*/
	tokenBalance, err := YToken.BalanceOf(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("YOCL Token Balance", tokenBalance.String())
}
