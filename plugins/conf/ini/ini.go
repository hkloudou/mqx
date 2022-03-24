package ini

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hkloudou/mqx/face"
	"gopkg.in/ini.v1"
)

type iniStruct struct {
	_file *ini.File
}

func New(custom string) (face.Conf, error) {
	obj := &iniStruct{}
	if err := obj.Init(custom); err != nil {
		return nil, err
	}
	return obj, nil
}

func MustNew(custom string) face.Conf {
	obj, err := New(custom)
	if err != nil {
		panic(err)
	}
	return obj
}

// File is the configuration object.
// var File *ini.File

func (m *iniStruct) Init(customConf string) error {
	data, err := Files.ReadFile("app.ini")
	if err != nil {
		return fmt.Errorf(`%v read default "app.ini"`, err)
	}

	m._file, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
		Insensitive:         true,
		AllowShadows:        true,
	}, data)
	if err != nil {
		return fmt.Errorf(`%v parse "app.ini"`, err.Error())
	}
	m._file.NameMapper = ini.SnackCase
	m._file.ValueMapper = os.ExpandEnv
	if customConf == "" {
		customConf = filepath.Join("", "conf", "app.ini")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return fmt.Errorf("%v get absolute path", err)
		}
	}
	// CustomConf = customConf

	if IsFile(customConf) {
		if err = m._file.Append(customConf); err != nil {
			return fmt.Errorf("%v append %q", err, customConf)
		}
	} else {
		log.Printf("Custom config %q not found. Ignore this warning if you're running for the first time\n", customConf)
	}

	// if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
	// 	return errors.Wrap(err, "mapping default section")
	// }
	return nil
}

func (m *iniStruct) MapTo(section string, source interface{}) error {
	_section := section
	if _section == "" {
		_section = ini.DefaultSection
	}
	if err := m._file.Section(_section).MapTo(source); err != nil {
		return fmt.Errorf("%v mapping ini config", err)
	}
	return nil
}

// IsFile returns true if given path exists as a file (i.e. not a directory).
func IsFile(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsDir returns true if given path is a directory, and returns false when it's
// a file or does not exist.
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}