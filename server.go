package main
import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
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
func build_javascript_data_obj(mdepth *MarketDepth) string {
	var output string

	var last_price float64 = 0.0
	for i:=0 ; i< mdepth.Asks {
		if len(output) > 0 {
			output = output + ","
		}
		output = output + mdepth.Asks[i].Price
		last_price = mdepth.Asks[i]
	}
	return output
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
	js_data := build_javascript_data_obj(&mdepth)
	c.HTML(http.StatusOK, "market_depth.html", gin.H{
			"title": "Market Depth",
			"Market": market_info,
			"Bids": mdepth.Bids,
			"Asks": mdepth.Asks,
			"JSData": js_data,
	})
}
