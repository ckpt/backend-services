## old database structure

Tables:

  * about
  * analyser
  * bilder
  * comments
  * daarligst
  * description
  * donations
  * driver
  * epostvarsel
  * epostvarsel_tag
  * favoritt
  * fortunecookie
  * fravar
  * gjeld
  * hands
  * historie
  * locations
  * mattips
  * michelin
  * mobbing
  * musikk
  * musikkpoeng
  * news
  * oppgaver
  * passenger
  * polls
  * pollsanswer
  * praten
  * rating
  * tournament
  * tournament_results
  * varsler
  * sitat
  * stats
  * strategi
  * tipping
  * tipping_results

  * users (seperate db)

## proposed backend services

### Member service

Local Data:

  * uuid
  * nick
  * username
  * pwhash
  * (api_token)
  * avatar/profile picture
  * email
  * date of birth
  * active/retired/locked
  * other accounts?
  * debts
  * quotes
  * votes
  * description
  * gossip
  * complaints
  * notification settings

Possible Redis storage format:

  * SADD members uuid
  * HMSET member:uuid:profile username myuser pwhash 312ae47c3d1dd13 email foo@example.com
  * SET member:uuid:picture `jpeg-data`
  * SET member:uuid:nick MyNick
  * SADD member:uuid:debts member-uuid
  * HMSET member:uuid:debt:member-uuid amount value due datetime
  * SADD member:uuid:quotes qoute
  * SET member:uuid:votes `full json`
  * SET member:uuid:gossip `full json`
  * SET member:uuid:complaints `full json`
  * HMSET member:uuid:notify news false analysis true
  
Operations:

  * create new member (admin only)
  * delete member (admin only)
  * get active/retired
  * set active/retired (admin only)
  * get/set nick (set is admin)
  * set/get basic profile info
  * (api_token management)
  * (re)set password
  * get totals (played tourneys/winnings/loss, season/all time)
  * get awards (season/all time)
  * compare to
  * get debt
  * set debt (only other members)
  * get quotes
  * set quotes (only other members)
  * set votes (e.g. favourite, loser)
  * get gossip
  * set gossip (only other members)
  * get complaints
  * set complaints (only other members)
  * get/set notification settings

### Tournament service

Local Data:

  * uuid
  * scheduled datetime
  * played datetime
  * stake
  * (payouts)
  * players
  * noshows
  * results
  * location
  * catering
  * season
  * betting pool

Possible Redis storage format:

  * SADD tournaments uuid
  * SADD seasons 2015
  * SADD season:2015:tournaments uuid
  * SET tournament:uuid:scheduled 01-01-2015Z19:30
  * SET tournament:uuid:played 01-01-2015Z19:30
  * SET tournament:uuid:location location-uuid
  * SET tournament:uuid:catering catering-uuid
  * SET tournament:uuid:stake 100
  * SET tournament:uuid:season 2015
  * SADD tournament:uuid:noshows member-uuid
  * ZADD tournament:uuid:winnings 800 member-uuid
  * ZADD tournament:uuid:winnings -200 member-uuid
  * SADD tournament:uuid:bettingpools member-uuid
  * SET tournament:uuid:bettingpool:member-uuid `full json`

Operations:

  * CRUD on local data
  * R for all, CUD is admin only
  * get result
  * set result (admin only)
  * get standings (season/date range/all time, leader/placement/points/heads up etc..)
  * set betting pool entry (once per player)
  * get betting pool entries
  * get betting pool results
  
### Location service

Local Data:

  * uuid
  * host
  * url
  * gps coordinates
  * name
  * description
  * facilities
  * pictures

### Catering service

Local Data:

  * uuid
  * caterer
  * tournament
  * meal
  * votes

### Transport service

### News service

Local Data:

  * uuid
  * datetime
  * author
  * tag (e.g. `news`, `analysis`, `strategy`, `recepie`, `golden hands`)
  * title
  * leadin
  * body
  * comments
  * picture
   
### Poll service

Local Data:

  * uuid
  * author
  * datetime
  * expires
  * id of external poll?

### Statistics service

### Music service

## proposed external services

  * Picture service?
  * Spotify
  * Poll service?
  * Notifications services
