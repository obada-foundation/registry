package main

import (
	"expvar"
	"fmt"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/obada-foundation/registry/cmd"
	log "github.com/obada-foundation/registry/system/logger"
	"go.uber.org/zap"
)

var revision = "unknown"

type opts struct {
	ServerCmd cmd.ServerCommand `command:"server"`
	ClientCmd cmd.ClientCommand `command:"client"`
}

func main() {
	fmt.Printf("DID Registry %s\n(c) OBADA Foundation %d\n\n", revision, time.Now().Year())

	logger, err := log.New("CLIENT-HELPER")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	var o opts

	if err := run(logger, &o); err != nil {
		if _, ok := err.(*flags.Error); !ok {
			logger.Errorw("startup", "ERROR", err)
			if err := logger.Sync(); err != nil {
				logger.Errorw("startup", "Sync LOG ERROR", err)
			}

			// nolint
			os.Exit(1)
		}
	}
}

func run(logger *zap.SugaredLogger, o *opts) error {
	expvar.NewString("build").Set(revision)

	p := flags.NewParser(o, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOpts{
			Revision: revision,
			Logger:   logger,
		})

		err := command.Execute(args)
		if err != nil {
			logger.Errorf("failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		return err
	}

	return nil
}
