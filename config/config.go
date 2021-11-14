package config

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/spf13/viper"
)

const (
	PRODUCTION = "PRODUCTION"
)

type targetServer struct {
	Address string
	Port    string
}

type dbServerInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
}


var (
	Environment            string
	CertificateFilePath    string
	CertificateKeyFilePath string
	WebDomainName          string
	ListenTo               targetServer
	CryptoKey              []byte
	LogPath                string
	DBInfo                 dbServerInfo
	DataPerPage            int
	CookieAuthName         string
	MaxSizeUploadPhotoByte int64
)

func init() {
	var err error
	envConfig := viper.New()

	envConfig.SetConfigName("env")
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath(".")

	if err = envConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	Environment = envConfig.GetString("environment")
	CertificateFilePath = envConfig.GetString("certificte_file_path")
	CertificateKeyFilePath = envConfig.GetString("certificte_key_file_path")
	ListenTo = targetServer{
		Address: envConfig.GetString("listen_address"),
		Port:    envConfig.GetString("listen_port"),
	}

	LogPath = envConfig.GetString("log_path")

	hasher := md5.New()
	hasher.Write([]byte(envConfig.GetString("crypto_key")))
	CryptoKey = []byte(hex.EncodeToString(hasher.Sum(nil)))

	WebDomainName = envConfig.GetString("web_domain_name")
	CookieAuthName = envConfig.GetString("cookie_auth_name")
	DBInfo.Host = envConfig.GetString("db_host")
	DBInfo.Port = envConfig.GetString("db_port")
	DBInfo.Username = envConfig.GetString("db_user")
	DBInfo.Password = envConfig.GetString("db_pass")
	DBInfo.DbName = envConfig.GetString("db_name")

	MaxSizeUploadPhotoByte = int64(envConfig.GetFloat64("max_size_upload_photo_mb") * 1000000)
}


