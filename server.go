package main
import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"html/template"
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
func build_javascript_data_obj(mdepth *MarketDepth) (template.JS,template.JS) {
	var asks_str string = "["
	var bids_str string = "["

	var last_price float64 = 0.0
	for i:=0 ; i < len(mdepth.Asks) ; i++ {
		if len(asks_str) > 1 {
			asks_str = asks_str + ","
		}
		var entry string
		entry = "{" +
				"x:" + fmt.Sprintf("%v",mdepth.Asks[i].Price*10)  + "," +
				"y:"  + fmt.Sprintf("%v",mdepth.Asks[i].Price) + "," +
				"price: " + fmt.Sprintf("%v",mdepth.Asks[i].Price) + "," +
				"addr: \"" + mdepth.Asks[i].EOAAddrSh + "\"," +
				"expires: \"" + mdepth.Asks[i].Expires + "\"," +
				"volume: " + fmt.Sprintf("%v",mdepth.Asks[i].Volume) + "," +
				"click: function() {load_order_data(\"" +
					mdepth.Asks[i].EOAAddrSh +"\",\"" +
					mdepth.Asks[i].WalletAddrSh + "\"," +
					fmt.Sprintf("%v,%v,%v,%v,%v,\"%v\",\"%v\"",
										mdepth.Asks[i].MktAid,
										mdepth.Asks[i].OutcomeIdx,
										mdepth.Asks[i].TotalBids,
										mdepth.Asks[i].TotalAsks,
										mdepth.Asks[i].TotalCancel,
										mdepth.Asks[i].DateCreated,
										mdepth.Asks[i].Expires,
					) +
				")}," +
				"toolTipContent: \"<div>User: {addr}<br/>Expires: {expires}<br/>ASK: {price}<br/>Volume: {volume}</div>\"" +
				"}"
		asks_str= asks_str + entry
		last_price = mdepth.Asks[i].Price
	}
// duplicate of above loop begins (ToDo: generalize it
	for i:=0 ; i < len(mdepth.Bids) ; i++ {
		if len(bids_str) > 1 {
			bids_str = bids_str + ","
		}
		var entry string
		entry = "{" +
				"x:" + fmt.Sprintf("%v",mdepth.Bids[i].Price*10)  + "," +
				"y:"  + fmt.Sprintf("%v",mdepth.Bids[i].Price) + "," +
				"price: " + fmt.Sprintf("%v",mdepth.Bids[i].Price) + "," +
				"addr: \"" + mdepth.Bids[i].EOAAddrSh + "\"," +
				"expires: \"" + mdepth.Bids[i].Expires + "\"," +
				"volume: " + fmt.Sprintf("%v",mdepth.Bids[i].Volume) + "," +
				"click: function() {load_order_data(\"" +
					mdepth.Bids[i].EOAAddrSh +"\",\"" +
					mdepth.Bids[i].WalletAddrSh + "\"," +
					fmt.Sprintf("%v,%v,%v,%v,%v,\"%v\",\"%v\"",
										mdepth.Bids[i].MktAid,
										mdepth.Bids[i].OutcomeIdx,
										mdepth.Bids[i].TotalBids,
										mdepth.Bids[i].TotalAsks,
										mdepth.Bids[i].TotalCancel,
										mdepth.Bids[i].DateCreated,
										mdepth.Bids[i].Expires,
					) +
				")}," +
				"toolTipContent: \"<div>User: {addr}<br/>Expires: {expires}<br/>BID: {price}<br/>Volume: {volume}</div>\"" +
				"}"
		bids_str= bids_str + entry
		last_price = mdepth.Bids[i].Price
	}
// duplicate of above loop ends


	asks_str = asks_str + "]"
	bids_str = bids_str + "]"
	_ = last_price
	return template.JS(bids_str),template.JS(asks_str)
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
	blknum,res := augur_srv.storage.get_last_block_num()
	_ = res
	c.HTML(http.StatusOK, "statistics.html", gin.H{
			"title": "Augur Market Statistics",
			"block_num" : blknum,
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
	market_info,_ := augur_srv.storage.get_market_info(market)
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
	market := c.Param("market")
	p_outcome := c.Param("outcome")
	var outcome uint8
	if len(p_outcome) > 0 {
		p, err := strconv.Atoi(p_outcome)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"title": "Augur Markets: Error",
				"ErrDescr": "Can't parse offset",
			})
			return
		}
		outcome = uint8(p)
	} else {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"title": "Augur Markets: Error",
			"ErrDescr": "Can't parse offset",
		})
		return
	}

	market_info,_ := augur_srv.storage.get_market_info(market)
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
