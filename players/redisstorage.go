package players

import (
	"encoding/json"
	"errors"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/m4rw3r/uuid"
	"os"
	"time"
)

type RedisPlayerStorage struct {
	pool *redigo.Pool
}

func (rps *RedisPlayerStorage) Store(p *Player) error {
	conn := rps.pool.Get()
	defer conn.Close()
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = conn.Do("SADD", "players", p.UUID)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", fmt.Sprintf("player:%s", p.UUID), b)
	if err != nil {
		return err
	}
	if p.User.Username != "" {
		_, err = conn.Do("SADD", "users", p.User.Username)
		if err != nil {
			return err
		}
		_, err = conn.Do("SET", fmt.Sprintf("user:%s:pwhash", p.User.Username),
			p.User.password)
		if err != nil {
			return err
		}
		_, err = conn.Do("SET", fmt.Sprintf("user:%s:player", p.User.Username),
			p.UUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rps *RedisPlayerStorage) Load(uuid uuid.UUID) (*Player, error) {
	conn := rps.pool.Get()
	defer conn.Close()
	b, err := redigo.Bytes(conn.Do("GET", fmt.Sprintf("player:%s", uuid)))
	if err != nil {
		return nil, err
	}
	p := new(Player)
	err = json.Unmarshal(b, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (rps *RedisPlayerStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (rps *RedisPlayerStorage) LoadAll() ([]*Player, error) {
	var players []*Player
	conn := rps.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", "players"))
	if err != nil {
		return nil, err
	}
	for _, player := range b {
		uuid, _ := uuid.FromString(player)
		p, err := rps.Load(uuid)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (rps *RedisPlayerStorage) LoadUser(username string) (*User, error) {
	conn := rps.pool.Get()
	defer conn.Close()
	player, err := redigo.String(conn.Do("GET", fmt.Sprintf("user:%s:player", username)))
	if err != nil {
		return nil, err
	}
	uuid, _ := uuid.FromString(player)
	p, err := rps.Load(uuid)
	if err != nil {
		return nil, err
	}
	return &p.User, nil
}

func NewRedisPlayerStorage() *RedisPlayerStorage {
	rps := new(RedisPlayerStorage)
	rps.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", os.Getenv("CKPT_REDIS"))
		},
	}
	return rps
}

func OLDcreateUUIDs(number int) []uuid.UUID {
	var uuids []uuid.UUID
	for number > 0 {
		uuid, _ := uuid.V4()
		uuids = append(uuids, uuid)
		number--
	}
	return uuids
}
