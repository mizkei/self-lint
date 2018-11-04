package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mizkei/self-lint/lint"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var (
		configPath  string
		includeTest bool
	)
	flag.StringVar(&configPath, "config", "", "required: config filepath")
	flag.BoolVar(&includeTest, "include-test", false, "optional: include test files")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return errors.New("argument not enough")
	}
	target := flag.Args()[0]

	if configPath == "" {
		flag.Usage()
		return errors.New("argument not enough")
	}

	dir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "get working directory failed")
	}

	conf, err := loadConf(configPath)
	if err != nil {
		return err
	}

	trifles, err := lint.Run(conf, target, includeTest)
	if err != nil {
		return errors.Wrap(err, "run lint failed")
	}

	for _, tf := range trifles {
		p, err := filepath.Rel(dir, tf.Position.Filename)
		if err != nil {
			return errors.Wrap(err, "get relative filepath failed")
		}
		fmt.Printf("%s:%d:%d: %s\n", p, tf.Position.Line, tf.Position.Column, tf.Text)
	}

	return nil
}

func loadConf(f string) (lint.Config, error) {
	cfile, err := os.Open(f)
	if err != nil {
		return lint.Config{}, errors.Wrap(err, "open config file failed")
	}
	defer cfile.Close()

	conf, err := lint.LoadConfig(cfile)
	if err != nil {
		return lint.Config{}, errors.Wrap(err, "load config failed")
	}

	return conf, nil
}
