// Copyright (c) 2017 The Alvalor Authors
//
// This file is part of Alvalor.
//
// Alvalor is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Alvalor is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Alvalor.  If not, see <http://www.gnu.org/licenses/>.

package network

import (
	"bytes"
	"net"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// Server represents the network component listening for incoming connections
// and performing the initial handshake to make sure we are dealing with a valid
// peer of our configured Alvalor network.
type Server struct {
	log       *zap.Logger
	wg        *sync.WaitGroup
	addresses <-chan string
	events    chan<- interface{}
	address   string
	network   []byte
	nonce     []byte
}

// NewServer will create a new server to listen for incoming peers and handling
// the handshake up to having a valid network connection for the given Alvalor
// network.
func NewServer(log *zap.Logger, wg *sync.WaitGroup, addresses <-chan string, events chan<- interface{}, options ...func(*Server)) *Server {
	server := &Server{
		log:       log,
		wg:        wg,
		addresses: addresses,
		events:    events,
		address:   "",
		network:   []byte{0, 0, 0, 0},
		nonce:     uuid.UUID{}.Bytes(),
	}
	for _, option := range options {
		option(server)
	}
	go server.listen()
	return server
}

// SetAddress allows us to define the local address we want to listen on with
// the server.
func SetAddress(address string) func(*Server) {
	return func(server *Server) {
		server.address = address
	}
}

// SetServerNetwork allows us to set the network to use during the initial connection
// handshake.
func SetServerNetwork(network []byte) func(*Server) {
	return func(server *Server) {
		server.network = network
	}
}

// SetServerNonce allows us to set our node nonce to make sure we never connect to
// ourselves.
func SetServerNonce(nonce []byte) func(*Server) {
	return func(server *Server) {
		server.nonce = nonce
	}
}

// Listen will start a listener on the configured network address and do the
// welcome handshake, forwarding valid peer connections.
func (server *Server) listen() {
	_, _, err := net.SplitHostPort(server.address)
	if err != nil {
		server.log.Error("invalid listen address", zap.String("server.address", server.address), zap.Error(err))
		return
	}
	ln, err := net.Listen("tcp", server.address)
	if err != nil {
		server.log.Error("could not create listener", zap.String("server.address", server.address), zap.Error(err))
		return
	}
Loop:
	for {
		tcpLn := ln.(*net.TCPListener)
		tcpLn.SetDeadline(time.Now().Add(1 * time.Second))
		conn, err := ln.Accept()
		netErr, ok := err.(*net.OpError)
		if ok && netErr.Timeout() {
			continue
		}
		if err != nil {
			server.log.Error("could not accept connection", zap.Error(err))
			break
		}
		address := conn.RemoteAddr().String()
		select {
		case _, ok := <-server.addresses:
			if !ok {
				break Loop
			}
		default:
			server.log.Info("no available connection slots", zap.String("address", address))
			conn.Close()
			continue
		}
		ack := append(server.network, server.nonce...)
		syn := make([]byte, len(ack))
		_, err = conn.Read(syn)
		if err != nil {
			server.log.Error("could not read syn packet", zap.Error(err))
			conn.Close()
			server.events <- Failure{Address: address}
			continue
		}
		network := syn[:len(server.network)]
		if !bytes.Equal(network, server.network) {
			server.log.Warn("dropping invalid network peer", zap.String("address", address), zap.ByteString("network", network))
			conn.Close()
			server.events <- Violation{Address: address}
			continue
		}
		nonce := syn[len(server.network):]
		if bytes.Equal(nonce, server.nonce) {
			server.log.Warn("dropping connection to self", zap.String("address", address))
			conn.Close()
			server.events <- Violation{Address: address}
			continue
		}
		_, err = conn.Write(ack)
		if err != nil {
			server.log.Error("could not write ack packet", zap.Error(err))
			conn.Close()
			server.events <- Failure{Address: address}
			continue
		}
		server.events <- Connection{Address: address, Conn: conn, Nonce: nonce}
	}
	server.wg.Done()
}
