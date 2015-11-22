package news

import (
	"errors"
	"os"
	"github.com/imdario/mergo"
	"github.com/m4rw3r/uuid"
	"github.com/ckpt/backend-services/utils"
	"time"
)

// We use dummy in memory storage for now
var storage NewsItemStorage = NewRedisNewsItemStorage()

// Init a message queue
var eventqueue utils.AMQPQueue = utils.NewRMQ(os.Getenv("CKPT_AMQP_URL"), "ckpt.events")

type NewsItem struct {
	UUID     uuid.UUID `json:"uuid"`
	Author   uuid.UUID `json:"author"`
	Created  time.Time `json:"created"`
	Tag      Tag       `json:"tag"`
	Title    string    `json:"title"`
	Leadin   string    `json:"leadin"`
	Body     string    `json:"body"`
	Picture  []byte    `json:"picture"`
	Comments []Comment `json:"comments"`
}

type Tag int

const (
	Article Tag = iota
	Analysis
	Strategy
	Recepie
	GoldenHand
)

type Comment struct {
	Player  uuid.UUID `json:player`
	Content string    `json:content`
}

// A storage interface for News
type NewsItemStorage interface {
	Store(*NewsItem) error
	Delete(uuid.UUID) error
	Load(uuid.UUID) (*NewsItem, error)
	LoadAll() ([]*NewsItem, error)
	LoadByAuthor(uuid.UUID) ([]*NewsItem, error)
}

//
// NewsItem related functions and methods
//

// Create a NewsItem
func NewNewsItem(itemdata NewsItem, author uuid.UUID) (*NewsItem, error) {
	c := new(NewsItem)
	if err := mergo.MergeWithOverwrite(c, itemdata); err != nil {
		return nil, errors.New(err.Error() + " - Could not set initial NewsItem data")
	}
	c.UUID, _ = uuid.V4()
	c.Author = author
	c.Created = time.Now()
	if err := storage.Store(c); err != nil {
		return nil, errors.New(err.Error() + " - Could not write NewsItem to storage")
	}
	eventqueue.Publish(utils.CKPTEvent{
		Type: utils.NEWS_EVENT,
		Subject: "Nytt bidrag lagt ut",
		Message: "Det er lagt ut et nytt bidrag på ckpt.no!",})
	return c, nil
}

func AllNewsItems() ([]*NewsItem, error) {
	return storage.LoadAll()
}

func DeleteByUUID(uuid uuid.UUID) bool {
	err := storage.Delete(uuid)
	if err != nil {
		return false
	}
	return true
}

func NewsItemByUUID(uuid uuid.UUID) (*NewsItem, error) {
	return storage.Load(uuid)
}

func (c *NewsItem) UpdateNewsItem(ci NewsItem) error {
	d := new(NewsItem)
	*d = *c
	if err := mergo.MergeWithOverwrite(c, ci); err != nil {
		return errors.New(err.Error() + " - Could not update NewsItem info")
	}
	c.UUID = d.UUID
	c.Author = d.Author
	c.Created = d.Created
	c.Tag = d.Tag
	err := storage.Store(c)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated NewsItem info")
	}
	return nil
}

// TODO: Comments
// func (c *NewsItem) AddVote(player uuid.UUID, score int) error {
// 	vote := Vote{Player: player, Score: score}
// 	c.Votes = append(c.Votes, vote)
// 	err := storage.Store(c)
// 	if err != nil {
// 		return errors.New(err.Error() + " - Could not store updated NewsItem info with added vote")
// 	}
// 	return nil
// }

// func (c *NewsItem) RemoveVote(player uuid.UUID) error {
// 	for i, v := range c.Votes {
// 		if v.Player == player {
// 			c.Votes = append(c.Votes[:i], c.Votes[i+1:]...)
// 		}
// 	}
// 	err := storage.Store(c)
// 	if err != nil {
// 		return errors.New(err.Error() + " - Could not store updated NewsItem info with removed vote")
// 	}
// 	return nil
// }
