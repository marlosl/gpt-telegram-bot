package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/marlosl/gpt-telegram-bot/consts"
	"github.com/marlosl/gpt-telegram-bot/utils/config"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var configVars = []string{
	consts.AwsRegion,
	consts.AwsAccessKeyId,
	consts.AwsSecretAccessKey,
	consts.PulumiAccessToken,
	consts.ProjectDir,
	consts.ProjectOutputDir,
	consts.CloudfareApiToken,
	consts.DnsZone,
	consts.DnsRecord,
	consts.CloudfareApiKey,
	consts.CloudfareApiEmail,
	consts.GptApiKey,
	consts.TelegramBotTextToken,
	consts.TelegramBotImageToken,
}

func InitConfig() {
	viper.AddConfigPath(GetExecutablePath())
	viper.SetConfigType("env")
	viper.SetConfigName(".config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Can't config file: %v", err)
		os.Exit(1)
	}

	InitEnvVars()
	config.NewConfig(config.File)
}

func InitEnvVars() {
	for _, envVar := range configVars {
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, viper.GetString(envVar))
		}
	}
}

func GetConfigValue(key string) string {
	return viper.GetString(key)
}

func GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func PrintJson(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func SPrintJson(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return string(b)
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func PrintRestyDebug(resp *resty.Response, err error) {
	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()

	fmt.Println("  Request URL   :", resp.Request.RawRequest.URL)
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
}
