//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package envs

import (
	"errors"
	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	"github.com/lastbackend/lastbackend/pkg/network"
	"github.com/lastbackend/lastbackend/pkg/node/state"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii"
	"github.com/lastbackend/lastbackend/pkg/runtime/cri"
	"github.com/lastbackend/lastbackend/pkg/runtime/csi"
)

var e Env

func Get() *Env {
	return &e
}

type Env struct {
	cri cri.CRI
	cii cii.CII
	csi map[string]csi.CSI

	state  *state.State
	net *network.Network

	client struct {
		node types.NodeClientV1
		rest types.ClientV1
	}

	mode struct {
		provision bool
	}
}

func (c *Env) SetCRI(cri cri.CRI) {
	c.cri = cri
}

func (c *Env) GetCRI() cri.CRI {
	return c.cri
}

func (c *Env) SetCII(iri cii.CII) {
	c.cii = iri
}

func (c *Env) GetCII() cii.CII {
	return c.cii
}

func (c *Env) SetNet(net *network.Network) {
	c.net = net
}

func (c *Env) GetNet() *network.Network {
	return c.net
}

func (c *Env) SetProvision(on bool) {
	c.mode.provision = on
}

func (c *Env) GetProvision() bool {
	return c.mode.provision
}

func (c *Env) SetCSI(kind string, si csi.CSI) {
	c.csi = make(map[string]csi.CSI)
	c.csi[kind] = si
}

func (c *Env) ListCSI () []string {
	var types = []string{}

	for t := range c.csi {
		types = append(types, t)
	}
	return types
}

func (c *Env) GetCSI(kind string) (csi.CSI, error) {
	if _, ok := c.csi[kind]; !ok {
		return nil, errors.New("storage container interface not supported")
	}
	return c.csi[kind], nil
}

func (c *Env) SetState(s *state.State) {
	c.state = s
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetClient(nc types.NodeClientV1, rc types.ClientV1) {
	c.client.node = nc
	c.client.rest = rc
}

func (c *Env) GetNodeClient() types.NodeClientV1 {
	return c.client.node
}

func (c *Env) GetRestClient() types.ClientV1 {
	return c.client.rest
}
