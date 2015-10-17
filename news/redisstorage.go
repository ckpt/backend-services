package news

import (
	"encoding/json"
	"errors"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/m4rw3r/uuid"
	"os"
	"time"
)

type RedisNewsItemStorage struct {
	pool *redigo.Pool
}

func (rnis *RedisNewsItemStorage) Store(c *NewsItem) error {
	conn := rnis.pool.Get()
	defer conn.Close()
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if _, err = conn.Do("SADD", "newsitems", c.UUID); err != nil {
		return err
	}
	if _, err = conn.Do("SET", fmt.Sprintf("newsitem:%s", c.UUID), b); err != nil {
		return err
	}
	return nil
}

func (rnis *RedisNewsItemStorage) Load(uuid uuid.UUID) (*NewsItem, error) {
	conn := rnis.pool.Get()
	defer conn.Close()
	b, err := redigo.Bytes(conn.Do("GET", fmt.Sprintf("newsitem:%s", uuid)))
	if err != nil {
		return nil, err
	}
	c := new(NewsItem)
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (rnis *RedisNewsItemStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (rnis *RedisNewsItemStorage) LoadAll() ([]*NewsItem, error) {
	var newsitems []*NewsItem
	conn := rnis.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", "newsitems"))
	if err != nil {
		return nil, err
	}
	for _, newsitem := range b {
		uuid, _ := uuid.FromString(newsitem)
		c, err := rnis.Load(uuid)
		if err != nil {
			return nil, err
		}
		newsitems = append(newsitems, c)
	}
	return newsitems, nil
}

func (rnis *RedisNewsItemStorage) LoadByAuthor(author uuid.UUID) ([]*NewsItem, error) {
	conn := rnis.pool.Get()
	defer conn.Close()
	found := make([]*NewsItem, 0)
	newsitems, err := rnis.LoadAll()
	if err != nil {
		return nil, err
	}
	for _, c := range newsitems {
		if c.Author == author {
			found = append(found, c)
		}
	}
	return found, nil
}

func NewRedisNewsItemStorage() *RedisNewsItemStorage {
	rnis := new(RedisNewsItemStorage)
	rnis.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", os.Getenv("CKPT_REDIS"))
		},
	}
	return rnis
}
