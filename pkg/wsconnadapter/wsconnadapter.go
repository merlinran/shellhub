package wsconnadapter

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	ws "nhooyr.io/websocket"
)

// an adapter for representing WebSocket connection as a net.Conn
// some caveats apply: https://github.com/gorilla/websocket/issues/441

var ErrUnexpectedMessageType = errors.New("unexpected websocket message type")

type Adapter struct {
	conn       *ws.Conn
	readMutex  sync.Mutex
	writeMutex sync.Mutex
	reader     io.Reader
}

func New(conn *ws.Conn) *Adapter {
	return &Adapter{
		conn: conn,
	}
}

func (a *Adapter) Read(b []byte) (int, error) {
	// Read() can be called concurrently, and we mutate some internal state here
	a.readMutex.Lock()
	defer a.readMutex.Unlock()

	if a.reader == nil {
		_, reader, _ := a.conn.Reader(context.TODO())
		a.reader = reader
	}

	bytesRead, err := a.reader.Read(b)
	if err != nil {
		a.reader = nil

		// EOF for the current Websocket frame, more will probably come so..
		if errors.Is(err, io.EOF) {
			// .. we must hide this from the caller since our semantics are a
			// stream of bytes across many frames
			err = nil
		}
	}

	return bytesRead, err
}

func (a *Adapter) Write(b []byte) (int, error) {
	a.writeMutex.Lock()
	defer a.writeMutex.Unlock()

	writer, _ := a.conn.Writer(context.TODO(), ws.MessageBinary)

	bytesWritten, err := writer.Write(b)
	writer.Close()

	return bytesWritten, err
}

func (a *Adapter) Close() error {
	//return a.conn.Close()
	return nil
}

type fakeAddr struct{}

func (fakeAddr) Network() string { return "revdial" }
func (fakeAddr) String() string  { return "revdialconn" }

func (a *Adapter) LocalAddr() net.Addr {
	return &fakeAddr{}
}

func (a *Adapter) RemoteAddr() net.Addr {
	return &fakeAddr{}
}

func (a *Adapter) SetDeadline(t time.Time) error {
	if err := a.SetReadDeadline(t); err != nil {
		return err
	}

	return a.SetWriteDeadline(t)
}

func (a *Adapter) SetReadDeadline(t time.Time) error {
	return nil
}

func (a *Adapter) SetWriteDeadline(t time.Time) error {
	a.writeMutex.Lock()
	defer a.writeMutex.Unlock()

	return nil
}
