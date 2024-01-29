package gpsd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
)

// DefaultAddress is the default address gpsd listens to.
const DefaultAddress = ":2947"

// Option is a client configuration option.
type Option func(c *Client)

// WithConn sets client connection.
func WithConn(conn net.Conn) Option {
	return func(c *Client) {
		c.conn = conn
	}
}

// WithLogger sets logger.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithChannel changes reports channel for adjusting buffer size, default is 10.
func WithChannel(ch chan Report) Option {
	return func(c *Client) {
		c.evch = ch
	}
}

// New creates a new gpsd client.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		evch:   make(chan Report, 1),
		done:   make(chan struct{}),
		logger: newStdLogger(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.conn == nil {
		return nil, errors.New("conn is nil")
	}
	go c.rx()
	return c, nil
}

// Dial dials the named address and returns a gpsd client for the connection.
func Dial(addr string, opts ...Option) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return New(append(opts, WithConn(conn))...)
}

// Client is a gpsd client implementation.
//
// https://gpsd.gitlab.io/gpsd/gpsd_json.html
// https://gpsd.gitlab.io/gpsd/client-howto.html
type Client struct {
	mu     sync.Mutex
	err    error
	conn   net.Conn
	evch   chan Report
	done   chan struct{}
	logger Logger
}

func (c *Client) rx() {
	defer close(c.evch)

	r := bufio.NewReader(c.conn)

	var b []byte
	var err error
	var v Report
	for {
		b, err = r.ReadSlice('\n')
		if err != nil {
			_ = c.close(err)
			return
		}
		// trim '\r' if present
		if b[len(b)-1] == '\r' {
			b = b[:len(b)-2]
		}
		if b[0] == '{' { // json
			c.debugf("RX %s", b)
			v, err = unmarshal(class(b), b)
			if err != nil {
				c.errorf("unmarshal error: %q", err)
				continue
			}
		} else {
			c.debugf("RX [RAW]")
			v = RAW(b) // raw data
		}
		select {
		case c.evch <- v:
		case <-c.done:
			return
		}
	}
}

// C returns a channel where report structures can be read from.
func (c *Client) C() <-chan Report {
	return c.evch
}

// Done reports where the client is stopped.
func (c *Client) Done() <-chan struct{} {
	return c.done
}

// Err returns client errors, makes sense to call it only when C is closed.
func (c *Client) Err() error {
	return c.err
}

// WatchFlag configures watch parameters.
type WatchFlag uint32

const (
	// WATCH_ENABLE enable streaming.
	WATCH_ENABLE = 0x0001

	// WATCH_DISABLE disable watching.
	WATCH_DISABLE = 0x0002

	// WATCH_JSON JSON output.
	WATCH_JSON = 0x0010

	// WATCH_NMEA output in NMEA.
	WATCH_NMEA = 0x0020

	// WATCH_RARE output of packets in hex.
	WATCH_RARE = 0x0040

	// WATCH_RAW output of raw packets.
	WATCH_RAW = 0x0080

	// WATCH_SCALED scale output to floats.
	WATCH_SCALED = 0x0100

	// WATCH_TIMING timing information.
	WATCH_TIMING = 0x0200

	// WATCH_DEVICE watch specific device.
	WATCH_DEVICE = 0x0800

	// WATCH_SPLIT24 split AIS Type 24s.
	WATCH_SPLIT24 = 0x1000

	// WATCH_PPS enable PPS JSON.
	WATCH_PPS = 0x2000
)

// Stream changes watch policy the second argument can only
// be a device path and it only implies WATCH_DEVICE flag.
func (c *Client) Stream(flags WatchFlag, devpath string) error {
	// cannot use WATCH struct unless all its attributes are pointers
	// that will force users to dereference each value and check if it's not a nil,
	// otherwise json.Marshall will discard zero values because of the omitempty option
	// or will send zero values every time if we get rid of it.
	s := "?WATCH={"
	if flags&WATCH_DISABLE != 0 {
		s += `"enable":false`
		if flags&WATCH_JSON != 0 {
			s += `,"json":false`
		}
		if flags&WATCH_NMEA != 0 {
			s += `,"nmea":false`
		}
		if flags&WATCH_RARE != 0 {
			s += `,"raw":1`
		}
		if flags&WATCH_RAW != 0 {
			s += `,"raw":2`
		}
		if flags&WATCH_SCALED != 0 {
			s += `,"scaled":false`
		}
		if flags&WATCH_TIMING != 0 {
			s += `,"scaled":false`
		}
		if flags&WATCH_SPLIT24 != 0 {
			s += `,"split24":false`
		}
		if flags&WATCH_PPS != 0 {
			s += `,"pps":false`
		}
	} else { // flags&WATCH_ENABLE
		s += `"enable":true`
		if flags&WATCH_JSON != 0 {
			s += `,"json":true`
		}
		if flags&WATCH_NMEA != 0 {
			s += `,"nmea":true`
		}
		if flags&WATCH_RARE != 0 {
			s += `,"raw":1`
		}
		if flags&WATCH_RAW != 0 {
			s += `,"raw":2`
		}
		if flags&WATCH_SCALED != 0 {
			s += `,"scaled":true`
		}
		if flags&WATCH_TIMING != 0 {
			s += `,"scaled":true`
		}
		if flags&WATCH_SPLIT24 != 0 {
			s += `,"split24":true`
		}
		if flags&WATCH_PPS != 0 {
			s += `,"pps":true`
		}
		if flags&WATCH_DEVICE != 0 {
			s += fmt.Sprintf(`,"device":%q`, devpath)
		}
	}
	s += "}"
	return c.Send([]byte(s))
}

// Send sends the given raw data to gpsd,
// use it only if you know what you're doing.
//
// Do not use marshalled report structs here, because they're
// supposed to be used to receive data from the daemon only.
func (c *Client) Send(b []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return c.err
	}
	if _, err := c.conn.Write(b); err != nil {
		_ = c.close(err)
		return err
	}
	c.debugf("TX %s", b)
	return nil
}

func (c *Client) errorf(format string, v ...interface{}) {
	if c.logger != nil {
		c.logger.Errorf(format, v...)
	}
}

func (c *Client) debugf(format string, v ...interface{}) {
	if c.logger != nil {
		c.logger.Debugf(format, v...)
	}
}

// ErrClosed signalizes that client is closed.
var ErrClosed = errors.New("closed")

func (c *Client) close(err error) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return nil
	default:
	}
	c.err = err
	close(c.done)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Close closes underlying gpsd connection.
func (c *Client) Close() error {
	return c.close(ErrClosed)
}

// unmarshal unmarshals the given json into a report corresponding to the class attribute.
func unmarshal(cls string, b []byte) (Report, error) {
	var v Report
	switch cls {
	case "TPV":
		v = &TPV{}
	case "SKY":
		v = &SKY{}
	case "GST":
		v = &GST{}
	case "ATT":
		v = &ATT{}
	case "VERSION":
		v = &VERSION{}
	case "DEVICES":
		v = &DEVICES{}
	case "WATCH":
		v = &WATCH{}
	case "POLL":
		v = &POLL{}
	case "TOFF":
		v = &TOFF{}
	case "PPS":
		v = &PPS{}
	case "OSC":
		v = &OSC{}
	case "DEVICE":
		v = &DEVICE{}
	case "SUBFRAME":
		v = &SUBFRAME{}
	case "RTCM2":
		v = &RTCM2{}
	case "ERROR":
		v = &ERROR{}
	default:
		return nil, fmt.Errorf("unknown class %q", cls)
	}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}

// class detects the given json report class name or returns an empty string when it fails.
// It's needed to avoid parsing json into map[string]interface{} for performance reasons.
func class(b []byte) string {
	const prefix = `"class":"`
	for i, j := 0, 0; i < len(b); i++ {
		// skip whitespace chars
		if b[i] == ' ' || b[i] == '\r' || b[i] == '\n' || b[i] == '\t' {
			continue
		}
		if b[i] == prefix[j] {
			// prefix detected
			if len(prefix) == j+1 {
				for k := i + 1; k < len(b); k++ {
					// end of name reached
					if b[k] == '"' {
						return string(b[i+1 : k])
					}
				}
				// '"' is not found
				break
			}
			j++
		} else {
			j = 0
		}
	}
	return ""
}
