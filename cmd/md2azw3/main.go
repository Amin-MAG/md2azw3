package main

import (
	"context"
	"fmt"

	"github.com/Amin-MAG/md2azw3/config"
	"github.com/Amin-MAG/md2azw3/internal/server"
	mdlog "github.com/Amin-MAG/md2azw3/pkg/log"

	"github.com/ilyakaznacheev/cleanenv"
)

const Banner = `
::::    ::::  :::::::::        ::::::::           :::     :::       ::: ::::::::: ::::::::  
+:+:+: :+:+:+ :+:    :+:      :+:    :+:        :+: :+:   :+:       :+:      :+: :+:    :+: 
+:+ +:+:+ +:+ +:+    +:+            +:+        +:+   +:+  +:+       +:+     +:+         +:+ 
+#+  +:+  +#+ +#+    +:+          +#+         +#++:++#++: +#+  +:+  +#+    +#+       +#++:  
+#+       +#+ +#+    +#+        +#+           +#+     +#+ +#+ +#+#+ +#+   +#+           +#+ 
#+#       #+# #+#    #+#       #+#            #+#     #+#  #+#+# #+#+#   #+#     #+#    #+# 
###       ### #########       ##########      ###     ###   ###   ###   ######### ########  

Starting MD2AZW3-%s...
Configuration: %+v
`

var cfg config.Config

func main() {
	ctx := context.Background()

	// Load the config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err.Error())
	}

	// Banner
	fmt.Printf(Banner, config.AppVersion, cfg.SecureClone())

	// Configure the log
	logger, err := mdlog.NewLogger(mdlog.Config{
		Level:          cfg.Logger.Level,
		IsPrettyPrint:  cfg.Logger.IsPrettyPrint,
		IsReportCaller: cfg.Logger.IsReportCallerMode,
	})
	if err != nil {
		panic(fmt.Errorf("can not initialize the mdlog with error: %s", err))
	}
	mdlog.SetupDefaultLogger(logger)
	logger.Info(ctx, "logger is setup successfully.")
	if cfg.Logger.IsPrettyPrint {
		logger.With("configuration", cfg.SecureClone()).Info(ctx, "launching Ravand")
	}

	// Start HTTP server
	e := server.New(cfg, logger)
	if err = server.Start(e, cfg, logger); err != nil {
		logger.WithError(err).Fatal(ctx, "HTTP server failed")
	}
}
