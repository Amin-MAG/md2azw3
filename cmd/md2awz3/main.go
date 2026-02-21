package main

import (
	"context"
	"fmt"
	"github.com/Amin-MAG/md2awz3/config"
	ravandlog "github.com/Amin-MAG/md2awz3/pkg/log"

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

Starting MD2AWZ3-%s...
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
	logger, err := ravandlog.NewLogger(ravandlog.Config{
		Level:          cfg.Logger.Level,
		IsPrettyPrint:  cfg.Logger.IsPrettyPrint,
		IsReportCaller: cfg.Logger.IsReportCallerMode,
	})
	if err != nil {
		panic(fmt.Errorf("can not initialize the ravandlog with error: %s", err))
	}
	ravandlog.SetupDefaultLogger(logger)
	logger.Info(ctx, "logger is setup successfully.")
	if cfg.Logger.IsPrettyPrint {
		logger.With("configuration", cfg.SecureClone()).Info(ctx, "launching Ravand")
	}

	//
}
