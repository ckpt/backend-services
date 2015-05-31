package main

import (
	"errors"
	"time"
	"github.com/m4rw3r/uuid"
)

func createUUIDs(number int) []uuid.UUID {
	var uuids []uuid.UUID
	for number > 0 {
		uuid, _ := uuid.V4()
		uuids = append(uuids, uuid)
		number--
	}
	return uuids
}

func getMembers() []Member {
	return dummyMembers[:]
}

func getMember(memberID uuid.UUID) (Member, error) {
	for _, member := range dummyMembers {
		if member.UUID == memberID {
			return member, nil
		}
	}
	return Member{}, errors.New("Not found")
}


var dummyUUIDs = createUUIDs(5)
var dummyMembers = [...]Member{
	Member{
		UUID: dummyUUIDs[0],
		Profile: MemberProfile{
			Birthday: time.Date(1979, time.April, 14, 0, 0, 0, 0, time.Local),
			Name: "Morten Knutsen",
			Email: "knumor@gmail.com",
		},
		Nick: "Panzer",
		Active: true,
	},
	Member{
		UUID: dummyUUIDs[1],
		Profile: MemberProfile{
			Birthday: time.Date(1979, time.October, 20, 0, 0, 0, 0, time.Local),
			Name: "Bjørn Helge Kopperud",
			Email: "bjorn@kjekkegutter.no",
		},
		Nick: "Bjøro",
		Active: true,
	},
}
