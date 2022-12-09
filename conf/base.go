package conf

type Base struct {
	Name        string `mapstructure:"name"`
	LogDir      string `mapstructure:"log_dir"`
	Port        string
	PprofToken  string
	Debug       bool `mapstructure:"debug"`
	LogPosition bool `mapstructure:"log_position"`
	Pprof       bool
	Watch       bool
}

const (
	// FileName 配置文件名
	FileName = "conf"
	// LogPrefix 日志前缀
	LogPrefix = ""
)

var (
	DefaultConf []interface{}
)

func init() {
	DefaultConf = append(DefaultConf, Base{
		Name:  "ZlsApp",
		Debug: true,
		Port:  "8181",
	})
}
