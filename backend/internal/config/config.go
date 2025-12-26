package config

type Config struct {
	Server   Server   `mapstructure:"server" json:"server" yaml:"server"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Cors     Cors     `mapstructure:"cors" json:"cors" yaml:"cors"`
	Storage  Storage  `mapstructure:"storage" json:"storage" yaml:"storage"`
}

type Server struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`
}

type Database struct {
	Driver  string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Source  string `mapstructure:"source" json:"source" yaml:"source"`
	LogMode string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
}

type JWT struct {
	SigningKey  string `mapstructure:"signing_key" json:"signing_key" yaml:"signing_key"`
	ExpiresTime string `mapstructure:"expires_time" json:"expires_time" yaml:"expires_time"`
	BufferTime  string `mapstructure:"buffer_time" json:"buffer_time" yaml:"buffer_time"`
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director" yaml:"director"`
	ShowLine      bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	EncodeLevel   string `mapstructure:"encode_level" json:"encode_level" yaml:"encode_level"`
	StacktraceKey string `mapstructure:"stacktrace_key" json:"stacktrace_key" yaml:"stacktrace_key"`
	LogInConsole  bool   `mapstructure:"log_in_console" json:"log_in_console" yaml:"log_in_console"`
}

type Cors struct {
	Mode      string      `mapstructure:"mode" json:"mode" yaml:"mode"`
	Whitelist []Whitelist `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`
}

type Whitelist struct {
	AllowOrigin  string `mapstructure:"allow-origin" json:"allow_origin" yaml:"allow-origin"`
	AllowMethods string `mapstructure:"allow-methods" json:"allow_methods" yaml:"allow-methods"`
	AllowHeaders string `mapstructure:"allow-headers" json:"allow_headers" yaml:"allow-headers"`
}

// Storage 存储配置
type Storage struct {
	Driver string       `mapstructure:"driver" json:"driver" yaml:"driver"` // local | oss | cos | r2
	Local  LocalStorage `mapstructure:"local" json:"local" yaml:"local"`
	OSS    OSSStorage   `mapstructure:"oss" json:"oss" yaml:"oss"`
	COS    COSStorage   `mapstructure:"cos" json:"cos" yaml:"cos"`
	R2     R2Storage    `mapstructure:"r2" json:"r2" yaml:"r2"`
}

// LocalStorage 本地存储配置
type LocalStorage struct {
	Path   string `mapstructure:"path" json:"path" yaml:"path"`
	Domain string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

// OSSStorage 阿里云OSS配置
type OSSStorage struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain          string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

// COSStorage 腾讯云COS配置
type COSStorage struct {
	Region    string `mapstructure:"region" json:"region" yaml:"region"`
	SecretID  string `mapstructure:"secret_id" json:"secret_id" yaml:"secret_id"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

// R2Storage Cloudflare R2配置
type R2Storage struct {
	AccountID       string `mapstructure:"account_id" json:"account_id" yaml:"account_id"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Domain          string `mapstructure:"domain" json:"domain" yaml:"domain"` // 自定义域名或 R2.dev 域名
}
