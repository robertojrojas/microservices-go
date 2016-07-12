package services

import "golang.org/x/net/context"

type RPCRequest struct {
	Ctx context.Context
}

type RPCResponse struct {
	Key  string
	Data []byte
}

// RPCService is an abstraction of the
// RPC to be performed.
type RPCService interface {
	RPC(*RPCRequest) (*RPCResponse, error)
}
