package config

import "github.com/spf13/viper"

func GetDbHost() string {
	return viper.GetString("database.host")
}

func GetDbName() string {
	return viper.GetString("database.dbname")
}

func GetDbUser() string {
	return viper.GetString("database.user")
}

func GetDbPassword() string {
	return viper.GetString("database.password")
}

func GetDbPort() string {
	return viper.GetString("database.port")
}

func GetJwtSecret() string {
	return viper.GetString("jwt.secret")
}

func GetJwtTTL() int64 {
	if viper.GetInt64("jwt.ttl") == 0 {
		return DefaultJwtTTL
	} else {
		return viper.GetInt64("jwt.ttl")
	}
}

func GetRedisHost() string {
	return viper.GetString("redis.host")
}

func GetRedisDB() int {
	return viper.GetInt("redis.db")
}
