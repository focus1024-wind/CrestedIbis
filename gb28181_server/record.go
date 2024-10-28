package gb28181_server

type Record struct {
	App       string  `json:"app" xml:"app" desc:"录制流的应用名"`
	FileName  string  `json:"file_name" xml:"file_name" desc:"录制文件名"`
	FileSize  int64   `json:"file_size" xml:"file_size" desc:"文件大小，单位字节"`
	StartTime int64   `json:"start_time" xml:"start_time" desc:"开始录制时间戳"`
	Stream    string  `json:"stream" xml:"stream" desc:"录制流ID"`
	TimeLen   float64 `json:"time_len" xml:"time_len" desc:"录制时长，单位: 秒"`
	Url       string  `json:"url" xml:"url" desc:"点播相对Url路径"`
}
