package main

import (
	"log"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/accounts/abi"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
)
var (
	RPC_URL = os.Getenv("AUGUR_ETH_NODE_RPC_URL")
	rpcclient *ethclient.Client

	augur_srv *AugurServer

	// contracts
	all_contracts map[string]interface{}

	dai_addr common.Address
	rep_addr common.Address

	ctrct_dai_token *DAICash
	ctrct_rep_token *RepTok

	// thes variables should be removed on the next code reorg task
	market_order_id int64 = 0
	fill_order_id int64 = 0
)
func initialize() {

	contract_addresses := new(ContractAddresses)
	contract_addresses.dai_addr = &dai_addr
	contract_addresses.reputation_addr = &rep_addr
	init_contract_addresses(contract_addresses)
}
func main() {

	initialize()

	if len(RPC_URL) == 0 {
		fmt.Printf("Configuration error: RPC URL of Ethereum node is not set."+
			"Calls to contracts are disabled. " +
			" Please set AUGUR_ETH_NODE_RPC environment variable")
	} else {
		var err error
		rpcclient, err = ethclient.Dial(RPC_URL)
		if err != nil {
			log.Fatal(err)
		}
		// init contracts
		fmt.Println("init DAI contract with addr %v\n",dai_addr.String())
		ctrct_dai_token,err = NewDAICash(dai_addr,rpcclient)
		if err != nil {
			Fatalf("Couldn't initialize DAI Cash contract: %v\n",err)
		}

		fmt.Println("init REP contract with addr %v\n",rep_addr.String())
		ctrct_rep_token,err = NewRepTok(rep_addr,rpcclient)
		if err != nil {
			Fatalf("Couldn't initialize Rep Token contract: %v\n",err)
		}
	}

	port_plain := os.Getenv("AUGUR_HTTP_PORT")
	port_secure := os.Getenv("AUGUR_HTTPS_PORT")

	if port_plain == "" {
		port_plain = "9090"
		log.Printf("Defaulting plain HTTP to port %s", port_plain)
	}
	if port_secure== "" {
		port_secure= "9443"
		log.Printf("Defaulting secure protocol to port %s", port_secure)
	}

	augur_srv = create_augur_server()

	r := gin.New()
	r.LoadHTMLGlob("html/templates/*")
	r.Use(gin.Logger())

	// Main HTML
	r.GET("/", main_page)
	r.GET("/index.html", main_page)
	r.GET("/index.htm", main_page)
	r.GET("/INDEX.HTM", main_page)

	r.GET("/markets.html",markets)
	r.GET("/statistics.html",statistics)
	r.GET("/categories.html",categories)
	r.GET("/explorer.html",explorer)
	r.GET("/market/:market",  market_trades)
	r.GET("/mdepth/:market/:outcome", market_depth)
	r.GET("/mphist/:market/:outcome", market_price_history)
	r.GET("/search", search)
	r.GET("/money/:addr",  read_money)
	r.GET("/order/:order",  order)

	r.Static("/imgs", "./html/imgs")
	r.Static("/res", "./html/res")			// resources (static)
	r.StaticFile("/favicon.ico", "./html/res/favicon.ico")

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port_plain)
	r.Run(":" + port_plain)
}
