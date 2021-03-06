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
	if setting.IsLocalTest() {
		host = "localhost:6379"
	} else {
		fullAddress := os.Getenv("REDISTOGO_URL")
		rawURL, err := url.Parse(fullAddress)
		if err != nil {
			logging.Error(err)
		}
		host = rawURL.Hostname() + ":" + rawURL.Port()
	}
	fmt.Println("\n ###REDIS HOST:", host)
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdelTimeout,
		Dial: func() (redis.Conn, error) {
			fmt.Println("\n### DIAL ###")
			c, err := redis.Dial("tcp", host)
			if err != nil {
				fmt.Println("\n### DIAL TCP ERROR: ", err)
				return nil, err
			}
			password := setting.RedisSetting.Paswword
			// 若不是本地 但是debug, 則輸入debug環境Redis密碼
			if !setting.IsLocalTest() && setting.ServerSetting.IsDebug() {
				password = "52b085c3ec98f1bdf6d1d9d66c5dcaec"
			} else {
				password = ""
			}
			if password != "" {
				fmt.Println("\n ##### REDIST PASSWORD: ", password)
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					fmt.Println("\n ### REDIS AUTH ERROR: ", err)
					return nil, err
				}
			}
			fmt.Println("## REDISCONNECT: ", err)
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
