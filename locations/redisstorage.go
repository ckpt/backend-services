package locations

import (
	"encoding/json"
	"errors"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/m4rw3r/uuid"
	"os"
	"time"
)

type RedisLocationStorage struct {
	pool *redigo.Pool
}

func (rls *RedisLocationStorage) Store(l *Location) error {
	conn := rls.pool.Get()
	defer conn.Close()
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	if _, err = conn.Do("SADD", "locations", l.UUID); err != nil {
		return err
	}
	if _, err = conn.Do("SET", fmt.Sprintf("location:%s", l.UUID), b); err != nil {
		return err
	}
	return nil
}

func (rls *RedisLocationStorage) Load(uuid uuid.UUID) (*Location, error) {
	conn := rls.pool.Get()
	defer conn.Close()
	b, err := redigo.Bytes(conn.Do("GET", fmt.Sprintf("location:%s", uuid)))
	if err != nil {
		return nil, err
	}
	l := new(Location)
	if err := json.Unmarshal(b, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (rls *RedisLocationStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (rls *RedisLocationStorage) LoadAll() ([]*Location, error) {
	var locations []*Location
	conn := rls.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", "locations"))
	if err != nil {
		return nil, err
	}
	for _, location := range b {
		uuid, _ := uuid.FromString(location)
		l, err := rls.Load(uuid)
		if err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (rls *RedisLocationStorage) LoadByPlayer(player uuid.UUID) (*Location, error) {
	conn := rls.pool.Get()
	defer conn.Close()
	locations, err := rls.LoadAll()
	if err != nil {
		return nil, err
	}
	for _, l := range locations {
		if l.Host == player {
			return l, nil
		}
	}
	return nil, errors.New("No location found for given player")
}

func NewRedisLocationStorage() *RedisLocationStorage {
	rls := new(RedisLocationStorage)
	rls.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", os.Getenv("CKPT_REDIS"))
		},
	}
	return rls
}
