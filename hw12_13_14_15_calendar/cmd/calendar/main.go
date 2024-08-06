package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/app"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage/sql"
)

// Args command-line parameters.
type Args struct {
	ConfigPath string
}

func main() {
	var config Config
	args, f := ProcessArgs(&config)
	if f.Arg(0) == "version" {
		printVersion()
		return
	}
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &config); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	logg, err := NewLogger(config, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	store, err := CreateStorage(ctx, config)
	if err != nil {
		cancel()
		log.Fatal(err) //nolint:gocritic
	}
	defer store.Close()

	calendar := app.New(logg, store)
	server := internalhttp.NewServer(config.App.URL, logg, calendar)

	go func() {
		<-ctx.Done()
		if err := server.Stop(); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}

func NewLogger(cfg Config, w io.Writer) (*logger.Logger, error) {
	var level logger.Level
	lvl := strings.ToLower(cfg.Logger.Level)
	switch lvl {
	case "debug":
		level = logger.LevelDebug
	case "info":
		level = logger.LevelInfo
	case "warn":
		level = logger.LevelWarn
	case "error":
		level = logger.LevelError
	default:
		return nil, fmt.Errorf("invalid logger level: %s", cfg.Logger.Level)
	}

	return logger.New(level, w), nil
}

func CreateStorage(ctx context.Context, cfg Config) (storage.Repo, error) {
	ds := cfg.Datasource.SQL
	if cfg.Datasource.Type == "in-memory" {
		return memorystorage.New(), nil
	}
	s := sqlstorage.New()
	err := s.Connect(ctx,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			ds.Host, ds.Port, ds.Username, ds.Password, ds.Name, ds.Ssl))
	if err != nil {
		return nil, err
	}
	if err = s.RunMigration(ds.MigrationsDir); err != nil {
		return nil, err
	}

	return s, nil
}

func ProcessArgs(cfg interface{}) (Args, *flag.FlagSet) {
	var a Args

	f := flag.NewFlagSet("Calendar app", 1)
	f.StringVar(&a.ConfigPath, "config", "/etc/calendar/config.toml", "Path to configuration file")

	// Embed config descriptions into command help
	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a, f
}
