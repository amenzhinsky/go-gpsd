package gpsd

import (
	"os"
	"testing"
)

func TestClass(t *testing.T) {
	for _, run := range []struct {
		b []byte
		w string
	}{
		{[]byte(`{"class":"VERSION","release":"3.16"}`), "VERSION"},
		{[]byte(`{"class":	"WATCH","enable":true}`), "WATCH"},
		{[]byte(`{"class": "TPV","device":"/dev/ttyUSB0"}`), "TPV"},
		{[]byte(`{   "class":   "SKY","device":"/dev/ttyUSB0"}`), "SKY"},
		{[]byte(`{"device":"/dev/ttyUSB0"}`), ""},
	} {
		if g := class(run.b); g != run.w {
			t.Errorf("class(%q) = %q, want %q", run.b, g, run.w)
		}
	}
}

func TestWatchJSON(t *testing.T) {
	addr := os.Getenv("TEST_ADDR")
	if addr == "" {
		t.Skip("TEST_ADDR is empty")
	}

	g, err := Dial(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()

	if err = g.Stream(WATCH_ENABLE|WATCH_JSON, ""); err != nil {
		t.Fatal(err)
	}

	// VERSION is received first when client connects
	r := <-g.C()
	if _, ok := r.(*VERSION); !ok {
		t.Fatalf("got %T, want a VERSION report", r)
	}
	v := r.(*VERSION)
	t.Logf("testing GPSD v%s, APIv%d.%d", v.Release, v.ProtoMajor, v.ProtoMinor)

	// DEVICES is received on WATCH_ENABLE
	r = <-g.C()
	if _, ok := r.(*DEVICES); !ok {
		t.Fatalf("got %T, want a DEVICES report", r)
	}

	if err = g.Err(); err != nil {
		t.Fatal(err)
	}
}
