package main
import (
	//"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)
type AugurServer struct {
	storage		*SQLStorage
}
func create_augur_server() *AugurServer {

	srv := new(AugurServer)
	srv.storage = connect_to_storage()
	return srv
}
func main_page(c *gin.Context) {
	blknum,_:= augur_srv.storage.get_last_block_num()
	c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Augur Prediction Markets",
			"block_num" : blknum,
	})
//	c.String(http.StatusOK, "Hello %s %s", name,out)
}
func markets(c *gin.Context) {
	blknum,_ := augur_srv.storage.get_last_block_num()
	markets := augur_srv.storage.get_active_market_list()
	c.HTML(http.StatusOK, "markets.html", gin.H{
			"title": "Augur Markets",
			"block_num" : blknum,
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
			"title": "Augur Market Categories",
			"block_num" : blknum,
	})
}
func market_trades(c *gin.Context) {
	market := c.Param("market")
	trades := augur_srv.storage.get_mkt_trades(market)
	c.HTML(http.StatusOK, "market_info.html", gin.H{
			"title": "Trades for market",
			"Trades" : trades,
			"market_address": market,
	})
}
