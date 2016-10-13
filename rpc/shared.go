// shared structs between client wrappers and server
package rpc

// aka namespace for GO RPC
type PPTRpc string

// empty args (nil isn't good)
type EmptyArgs struct{}

// arg to describe remove node params
type RemoveNodeArgs struct {
	// Host to remove
	Host string
}

// arg to describe status node params
type GetStatusArgs struct {
	// print out errors while analyzing status
	// default: false
	Errors bool
}
