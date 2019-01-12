package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

func Get(key string) interface{} {
	return viper.Get(key)
}
func GetBool(key string) bool {
	return viper.GetBool(key)
}
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
func GetInt(key string) int {
	return viper.GetInt(key)
}
func GetU64(key string) uint64 {
	return uint64(viper.GetInt(key))
}
func GetString(key string) string {
	return viper.GetString(key)
}
func GetMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
func GetStringMap(key string) map[string]string {
	return viper.GetStringMapString(key)
}
func GetListMap(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}

func MustGet(key string) interface{} {
	return panicIfEmpty(key, viper.Get(key))
}
func MustGetDuration(key string) time.Duration {
	return panicIfEmpty(key, viper.GetDuration(key)).(time.Duration)
}
func MustGetInt(key string) int {
	return panicIfEmpty(key, viper.GetInt(key)).(int)
}
func MustGetU64(key string) uint64 {
	return uint64(MustGetInt(key))
}
func MustGetString(key string) string {
	return panicIfEmpty(key, viper.GetString(key)).(string)
}
func MustGetMap(key string) map[string]interface{} {
	return panicIfLenZero(key, viper.GetStringMap(key)).(map[string]interface{})
}
func MustGetStringMap(key string) map[string]string {
	return panicIfLenZero(key, viper.GetStringMapString(key)).(map[string]string)
}
func MustGetListMap(key string) map[string][]string {
	return panicIfLenZero(key, viper.GetStringMapStringSlice(key)).(map[string][]string)
}

func panicIfEmpty(key string, val interface{}) interface{} {
	if val == reflect.Zero(reflect.TypeOf(val)).Interface() {
		panic(fmt.Sprintf("configuration key %s is not defined", key))
	}
	return val
}

func panicIfLenZero(key string, val interface{}) interface{} {
	if val == nil || reflect.ValueOf(val).Len() == 0 {
		panic(fmt.Sprintf("configuration key %s is not defined or maps to empty value", key))
	}
	return val
}
