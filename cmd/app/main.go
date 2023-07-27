package main

import (
	"flag"
	"log"

	"github.com/antonmisa/1cctl_cli/config"
	"github.com/antonmisa/1cctl_cli/internal/app"
)

func main() {
	var prepare bool

	flag.BoolVar(&prepare, "prepare", false, "creating default environment and config")

	var clusterConnection string

	flag.StringVar(&clusterConnection, "clusterConnection", "localhost:1545", "cluster host:port to connect in cli mode")

	var clusterName string

	flag.StringVar(&clusterName, "clusterName", "localhost:1541", "cluster host:port to make a backup in cli mode")

	var clusterAdmin string

	flag.StringVar(&clusterAdmin, "clusterAdmin", "", "cluster admin name")

	var clusterPwd string

	flag.StringVar(&clusterPwd, "clusterPwd", "", "cluster password")

	var infobase string

	flag.StringVar(&infobase, "infobase", "", "Infobase name in cluster to make a backup in cli mode")

	var infobaseUser string

	flag.StringVar(&infobaseUser, "infobaseUser", "robot", "infobase admin name")

	var infobasePwd string

	flag.StringVar(&infobasePwd, "infobasePwd", "", "infobase password")

	var outputPath string

	flag.StringVar(&outputPath, "output", "", "directory backup move to")

	flag.Parse()

	// Just prepare env, config and exit
	if prepare {
		err := config.Prepare()
		if err != nil {
			log.Fatalf("Prepare error: %s", err)
		}
	}

	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg, clusterConnection, clusterName, clusterAdmin, clusterPwd,
		infobase, infobaseUser, infobasePwd, outputPath)
}
