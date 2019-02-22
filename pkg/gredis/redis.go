package gredis

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/setting"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

func Setup() error {
	host := ""
	if setting.ServerSetting.RunMode == "debug" {
		host = "localhost:6379"
	} else {
		fullAddress := os.Getenv("REDISTOGO_URL")
		rawURL, err := url.Parse(fullAddress)
		if err != nil {
			logging.Error(err)
		}
		host = rawURL.Hostname()
	}
	fmt.Println("\n ###REDIS HOST:", host)
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdelTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}

			if setting.RedisSetting.Paswword != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Paswword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		fmt.Println("ERROR: ", err)
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	// redis.Bytes(reply interface{}, err error)：将命令返回转为 Bytes
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, err
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	// redis.Bool(reply interface{}, err error)：将命令返回转为布尔值
	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()
	// redis.Strings(reply interface{}, err error)：将命令返回转为 []string
	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
