package protoudp

import (
	"net"
	"strconv"
)

const (
	bufferSize = 2048
)

type Server struct {
	ln net.PacketConn
}

func NewServer(port int) (*Server, error) {
	ln, err := net.ListenPacket("udp", ":"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		return nil, err
	}

	return &Server{
		ln: ln,
	}, nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Port() int {
	return s.ln.LocalAddr().(*net.UDPAddr).Port
}

func (s *Server) ReadFrame() (*Frame, *net.UDPAddr, error) {
	buf := make([]byte, bufferSize)

	n, source, err := s.ln.ReadFrom(buf)
	if err != nil {
		return nil, nil, err
	}

	f, err := frameDecode(buf[:n])
	if err != nil {
		return nil, nil, err
	}

	return f, source.(*net.UDPAddr), nil
}

func (s *Server) WriteFrame(f *Frame, dest *net.UDPAddr) error {
	byts, err := frameEncode(f)
	if err != nil {
		return err
	}

	_, err = s.ln.WriteTo(byts, dest)
	return err
}
