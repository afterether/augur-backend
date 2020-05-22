package main
import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"html/template"

	"github.com/ethereum/go-ethereum/common"
)
const (
	DEFAULT_MARKET_ROWS_LIMIT int	= 10
)
type AugurServer struct {
	storage		*SQLStorage
}
func create_augur_server() *AugurServer {

	srv := new(AugurServer)
	srv.storage = connect_to_storage()
	return srv
}
func mkt_depth_entry_to_js_obj(de *DepthEntry) string {

	var output string
	output = "{" +
				"x:" + fmt.Sprintf("%v",de.Price)  + "," +
				"y:"  + fmt.Sprintf("%v",de.AccumVol) + "," +
				"price: " + fmt.Sprintf("%v",de.Price) + "," +
				"addr: \"" + de.EOAAddrSh + "\"," +
				"expires: \"" + de.Expires + "\"," +
				"volume: " + fmt.Sprintf("%v",de.Volume) + "," +
				"click: function() {load_order_data(\"" +
					de.EOAAddrSh +"\",\"" +
					de.WalletAddrSh + "\"," +
					fmt.Sprintf("%v,%v,%v,%v,%v,%v,\"%v\",\"%v\"",
										de.MktAid,
										de.OutcomeIdx,
										de.Volume,
										de.TotalBids,
										de.TotalAsks,
										de.TotalCancel,
										de.DateCreated,
										de.Expires,
					) +
				")}" +
				"}"
	return output
}
func build_javascript_data_obj(mdepth *MarketDepth) (template.JS,template.JS) {
	var asks_str string = "["
	var bids_str string = "["

	var last_price float64 = 0.0
	for i:=0 ; i < len(mdepth.Asks) ; i++ {
		if len(asks_str) > 1 {
			asks_str = asks_str + ","
		}
		var entry string
		entry = mkt_depth_entry_to_js_obj(&mdepth.Asks[i])
		asks_str= asks_str + entry
		last_price = mdepth.Asks[i].Price
	}
/*	// Possibly replace this with a line indicating the spread, as another Serie
	if len(mdepth.Asks) > 0 {
		if len(mdepth.Bids) > 0 {
			// add fake BID entry to fill the hole for the spread
			last_elt:=len(mdepth.Asks)-1
			fake_entry := mdepth.Asks[last_elt]
			fake_entry.Price = mdepth.Bids[0].Price*10
			bids_str = "[" + mkt_depth_entry_to_js_obj(&fake_entry)
		}
	}
*/
	for i:=0 ; i < len(mdepth.Bids) ; i++ {
		if len(bids_str) > 1 {
			bids_str = bids_str + ","
		}
		var entry string
		entry = mkt_depth_entry_to_js_obj(&mdepth.Bids[i])
		bids_str= bids_str + entry
		last_price = mdepth.Bids[i].Price
	}

	asks_str = asks_str + "]"
	bids_str = bids_str + "]"
	_ = last_price
	fmt.Printf("asks JS string: %v\n",asks_str)
	fmt.Printf("bids JS string: %v\n",bids_str)
	return template.JS(bids_str),template.JS(asks_str)
}
func build_javascript_price_history(orders *[]MarketOrder) template.JS {
	var data_str string = "["

	for i:=0 ; i < len(*orders) ; i++ {
		if len(data_str) > 1 {
			data_str = data_str + ","
		}
		var e = &(*orders)[i];
		var entry string
		entry = "{" +
//				"x:" + fmt.Sprintf("\"%v\"",e.Date)  + "," +
				"x:" + fmt.Sprintf("%v",i)  + "," +
				"y:"  + fmt.Sprintf("%v",e.Price) + "," +
				"price: " + fmt.Sprintf("%v",e.Price) + "," +
				"volume: " + fmt.Sprintf("%v",e.Volume) + "," +
				"click: function() {load_order_data(\"" +
					e.SellerEOAAddrSh +"\",\"" +
					e.SellerWalletAddrSh + "\",\"" +
					e.BuyerEOAAddrSh + "\",\"" +
					e.BuyerWalletAddrSh + "\"," +
					fmt.Sprintf("%v,%v,%v,\"%v\"",e.MktAid,e.Price,e.Volume,e.Date) +
				")}" +
				"}"
		fmt.Printf("\nentry = %v\n",entry)
		data_str= data_str + entry
	}
	data_str = data_str + "]"
	fmt.Printf("JS price history string: %v\n",data_str)
	return template.JS(data_str)
}
func main_page(c *gin.Context) {
	blknum,_:= augur_srv.storage.get_last_block_num()
	c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Augur Prediction Markets",
			"block_num" : blknum,
	})
}
func markets(c *gin.Context) {
	off_str := c.Query("off")
	var off int = 0
	var err error
	fmt.Printf("off_str = %v\n",off_str)
	if len(off_str) > 0 {
		off, err = strconv.Atoi(off_str)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Can't parse offset",
			})
			return
		}
	}
	fmt.Printf("off = %v\n",off)
	markets := augur_srv.storage.get_active_market_list(off,DEFAULT_MARKET_ROWS_LIMIT)
	c.HTML(http.StatusOK, "markets.html", gin.H{
			"title": "Augur Markets",
			"Markets" : markets,
	})
}
func categories(c *gin.Context) {
	blknum,_ := augur_srv.storage.get_last_block_num()
	categories := augur_srv.storage.get_categories()
	c.HTML(http.StatusOK, "categories.html", gin.H{
			"title": "Augur Market Categories",
			"block_num" : blknum,
			"Categories" : categories,
	})
}
func statistics(c *gin.Context) {
	stats := augur_srv.storage.get_main_stats()
	c.HTML(http.StatusOK, "statistics.html", gin.H{
			"title": "Augur Market Statistics",
			"MainStats" : stats,
	})
}
func explorer(c *gin.Context) {
	blknum,res := augur_srv.storage.get_last_block_num()
	_ = res
	c.HTML(http.StatusOK, "explorer.html", gin.H{
			"title": "Augur Event Explorer",
			"block_num" : blknum,
	})
}
func market_trades(c *gin.Context) {
	market := c.Param("market")
	fmt.Printf("getting trades for market %v\n",market)
	market_info,_ := augur_srv.storage.get_market_info(market,0,false)
	fmt.Printf("market info = %+v",market_info)
	trades := augur_srv.storage.get_mkt_trades(market)
	outcome_vols,_ := augur_srv.storage.get_outcome_volumes(market)
	c.HTML(http.StatusOK, "market_info.html", gin.H{
			"title": "Trades for market",
			"Trades" : trades,
			"Market": market_info,
			"OutcomeVols" : outcome_vols,
	})
}
func market_depth(c *gin.Context) {

	// Market Depth Info: https://en.wikipedia.org/wiki/Order_book_(trading)
	market := c.Param("market")
	p_outcome := c.Param("outcome")
	var outcome int
	if len(p_outcome) > 0 {
		p, err := strconv.Atoi(p_outcome)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Can't parse outcome",
			})
			return
		}
		outcome = int(p)
	} else {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Augur Markets: Error",
			"ErrDescr": "Can't parse outcome",
		})
		return
	}

	market_info,err := augur_srv.storage.get_market_info(market,outcome,true)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Augur Markets: Error",
			"ErrDescr": "Can't find this market, address not registered",
		})
		return
	}
	mdepth := augur_srv.storage.get_mkt_depth(market,outcome)
	js_bid_data,js_ask_data := build_javascript_data_obj(mdepth)
	fmt.Printf("js ask_data = %v\n",js_ask_data)
	fmt.Printf("js bid_data = %v\n",js_bid_data)
	c.HTML(http.StatusOK, "market_depth.html", gin.H{
			"title": "Market Depth",
			"Market": market_info,
			"Bids": mdepth.Bids,
			"Asks": mdepth.Asks,
			"JSAskData": js_ask_data,
			"JSBidData": js_bid_data,
	})
}
func market_price_history(c *gin.Context) {

	market := c.Param("market")
	p_outcome := c.Param("outcome")
	var outcome int
	if len(p_outcome) > 0 {
		p, err := strconv.Atoi(p_outcome)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Can't parse outcome",
			})
			return
		}
		outcome = int(p)
	} else {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Augur Markets: Error",
			"ErrDescr": "Can't parse outcome",
		})
		return
	}

	market_info,err := augur_srv.storage.get_market_info(market,outcome,true)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Augur Markets: Error",
			"ErrDescr": "Can't find this market, address not registered",
		})
		return
	}
	mkt_price_hist := augur_srv.storage.get_price_history_for_outcome(market_info.MktAid,outcome)
	js_price_history := build_javascript_price_history(&mkt_price_hist)
	fmt.Printf("js price history = %v\n",js_price_history)
	c.HTML(http.StatusOK, "price_history.html", gin.H{
			"title": "Market Price History",
			"Market": market_info,
			"Prices": mkt_price_hist,
			"JSPriceData": js_price_history,
	})
}
func serve_user_info_page(c *gin.Context,addr string) {

	aid,err := augur_srv.storage.nonfatal_lookup_address(addr)
	if err == nil {
		user_info := augur_srv.storage.get_user_info(aid)
		c.HTML(http.StatusOK, "user_info.html", gin.H{
			"title": "User "+addr,
			"user_addr": addr,
			"UserInfo" : user_info,
		})
	} else {
		c.HTML(http.StatusOK, "user_not_found.html", gin.H{
			"title": "User "+addr,
			"user_addr": addr,
		})
	}
}
func serve_tx_info_page(c *gin.Context,tx_hash string) {

	c.HTML(http.StatusOK, "tx_info.html", gin.H{
			"title": "Transaction " + tx_hash,
			"tx_hash" : tx_hash,
	})
}
func search(c *gin.Context) {

	keyword := c.Query("q")
	fmt.Printf("Searching for %v\n",keyword)
	if (len(keyword) == 40) || (len(keyword) == 42) { // address
		if len(keyword) == 42 {	// strip 0x prefix
			keyword = keyword[2:]
		}
		addr_bytes,err := hex.DecodeString(keyword)
		if err == nil {
			addr := common.BytesToAddress(addr_bytes)
			addr_str := addr.String()
			serve_user_info_page(c,addr_str)
		} else {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Invalid HEX string in address parameter",
			})
			return
		}
	}
	if (len(keyword) == 64) || (len(keyword) == 66) { // Hash (Tx hash)
		if len(keyword) == 66 {	// strip 0x prefix
			keyword = keyword[2:]
		}
		hash_bytes,err := hex.DecodeString(keyword)
		if err != nil {
			hash := common.BytesToHash(hash_bytes)
			hash_str := hash.String()
			serve_tx_info_page(c,hash_str)
		} else {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Invalid HEX string in hash parameter",
			})
			return
		}
	}
}
