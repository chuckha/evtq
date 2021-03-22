package infrastructure

import (
	"bytes"
	"net"

	"github.com/pkg/errors"

	"github.com/chuckha/evtq/core"
)

type Offsets interface {
	GetOffsets(name string) []*core.Offset
}

type ConnectorBuilder struct {
	OffsetRepo Offsets
}

func (c *ConnectorBuilder) BuildConnector(info *core.ConnectorBuilderInfo) (core.Connector, error) {
	offsets := c.OffsetRepo.GetOffsets(info.Name)
	switch v := info.Info.(type) {
	case LocalConnectorInfo:
		return NewLocalConnector(info, offsets)
	case TCPConnectorInfo:
		return NewTCPConnector(info, offsets)
	default:
		return nil, errors.Errorf("unsupported info type %T", v)
	}
}

type LocalConnectorInfo struct{}

func (l LocalConnectorInfo) Info() {}

type TCPConnectorInfo struct {
	Network string
	Address string
}

func (t TCPConnectorInfo) Info() {}

func NewTCPConnector(info *core.ConnectorBuilderInfo, offsets []*core.Offset) (core.Connector, error) {
	tcpInfo, ok := info.Info.(TCPConnectorInfo)
	if !ok {
		return nil, errors.Errorf("info needs to be of type TCPConnectorInfo, not %T", info.Info)
	}
	// TODO dial with retries
	conn, err := net.Dial(tcpInfo.Network, tcpInfo.Address)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return core.NewConnector(info.Name, conn, EncDecFactory(info.EncodingType), offsets)
}

func NewLocalConnector(info *core.ConnectorBuilderInfo, offsets []*core.Offset) (core.Connector, error) {
	_, ok := info.Info.(LocalConnectorInfo)
	if !ok {
		return nil, errors.Errorf("info needs to be of type LocalConnectorInfo, not %T", info)
	}
	var IO bytes.Buffer
	return core.NewConnector(info.Name, &IO, EncDecFactory(info.EncodingType), offsets)
}
