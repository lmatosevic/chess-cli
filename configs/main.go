package configs

import (
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type general struct {
	AppName     string `yaml:"appName"`
	Version     string
	Author      string
	Description string
}

type server struct {
	Hostname string
	Host     string
	Port     uint16
	Debug    bool
}

type database struct {
	Host          string
	Port          uint16
	Name          string
	Schema        string
	Username      string
	Password      string
	MigrationsDir string `yaml:"migrationsDir"`
	AutoMigrate   bool   `yaml:"autoMigrate"`
}

type rules struct {
	DefaultTurnDurationSeconds int32 `yaml:"defaultTurnDurationSeconds"`
	DrawRequestTimeoutTurns    int32 `yaml:"drawRequestTimeoutTurns"`
	MaxCreatedGames            int32 `yaml:"maxCreatedGames"`
	MaxJoinedGames             int32 `yaml:"maxJoinedGames"`
}

type Config struct {
	General  general
	Server   server
	Database database
	Rules    rules
}

const defaultConfigPath = "./config.yaml"

var config Config

func GetConfig() *Config {
	// If config already parsed, return existing value
	if config != (Config{}) {
		return &config
	}

	// Load external or default config file
	var tempConfig Config
	if len(os.Args) > 1 {
		configPath := os.Args[1]
		tempConfig = *GetConfigFromFile(configPath)
		log.Printf("Config file loaded from %q\n", configPath)
	} else {
		tempConfig = *GetConfigFromFile(defaultConfigPath)
		log.Printf("Using default config file\n")
	}

	OverrideWithEnvVariables(&tempConfig, "")

	config = tempConfig

	return &config
}

func OverrideWithEnvVariables(obj any, path string) {
	pathSep := ""
	if path != "" {
		pathSep = "_"
	}

	// Resolve reflected values based on object type
	values := reflect.ValueOf(obj)
	if values.Kind() == reflect.Ptr {
		values = reflect.Indirect(values)
	} else if values.Kind() == reflect.Struct {
		values = reflect.Indirect(obj.(reflect.Value))
	} else if values.Kind() == reflect.Interface {
		values = values.Elem()
	}
	typ := values.Type()

	// Iterate over object fields and call this function recursively for struct types and try to update primitive types
	// from environment variables if they are set (e.g. Database.Password -> DATABASE_PASSWORD)
	for i := 0; i < values.NumField(); i++ {
		fieldName := typ.Field(i).Name
		envVarName := fmt.Sprintf("%s%s%s", path, pathSep, fieldName)
		if typ.Field(i).Type.Kind() == reflect.Struct {
			inf := values.Field(i).Addr()
			OverrideWithEnvVariables(inf, envVarName)
		} else {
			env, ok := os.LookupEnv(strings.ToUpper(envVarName))
			if ok {
				f := values.Field(i)
				if f.IsValid() && f.CanSet() {
					switch f.Kind() {
					case reflect.String:
						f.SetString(env)
					case reflect.Bool:
						f.SetBool(strings.ToLower(env) == "true")
					case reflect.Int:
					case reflect.Int16:
					case reflect.Int32:
					case reflect.Int64:
					case reflect.Uint:
					case reflect.Uint16:
					case reflect.Uint32:
					case reflect.Uint64:
						val, e := strconv.Atoi(env)
						if e == nil {
							f.SetInt(int64(val))
						}
					case reflect.Float32:
						val, e := strconv.ParseFloat(env, 32)
						if e == nil {
							f.SetInt(int64(val))
						}
					case reflect.Float64:
						val, e := strconv.ParseFloat(env, 64)
						if e == nil {
							f.SetInt(int64(val))
						}
					default:
						log.Printf("Unsupported env variable type: %s -> %s\n", strings.ToUpper(envVarName), f.Kind())
					}
				}
			}
		}
	}
}

func GetConfigFromFile(filePath string) *Config {
	conf, err := readConfig(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	return conf
}

func readConfig(filePath string) (*Config, error) {
	buf, err := utils.ReadFromFile(filePath)

	conf := &Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("config error in file %q: %v", filePath, err)
	}

	return conf, nil
}
