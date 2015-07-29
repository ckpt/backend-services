package caterings

import (
	"encoding/json"
	"errors"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/m4rw3r/uuid"
	"os"
	"time"
)

type RedisCateringStorage struct {
	pool *redigo.Pool
}

func (rcs *RedisCateringStorage) Store(c *Catering) error {
	conn := rcs.pool.Get()
	defer conn.Close()
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if _, err = conn.Do("SADD", "caterings", c.UUID); err != nil {
		return err
	}
	if _, err = conn.Do("SET", fmt.Sprintf("catering:%s", c.UUID), b); err != nil {
		return err
	}
	return nil
}

func (rcs *RedisCateringStorage) Load(uuid uuid.UUID) (*Catering, error) {
	conn := rcs.pool.Get()
	defer conn.Close()
	b, err := redigo.Bytes(conn.Do("GET", fmt.Sprintf("catering:%s", uuid)))
	if err != nil {
		return nil, err
	}
	c := new(Catering)
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (rcs *RedisCateringStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (rcs *RedisCateringStorage) LoadAll() ([]*Catering, error) {
	var caterings []*Catering
	conn := rcs.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", "caterings"))
	if err != nil {
		return nil, err
	}
	for _, catering := range b {
		uuid, _ := uuid.FromString(catering)
		c, err := rcs.Load(uuid)
		if err != nil {
			return nil, err
		}
		caterings = append(caterings, c)
	}
	return caterings, nil
}

func (rcs *RedisCateringStorage) LoadByTournament(tournament uuid.UUID) (*Catering, error) {
	conn := rcs.pool.Get()
	defer conn.Close()
	caterings, err := rcs.LoadAll()
	if err != nil {
		return nil, err
	}
	for _, c := range caterings {
		if c.Tournament == tournament {
			return c, nil
		}
	}
	return nil, errors.New("No catering found for given tournament")
}

func NewRedisCateringStorage() *RedisCateringStorage {
	rcs := new(RedisCateringStorage)
	rcs.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", os.Getenv("CKPT_REDIS"))
		},
	}
	return rcs
}
