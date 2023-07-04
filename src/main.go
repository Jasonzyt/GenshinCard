package main

import (
	"flag"
	"os"

	"github.com/jasonzyt/genshincard/config"
	"github.com/jasonzyt/genshincard/image"
	"github.com/jasonzyt/genshincard/net"
)

type Flags struct {
	UseConfig    bool
	ConfigPath   string
	CookieSecret string
	//ServerMode   bool
	UserId        string
	SvgOutputPath string
}

func parseFlags() *Flags {
	flags := &Flags{}
	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	flagSet.BoolVar(&flags.UseConfig, "use-config", false, "use config file")
	flagSet.StringVar(&flags.ConfigPath, "config-path", "./config.yaml", "config file path")
	flagSet.StringVar(&flags.CookieSecret, "cookie-secret", "", "cookie secret")
	//flagSet.BoolVar(&flags.ServerMode, "server-mode", false, "http server mode")
	flagSet.StringVar(&flags.UserId, "user-id", "", "user id")
	flagSet.StringVar(&flags.SvgOutputPath, "output-path", "out.svg", "output path")
	flagSet.Parse(os.Args[1:])
	return flags
}

func detectRegion(uid string) string {
	if len(uid) == 0 {
		return ""
	}
	switch uid[0] {
	case '1':
		return "cn_gf01"
	case '2':
		return "cn_gf01"
	case '5':
		return "cn_qd01"
	case '6':
		return "os_usa"
	case '7':
		return "os_euro"
	case '8':
		return "os_asia"
	case '9':
		return "os_cht"
	}
	return ""
}

func detectApiType(region string) int {
	country := region[:2]
	switch country {
	case "cn":
		return net.MiYouSheApi
	case "os":
		return net.HoYoLabApi
	}
	return -1
}

func main() {

	flags := parseFlags()

	if flags.UseConfig {
		config.Init()
	} else {
		config.InitWithCookie(flags.CookieSecret)
	}

	if flags.UserId == "" {
		panic("user id is empty")
	}
	if flags.SvgOutputPath == "" {
		panic("output path is empty")
	}

	region := detectRegion(flags.UserId)
	if region == "" {
		panic("invalid user id")
	}

	apiType := detectApiType(region)
	if apiType == -1 {
		panic("invalid region")
	}

	profile, err := net.QueryPlayerProfile(apiType, region, flags.UserId)
	if err != nil {
		panic(err)
	}
	if profile == nil {
		panic("failed to query player profile")
	}

	outFile, err := os.OpenFile(flags.SvgOutputPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	outFile.Truncate(0)
	outFile.Seek(0, 0)

	image.GenerateSvg(profile, outFile)
}
