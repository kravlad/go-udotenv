package udotenv

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const defaultEnvPath = ".env"
const (
	envsId = iota + 1
	overloadId
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Config represents the configuration settings for managing environment variables.
// It provides options for specifying environment flags, overload flags, a default
// environment file path, and whether to overload by default.
//
// Fields:
//   - EnvFlags: A list of environment variable flags to be used.
//   - OverloadFlags: A list of flags that determine whether environment variables
//     should be overloaded.
//   - DefaultEnvPath: The default file path to the environment file.
//   - OverloadByDefault: A boolean indicating whether environment variables should
//     be overloaded by default.
type Config struct {
	EnvFlags          []string
	OverloadFlags     []string
	DefaultEnvPath    string
	OverloadByDefault bool
}

// udotEnvType represents the environment configuration structure for the application.
// It contains the configuration settings, an environment parameter, and a flag
// to determine whether to overload existing parameters.
//
// Fields:
// - Config: A pointer to the Config structure that holds the application's configuration settings.
// - EnvParam: A string representing the environment parameter to be used.
// - OverloadParam: A boolean flag indicating whether to overwrite existing environment parameters.
type udotEnvType struct {
	Config        *Config
	EnvParam      stringSlice
	OverloadParam bool
}

// Load reads environment variables from a specified file and loads them into
// the application's environment. If the `OverloadParam` field is set to true,
// it will overwrite existing environment variables with the values from the file.
//
// The method uses the `godotenv` package to handle the loading process. If the
// `EnvParam` field is empty, the method returns immediately without performing
// any action. If an error occurs while loading the file, the method will panic
// with an error message.
//
// Note: Ensure that `EnvParam` is set to the path of the environment file before
// calling this method.
//
// Example:
//
//	ue := &udotEnv{
//	    EnvParam:      []string{".env"},
//	    OverloadParam: false,
//	}
//	ue.Load() // Loads environment variables from the .env file.
func (ue *udotEnvType) Load() {
	if len(ue.EnvParam) == 0 {
		return
	}

	f := godotenv.Load
	if ue.OverloadParam == true {
		f = godotenv.Overload
	}

	err := f(ue.EnvParam...)
	if err != nil {
		panic(fmt.Sprintln("error loading file '", ue.EnvParam, "'"))
	}
}

// GetDefaultConfig returns a pointer to a Config struct initialized with
// default values. The default configuration includes predefined flags for
// environment variables and overload options, as well as a default path
// for the environment file.
func GetDefaultConfig() *Config {
	return &Config{
		EnvFlags:       []string{"envs", "e"},
		OverloadFlags:  []string{"env-overload", "eo", "o"},
		DefaultEnvPath: defaultEnvPath,
	}
}

// New creates and initializes a new instance of udotEnvType with the provided configuration.
//
// Parameters:
//   - parseFlags: A boolean indicating whether to parse command-line flags immediately.
//   - config: Optional variadic parameter to pass a single *Config instance. If no configuration
//     is provided, a default configuration will be used. If more than one configuration is passed,
//     the function will panic.
//
// Behavior:
//   - If no configuration is provided, the default configuration is used.
//   - If a configuration is provided, it is used to initialize the udotEnvType instance. If the
//     DefaultEnvPath in the configuration is empty, it is set to a predefined default value.
//   - Command-line flags are registered based on the EnvFlags and OverloadFlags in the configuration.
//     Flags are stored in a map to ensure that only one flag per parameter is passed.
//   - If the `parseFlags` parameter is true, the function will parse the command-line flags.
//
// Panics:
//   - If more than one configuration is passed.
//   - If multiple flags for the same parameter are passed.
//
// Returns:
//   - A pointer to the initialized udotEnvType instance.
func New(parseFlags bool, config ...*Config) (udotEnv *udotEnvType) {
	udotEnv = &udotEnvType{}
	if len(config) == 0 {
		udotEnv.Config = GetDefaultConfig()

	} else if len(config) == 1 {
		udotEnv.Config = config[0]
		if udotEnv.Config.DefaultEnvPath == "" {
			udotEnv.Config.DefaultEnvPath = defaultEnvPath
		}

	} else {
		panic("only 1 config must be passed")
	}

	flagStorage := make(map[string]int, len(udotEnv.Config.EnvFlags)+len(udotEnv.Config.OverloadFlags))
	for _, v := range udotEnv.Config.EnvFlags {
		flag.Var(&udotEnv.EnvParam, v, "help message for flag n")
		flagStorage[v] = envsId
	}

	for _, v := range udotEnv.Config.OverloadFlags {
		flag.BoolVar(&udotEnv.OverloadParam, v, udotEnv.Config.OverloadByDefault, "help message for flag n")
		flagStorage[v] = overloadId
	}

	if len(os.Args) <= 1 {
		return
	}

	newArgs := make([]string, 1, len(os.Args)+1) // add 1 for case if envParam passed without a value
	newArgs[0] = os.Args[0]

	passedParams := make(map[int]bool, 2)
	for i, argName := range os.Args[1:] {
		newArgs = append(newArgs, argName)
		if !strings.HasPrefix(argName, "-") || len(argName) < 2 {
			continue
		}

		argId, ok := flagStorage[argName[1:]]
		if !ok {
			argId, ok = flagStorage[argName[2:]]
		}
		_, passed := passedParams[argId]

		if ok && passed {
			panic("only one flag per param must be passed")
		} else if ok && argId != envsId {
			passedParams[argId] = true
		}

		if (argId == envsId) &&
			((len(os.Args)-2 == i) ||
				((len(os.Args)-2 > i) && (strings.HasPrefix(os.Args[i+2], "-")))) {
			newArgs = append(newArgs, udotEnv.Config.DefaultEnvPath)
		}
	}
	os.Args = newArgs

	if parseFlags {
		flag.Parse()
	}
	return
}
