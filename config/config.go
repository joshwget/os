package config

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	yaml "github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/rancher/os/util"
)

func (c *CloudConfig) Import(bytes []byte) (*CloudConfig, error) {
	data, err := readConfig(bytes, false, CloudConfigPrivateFile)
	if err != nil {
		return c, err
	}

	return NewConfig().Merge(data)
}

func (c *CloudConfig) MergeBytes(bytes []byte) (*CloudConfig, error) {
	data, err := readConfig(bytes, false)
	if err != nil {
		return c, err
	}
	return c.Merge(data)
}

var keysToStringify = []string{
	"command",
	"dns",
	"dns_search",
	"entrypoint",
	"env_file",
	"environment",
	"labels",
	"links",
}

func isPathToStringify(path []interface{}) bool {
	l := len(path)
	if l == 0 {
		return false
	}
	if sk, ok := path[l-1].(string); ok {
		return util.Contains(keysToStringify, sk)
	}
	return false
}

func stringifyValue(data interface{}, path []interface{}) interface{} {
	switch data := data.(type) {
	case map[interface{}]interface{}:
		result := make(map[interface{}]interface{}, len(data))
		if isPathToStringify(path) {
			for k, v := range data {
				switch v := v.(type) {
				case []interface{}:
					result[k] = stringifyValue(v, append(path, k))
				case map[interface{}]interface{}:
					result[k] = stringifyValue(v, append(path, k))
				default:
					result[k] = fmt.Sprint(v)
				}
			}
		} else {
			for k, v := range data {
				result[k] = stringifyValue(v, append(path, k))
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(data))
		if isPathToStringify(path) {
			for k, v := range data {
				result[k] = fmt.Sprint(v)
			}
		} else {
			for k, v := range data {
				result[k] = stringifyValue(v, append(path, k))
			}
		}
		return result
	default:
		return data
	}
}

func StringifyValues(data map[interface{}]interface{}) map[interface{}]interface{} {
	return stringifyValue(data, nil).(map[interface{}]interface{})
}

func (c *CloudConfig) Merge(values map[interface{}]interface{}) (*CloudConfig, error) {
	d := map[interface{}]interface{}{}
	if err := util.Convert(c, &d); err != nil {
		return c, err
	}
	r := util.MapsUnion(d, StringifyValues(values))
	t := &CloudConfig{}
	if err := util.Convert(r, t); err != nil {
		return c, err
	}
	return t, nil
}

func Dump(private, full bool) (string, error) {
	var cfg *CloudConfig
	var err error

	if full {
		cfg, err = LoadConfig()
	} else {
		files := []string{CloudConfigBootFile, CloudConfigPrivateFile, CloudConfigFile}
		if !private {
			files = util.FilterStrings(files, func(x string) bool { return x != CloudConfigPrivateFile })
		}
		cfg, err = ChainCfgFuncs(nil,
			func(_ *CloudConfig) (*CloudConfig, error) { return ReadConfig(nil, true, files...) },
			amendNils,
		)
	}

	if err != nil {
		return "", err
	}

	bytes, err := yaml.Marshal(*cfg)
	return string(bytes), err
}

func (c *CloudConfig) Get(key string) (interface{}, error) {
	data := map[interface{}]interface{}{}
	if err := util.Convert(c, &data); err != nil {
		return nil, err
	}

	v, _ := getOrSetVal(key, data, nil)
	return v, nil
}

func (c *CloudConfig) GetIgnoreOmitEmpty(key string) (interface{}, error) {
	data := map[interface{}]interface{}{}
	if err := util.ConvertIgnoreOmitEmpty(c, &data); err != nil {
		return nil, err
	}

	v, _ := getOrSetVal(key, data, nil)
	return v, nil
}

func (c *CloudConfig) Set(key string, value interface{}) (map[interface{}]interface{}, error) {
	data := map[interface{}]interface{}{}
	_, data = getOrSetVal(key, data, value)

	return data, nil
}

func (c *CloudConfig) Save(cfgDiffs ...map[interface{}]interface{}) error {
	files := append([]string{OsConfigFile, OemConfigFile}, CloudConfigDirFiles()...)
	files = util.FilterStrings(files, func(x string) bool { return x != CloudConfigPrivateFile })
	exCfg, err := ChainCfgFuncs(nil,
		func(_ *CloudConfig) (*CloudConfig, error) {
			return ReadConfig(nil, true, files...)
		},
		readCmdline,
		amendNils)
	if err != nil {
		return err
	}

	exCfg = mergeMetadata(exCfg, readMetadata())

	exData := map[interface{}]interface{}{}
	if err := util.Convert(exCfg, &exData); err != nil {
		return err
	}

	data := map[interface{}]interface{}{}
	if err := util.Convert(c, &data); err != nil {
		return err
	}

	data = util.MapsDifference(data, exData)

	// Apply any additional config diffs
	for _, diff := range cfgDiffs {
		data = util.MapsUnion(data, diff)
	}

	log.WithFields(log.Fields{"diff": data}).Debug("The diff we're about to save")
	if err := saveToDisk(data); err != nil {
		return err
	}
	return nil
}
