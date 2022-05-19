package environment

import (
	"fmt"
	"os"
	"strconv"
)

func GetValue(name string) string {
	return os.Getenv(name)
}

func GetValueOrDefault(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

func MustGet(name string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	err := fmt.Errorf("\"%s\" environment variable is not defined", name)
	panic(err)
}

func GetBool(name string) (bool, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		return false, err
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}

func GetBoolOrDefault(name string, fallback bool) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolValue
}

func MustGetBool(name string) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		panic(err)
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}
	return boolValue
}

func GetDouble(name string) (float64, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		return 0, err
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return floatValue, nil
}

func GetDoubleOrDefault(name string, fallback float64) float64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return floatValue
}

func MustGetDouble(name string) float64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		panic(err)
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return floatValue
}

func GetInt(name string) (int64, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		return 0, err
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func GetIntOrDefault(name string, fallback int64) int64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return intValue
}

func MustGetInt(name string) int64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		err := fmt.Errorf("\"%s\" environment variable is not defined", name)
		panic(err)
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return intValue
}
