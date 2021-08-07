package mjolnir

import (
	"fmt"
	"github.com/complone/blast/common"
	"github.com/complone/blast/http_api"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

type option struct {
	Verbose bool   `short:"v" long:"verbose" description:"Show verbose debug message"`
	Config  string `short:"c" description:"Specified config file" default:"../conf/openapi.json"`
}

func RegisterAPIHandle(apiNode string, handles http_api.ModuleHandles) {
	http_api.ModuleHandleContainer[apiNode] = handles
}

func RunApp(cfg common.IConfig) int {
	if common.AppName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "common.AppName must set in build.sh\n")
		return -1
	}

	initLog()

	var opt option
	p := flags.NewParser(&opt, flags.Default & ^flags.PrintErrors)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				_, _ = fmt.Fprintln(os.Stdout, err)
				return 0
			default:
				_, _ = fmt.Fprintf(os.Stderr, "error when parsing command: %s\n", err)
				return -1
			}
		}
	}

	if opt.Verbose {
		printVersion()
		return 0
	}

	err := common.InitConfigFromJson(opt.Config, cfg)
	if err != nil {
		log.Warningf("InitConfigFromJson failed. err = %s", err)
		return -1
	}

	err = http_api.SyncJsonSchema(common.GlobalConfig.OpenAPI.SchemaPath)
	if err != nil {
		log.Warningf("SyncJsonSchema failed. err = %s", err)
		return -1
	}

	pidFile := fmt.Sprintf("%s.pid", common.AppName)
	_ = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0666)

	log.Infof("%s begin serve on addr %s with pid %d",
		common.AppName,
		common.GlobalConfig.OpenAPI.HttpAddr,
		os.Getpid())

	err = http_api.RunGin(common.GlobalConfig.OpenAPI.HttpAddr)
	if err != nil {
		log.Warningf("%s run failed with pid %d, err %s", common.AppName, os.Getpid(), err)
		return -1
	}

	log.Infof("%s stopped with pid %d", common.AppName, os.Getpid())
	return 0
}

func printVersion() {
	fmt.Printf("%-12s\t%s\n", "AppName:", common.AppName)
	fmt.Printf("%-12s\t%s\n", "AppVersion:", common.AppVersion)
	fmt.Printf("%-12s\t%s\n", "Built:", common.BuildTime)
	fmt.Printf("%-12s\t%s\n", "GoVersion:", common.GoVersion)
}

func initLog() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := dir + "/../logs/" + common.AppName + ".log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(24*7)*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	log.SetOutput(writer)
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(common.NewLogFormatter())

	gin.DefaultWriter = log.StandardLogger().Out
}