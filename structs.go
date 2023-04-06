package gpsd

import "time"

// Report is a grouping interface for all types of reports.
type Report interface {
	class() string
}

type TPV struct {
	Class  string    `json:"class"`
	Device string    `json:"device,omitempty"`
	Status float64   `json:"status,omitempty"`
	Mode   float64   `json:"mode"`
	Time   time.Time `json:"time,omitempty"`
	Ept    float64   `json:"ept,omitempty"`
	Lat    float64   `json:"lat,omitempty"`
	Lon    float64   `json:"lon,omitempty"`
	Alt    float64   `json:"alt,omitempty"`
	Eph    float64   `json:"eph,omitempty"`
	Epx    float64   `json:"epx,omitempty"`
	Epy    float64   `json:"epy,omitempty"`
	Epv    float64   `json:"epv,omitempty"`
	Track  float64   `json:"track,omitempty"`
	Speed  float64   `json:"speed,omitempty"`
	Climb  float64   `json:"climb,omitempty"`
	Epd    float64   `json:"epd,omitempty"`
	Eps    float64   `json:"eps,omitempty"`
	Epc    float64   `json:"epc,omitempty"`
}

func (r *TPV) class() string {
	return r.Class
}

type SKY struct {
	Class      string      `json:"class"`
	Device     string      `json:"device,omitempty"`
	Time       time.Time   `json:"time,omitempty"`
	Xdop       float64     `json:"xdop,omitempty"`
	Ydop       float64     `json:"ydop,omitempty"`
	Vdop       float64     `json:"vdop,omitempty"`
	Tdop       float64     `json:"tdop,omitempty"`
	Hdop       float64     `json:"hdop,omitempty"`
	Pdop       float64     `json:"pdop,omitempty"`
	Gdop       float64     `json:"gdop,omitempty"`
	Satellites []Satellite `json:"satellites"`
}

func (r *SKY) class() string {
	return r.Class
}

type Satellite struct {
	PRN  float64 `json:"PRN"`
	Az   float64 `json:"az"`
	El   float64 `json:"el"`
	Ss   float64 `json:"ss"`
	Used bool    `json:"used"`
}

type GST struct {
	Class  string    `json:"class"`
	Device string    `json:"device,omitempty"`
	Time   time.Time `json:"time,omitempty"`
	Rms    float64   `json:"rms,omitempty"`
	Major  float64   `json:"major,omitempty"`
	Minor  float64   `json:"minor,omitempty"`
	Orient float64   `json:"orient,omitempty"`
	Lat    float64   `json:"lat,omitempty"`
	Lon    float64   `json:"lon,omitempty"`
	Alt    float64   `json:"alt,omitempty"`
}

func (r *GST) class() string {
	return r.Class
}

type ATT struct {
	Class   string    `json:"class"`
	Device  string    `json:"device"`
	Time    time.Time `json:"time,omitempty"`
	Heading float64   `json:"heading,omitempty"`
	MagSt   string    `json:"mag_st,omitempty"`
	Pitch   float64   `json:"pitch,omitempty"`
	PitchSt string    `json:"pitch_st,omitempty"`
	Yaw     float64   `json:"yaw,omitempty"`
	YawSt   string    `json:"yaw_st,omitempty"`
	Roll    float64   `json:"roll,omitempty"`
	RollSt  string    `json:"roll_st,omitempty"`
	Dip     float64   `json:"dip,omitempty"`
	MagLen  float64   `json:"mag_len,omitempty"`
	MagX    float64   `json:"mag_x,omitempty"`
	MagY    float64   `json:"mag_y,omitempty"`
	MagZ    float64   `json:"mag_z,omitempty"`
	AccLen  float64   `json:"acc_len,omitempty"`
	AccX    float64   `json:"acc_x,omitempty"`
	AccY    float64   `json:"acc_y,omitempty"`
	AccZ    float64   `json:"acc_z,omitempty"`
	GyroX   float64   `json:"gyro_x,omitempty"`
	GyroY   float64   `json:"gyro_y,omitempty"`
	Depth   float64   `json:"depth,omitempty"`
	Temp    float64   `json:"temp,omitempty"`
}

func (r *ATT) class() string {
	return r.Class
}

type VERSION struct {
	Class      string  `json:"class"`
	Release    string  `json:"release"`
	Rev        string  `json:"rev"`
	ProtoMajor float64 `json:"proto_major"`
	ProtoMinor float64 `json:"proto_minor"`
	Remove     string  `json:"remote,omitempty"`
}

func (r *VERSION) class() string {
	return r.Class
}

type DEVICES struct {
	Class   string   `json:"class"`
	Devices []DEVICE `json:"devices"`
	Remote  string   `json:"remote,omitempty"`
}

func (r *DEVICES) class() string {
	return r.Class
}

type WATCH struct {
	Class   string  `json:"class"`
	Enable  bool    `json:"enable,omitempty"`
	Json    bool    `json:"json,omitempty"`
	Nmea    bool    `json:"nmea,omitempty"`
	Raw     float64 `json:"raw,omitempty"`
	Scaled  bool    `json:"scaled,omitempty"`
	Split24 bool    `json:"split24,omitempty"`
	Pps     bool    `json:"pps,omitempty"`
	Device  string  `json:"device,omitempty"`
	Remote  string  `json:"remote,omitempty"`
}

func (r *WATCH) class() string {
	return r.Class
}

type POLL struct {
	Class  string    `json:"class"`
	Time   time.Time `json:"time"`
	Active float64   `json:"active"`
	Tpv    []TPV     `json:"tpv"`
	Sky    []SKY     `json:"sky"`
	Gst    []GST     `json:"gst"`
}

func (r *POLL) class() string {
	return r.Class
}

type TOFF struct {
	Class     string  `json:"class"`
	Device    string  `json:"device"`
	RealSec   float64 `json:"real_sec"`
	RealNSec  float64 `json:"real_nsec"`
	ClockSec  float64 `json:"clock_sec"`
	ClockNSec float64 `json:"clock_nsec"`
}

func (r *TOFF) class() string {
	return r.Class
}

type PPS struct {
	Class     string  `json:"class"`
	Device    string  `json:"device"`
	RealSec   float64 `json:"real_sec"`
	RealNSec  float64 `json:"real_nsec"`
	ClockSec  float64 `json:"clock_sec"`
	ClockNSec float64 `json:"clock_nsec"`
	Precision float64 `json:"precision"`
}

func (r *PPS) class() string {
	return r.Class
}

type OSC struct {
	Class       string  `json:"class"`
	Device      string  `json:"device"`
	Running     bool    `json:"running"`
	Reference   bool    `json:"reference"`
	Disciplined bool    `json:"disciplined"`
	Delta       float64 `json:"delta"`
}

func (r *OSC) class() string {
	return r.Class
}

type DEVICE struct {
	Class string `json:"class"`

	// can be string whet it's a member of DEVICES and float64 when used separately
	Activated interface{} `json:"activated,omitempty"`
	Path      string      `json:"path,omitempty"`
	Flags     float64     `json:"flags,omitempty"`
	Driver    string      `json:"driver,omitempty"`
	Subtype   string      `json:"subtype,omitempty"`
	Bps       float64     `json:"bps,omitempty"`
	Parity    string      `json:"parity,omitempty"`
	Stopbits  float64     `json:"stopbits"`
	Native    float64     `json:"native,omitempty"`
	Cycle     float64     `json:"cycle,omitempty"`
	Mincycle  float64     `json:"mincycle,omitempty"`
}

func (r *DEVICE) class() string {
	return r.Class
}

type ERROR struct {
	Class   string `json:"class"`
	Message string `json:"message"`
}

func (r *ERROR) class() string {
	return r.Class
}

type RAW []byte

func (r RAW) class() string {
	return "RAW"
}

type SUBFRAME struct {
	Class  string `json:"class"`
	Device string `json:"device"`
	GnssId int    `json:"gnssId"`
	TSV    int    `json:"tSV"`
	TOW17  int    `json:"TOW17"`
	Frame  int    `json:"frame,omitempty"`
	Scaled bool   `json:"scaled"`
}

func (r *SUBFRAME) class() string {
	return r.Class
}

type RTCM2 struct {
	Class         string  `json:"class"`
	Type          int     `json:"type"`
	StationId     int     `json:"station_id"`
	Zcount        float64 `json:"zcount"`
	Seqnum        int     `json:"seqnum"`
	Length        int     `json:"length"`
	StationHealth int     `json:"station_health"`
	X             float64 `json:"x,omitempty"`
	Y             float64 `json:"y,omitempty"`
	Z             float64 `json:"z,omitempty"`
	Satellites    []struct {
		Ident      int  `json:"ident"`
		Iodl       bool `json:"iodl"`
		Health     int  `json:"health"`
		Snr        int  `json:"snr"`
		HealthEn   bool `json:"health_en"`
		NewData    bool `json:"new_data"`
		LosWarning bool `json:"los_warning"`
		Tou        int  `json:"tou"`
	} `json:"satellites,omitempty"`
	System  string  `json:"system,omitempty"`
	Sense   int     `json:"sense,omitempty"`
	Datum   string  `json:"datum,omitempty"`
	Dx      float64 `json:"dx,omitempty"`
	Dy      float64 `json:"dy,omitempty"`
	Dz      float64 `json:"dz,omitempty"`
	Message string  `json:"message,omitempty"`
	Ar      string  `json:"ar,omitempty"`
	Sid     string  `json:"sid,omitempty"`
}

func (r *RTCM2) class() string {
	return r.Class
}
