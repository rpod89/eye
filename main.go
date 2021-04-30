// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
// STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF
// THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

//import "os"
import "fmt"
import "time"
import "crypto/sha1"
import "github.com/romana/rlog"

import "etcd.io/bbolt"

import "github.com/deroproject/derohe/rpc"

import "github.com/ybbus/jsonrpc"

import "oracle"

const PLUGIN_NAME = "basic_oracle_server"

const DEST_PORT = uint64(0xDEADBEEF12345678) // 16045690981402826360

var expected_arguments = rpc.Arguments{
	{rpc.RPC_DESTINATION_PORT, rpc.DataUint64, DEST_PORT},
}

var result = string("0.00")
var response = rpc.Arguments{
	{rpc.RPC_DESTINATION_PORT, rpc.DataUint64, uint64(0)},
	{rpc.RPC_SOURCE_PORT, rpc.DataUint64, DEST_PORT},
	{rpc.RPC_COMMENT, rpc.DataString, result},
}

var rpcClient1 = jsonrpc.NewClient("http://127.0.0.1:40404/json_rpc") // DUMMY WALLET that receive every requests (don't store dero on it)
var rpcClient2 = jsonrpc.NewClient("http://127.0.0.1:40405/json_rpc") // OPERATOR PRIVATE WALLET that receive rewards after each request

func main() {
	var err error
	fmt.Printf("basic oracle over dero chain.\n")
	var addr *rpc.Address
	var addr_result rpc.GetAddress_Result
	err = rpcClient1.CallFor(&addr_result, "GetAddress")
	if err != nil || addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return
	}

	if addr, err = rpc.NewAddress(addr_result.Address); err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", addr_result.Address, err)
		return
	}

	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(addr.String())))

	db_name := fmt.Sprintf("%s_%s.bbolt.db", PLUGIN_NAME, shasum)
	db, err := bbolt.Open(db_name, 0600, nil)
	if err != nil {
		fmt.Printf("could not open db err:%s\n", err)
		return
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("SEEK"))
		return err
	})
	if err != nil {
		fmt.Printf("err creating bucket. err %s\n", err)
	}

	fmt.Printf("Persistant store created in '%s'\n", db_name)

	service_address := addr.Clone()
	service_address.Arguments = expected_arguments
	fmt.Printf("Integrated address to activate '%s', service: \n%s\n", PLUGIN_NAME, service_address.String())

	processing_thread(db) // rkeep processing
}

func processing_thread(db *bbolt.DB) {

	var err error

	for { // currently we traverse entire history

		time.Sleep(time.Second)

		var transfers rpc.Get_Transfers_Result
		err = rpcClient1.CallFor(&transfers, "GetTransfers", rpc.Get_Transfers_Params{In: true, DestinationPort: DEST_PORT})
		if err != nil {
			rlog.Warnf("Could not obtain gettransfers from wallet err %s\n", err)
			continue
		}

		for _, e := range transfers.Entries {
			if e.Coinbase || !e.Incoming { // skip coinbase or outgoing, self generated transactions
				continue
			}

			// check whether the entry has been processed before, if yes skip it
			var already_processed bool
			db.View(func(tx *bbolt.Tx) error {
				if b := tx.Bucket([]byte("SEEK")); b != nil {
					if ok := b.Get([]byte(e.TXID)); ok != nil { // if existing in bucket
						already_processed = true
					}
				}
				return nil
			})

			if already_processed { // if already processed skip it
				continue
			}

			// check whether this service should handle the transfer
			if !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) ||
				DEST_PORT != e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64) {
				continue

			}

			rlog.Infof("tx should be processed %s\n", e.TXID)

			if e.Payload_RPC.Has("NAME", rpc.DataString) && e.Payload_RPC.Has("CURR", rpc.DataString) {
				name := e.Payload_RPC.Value("NAME", rpc.DataString).(string)
				curr := e.Payload_RPC.Value("CURR", rpc.DataString).(string)

				if len(name) == 0 || len(curr) == 0 {
					rlog.Warnf("One argument was missing NAME: [%s] or CURR: [%s] \n", name, curr)
					continue
				}

				price, err := oracle.SimplePrice(name, curr)
				if err != nil {
					// return the error instead of the price
					result = fmt.Sprintf("ERROR - coin:(%s) and currency:(%s) are unknowns, please try again!", name, curr)
				} else {
					// price (2 decimal) is send back
					result = fmt.Sprintf("%s %0.02f %s", name, price, curr)
				}

				response[1].Value = e.SourcePort
				response[2].Value = result

				var str string
				tparams := rpc.Transfer_Params{Transfers: []rpc.Transfer{{Destination: e.Sender, Amount: uint64(1), Payload_RPC: response}}}
				err = rpcClient2.CallFor(&str, "Transfer", tparams)
				if err != nil {
					rlog.Warnf("sending reply tx err %s\n", err)
					continue
				}

				err = db.Update(func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte("SEEK"))
					return b.Put([]byte(e.TXID), []byte("done"))
				})
				if err != nil {
					rlog.Warnf("err updating db to err %s\n", err)
				} else {
					rlog.Infof("oracle replied successfully with result: [ %s ]", result)
				}

			}
		}

	}
}
