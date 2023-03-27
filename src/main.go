package main

import (
	"expvar"
	"fmt"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/obada-foundation/registry/cmd"
	_ "github.com/obada-foundation/registry/db"
	"github.com/obada-foundation/registry/system/logger"
	"go.uber.org/zap"
)

//	@title			OBADA DID Registry
//	@version		v0.1
//	@contact.name	techops@obada.io
//	@contact.url	https://www.obada.io
//	@host			registry.obada.io
//	@BasePath		/0.1/
//	@schemes		https http

var revision = "unknown"

type opts struct {
	ServerCmd cmd.ServerCommand `command:"server"`
}

func main() {
	fmt.Printf("DID Registry %s\n(c) OBADA Foundation %d\n\n", revision, time.Now().Year())

	log, err := logger.New("CLIENT-HELPER")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	var o opts

	if err := run(log, &o); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Errorw("startup", "ERROR", err)
			log.Sync()
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
