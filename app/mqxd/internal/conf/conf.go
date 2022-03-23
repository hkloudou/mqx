package conf

import (
	"path/filepath"

	"github.com/hkloudou/mqx/app/mqxd/conf"
	"github.com/hkloudou/mqx/app/mqxd/internal/osutil"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	log "unknwon.dev/clog/v2"
)

func init() {
	// Initialize the primary logger until logging service is up.
	err := log.NewConsole()
	if err != nil {
		panic("init console logger: " + err.Error())
	}
}

// File is the configuration object.
var File *ini.File

func Init(customConf string) error {
	data, err := conf.Files.ReadFile("app.ini")
	if err != nil {
		return errors.Wrap(err, `read default "app.ini"`)
	}

	File, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)
	if err != nil {
		return errors.Wrap(err, `parse "app.ini"`)
	}
	File.NameMapper = ini.SnackCase

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.ini")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	// CustomConf = customConf

	if osutil.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {

		log.Warn("Custom config %q not found. Ignore this warning if you're running for the first time\n", customConf)
	}

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}
	return nil
}
