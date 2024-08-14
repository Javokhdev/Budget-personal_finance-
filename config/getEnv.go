package config

import (
  "fmt"
  "os"
  "strings"

  "github.com/joho/godotenv"
  "github.com/spf13/cast"
)

type Config struct {
  HTTPPort string

  DB_HOST     string
  DB_PORT     int
  DB_USER     string
  DB_PASSWORD string
  DB_NAME string

  KafkaBrokers []string

  DefaultOffset string
  DefaultLimit  string

  TokenKey string
}

func Load() Config {
  if err := godotenv.Load(); err != nil {
    fmt.Println("No .env file found")
  }

  config := Config{}

  config.HTTPPort = cast.ToString(GetOrReturnDefaultValue("HTTP_PORT", ":8070"))

  config.DB_HOST = cast.ToString(GetOrReturnDefaultValue("DB_HOST", "localhost"))
  config.DB_PORT = cast.ToInt(GetOrReturnDefaultValue("DB_PORT", 5432))
  config.DB_USER = cast.ToString(GetOrReturnDefaultValue("DB_USER", "postgres"))
  config.DB_PASSWORD = cast.ToString(GetOrReturnDefaultValue("DB_PASSWORD", "root"))
  config.DB_NAME = cast.ToString(GetOrReturnDefaultValue("DB_NAME", "olympy"))

  config.KafkaBrokers = parseKafkaBrokers(GetOrReturnDefaultValue("KAFKA_BROKERS", "localhost:9092"))

  config.DefaultOffset = cast.ToString(GetOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
  config.DefaultLimit = cast.ToString(GetOrReturnDefaultValue("DEFAULT_LIMIT", "10"))
  config.TokenKey = cast.ToString(GetOrReturnDefaultValue("TokenKey", "my_secret_key"))
  return config
}

func GetOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
  val, exists := os.LookupEnv(key)

  if exists {
    return val
  }

  return defaultValue
}

func parseKafkaBrokers(brokers interface{}) []string {
  switch v := brokers.(type) {
  case string:
    return strings.Split(v, ",")
  case []string:
    return v
  default:
    return []string{"kafka:9092"}
  }
}