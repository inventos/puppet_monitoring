// this's wrapper to call master process methods via GO RPC
package rpc

import (
	"log"
	"net/rpc"
	"os"
	"puppet_monitoring/impl"
)

// struct to group methods
type RPCClient struct {
	Conf *impl.Settings
}

// call GetStatus method and return info
func (p *RPCClient) GetStatus(with_errors bool) string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.GetStatus", &GetStatusArgs{Errors: with_errors}, &result)
	return result
}

// call RemoveNode method and return info
func (p *RPCClient) RemoveNode(host string) string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.RemoveNode", &RemoveNodeArgs{Host: host}, &result)
	return result
}

// call GetInfo method and return info
func (p *RPCClient) GetInfo() string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.GetInfo", &EmptyArgs{}, &result)
	return result
}

// call StopMasterProcess method and return success state and error info
func (p *RPCClient) StopMasterProcess() (bool, error) {
	cl := p.createClient()
	var result bool
	err := cl.Call("PPTRpc.StopMasterProcess", &EmptyArgs{}, &result)
	return result && err == nil, err
}

// private method to create common rpc dialup
func (c RPCClient) createClient() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", c.Conf.RpcComputed)
	if err != nil {
		log.Fatal("connect to master process failed:", err)
		os.Exit(1)
	}
	return client
}
