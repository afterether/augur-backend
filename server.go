package main
import (
	"fmt"
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
	blknum,res := augur_srv.storage.get_last_block_num()
	_ = res
	markets := augur_srv.storage.get_active_market_list()
	c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Augur Prediction Markets",
			"block_num" : blknum,
			"Markets" : markets,
	})
//	c.String(http.StatusOK, "Hello %s %s", name,out)
}
func contract_info(c *gin.Context) {
	name := c.Param("contract")
	blknum,res := augur_srv.storage.get_last_block_num()
	_ = res
	out := fmt.Sprintf("block = %v",blknum)
	c.String(http.StatusOK, "Hello %s %s", name,out)
}
