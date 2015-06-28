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

### Player service

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

  * SADD players uuid
  * SET player:uuid:user username
  * SADD users username
  * HMSET user:username username myuser pwhash 312ae47c3d1dd13 admin false apitoken xxxxxxxx
  * HMSET user:username:notify news false analysis true
  * HMSET player:uuid:profile email foo@example.com name "Morten Knutsen" description "The boss"
  * SET player:uuid:picture `jpeg-data`
  * SET player:uuid:nick MyNick
  * SADD player:uuid:debts debt-uuid
  * HMSET player:uuid:debt:debt-uuid who player-uuid amount value due datetime
  * SADD player:uuid:quotes qoute
  * SET player:uuid:votes `full json`
  * SET player:uuid:gossip `full json`
  * SET player:uuid:complaints `full json`
 
Operations:

  What                                 API endpoint

  * create new player (admin only)  -> POST /players
  * get all players                 -> GET /players
  * delete player (admin only)      -> DELETE /players/:uuid
  * get player info                 -> GET /players/:uuid
  * set player nick (admin only)    -> PUT /players/:uuid/
  * set player active status (admin)-> PUT /players/:uuid/
  * set user for player (admin)     -> PUT /players/:uuid/user
  * get basic profile info          -> GET /players/:uuid/profile
  * set basic profile info          -> PUT /players/:uuid/profile
  * get totals                      -> GET /players/:uuid/winnings
    (played tourneys/winnings/loss,
     season/all time)
  * get awards                      -> GET /players/:uuid/awards
    (season/all time)
  * compare to                      -> GET /players/:uuid/compare/:otheruuid
  * get debt                        -> GET /players/:uuid/debts
  * ceate debt (other players)      -> POST /players/:uuid/debts
  * delete debt (other players)     -> DELETE /players/:uuid/debts/:debtuuid
  * get quotes                      -> GET /players/quotes
                                    -> GET /players/:uuid/quotes
  * set quotes (only other players) -> POST /players/:uuid/quotes
  * set votes                       -> PUT /players/:uuid/votes
    (e.g. favourite, loser)
  * get gossip                      -> GET /players/gossip
                                    -> GET /players/:uuid/gossip
  * set gossip (only other players) -> PUT /players/:uuid/gossip
  * set complaints                  -> POST /players/:uuid/complaints
  * get complaints                  -> GET /players/:uuid/complaints
  * get complaints by other player  -> GET /players/:uuid/complaints/:playeruuid

  * get tasks                       -> GET /players/tasks
                                    -> GET /players/:uuid/tasks
  * create user account             -> POST /users
  * remove user account             -> DELETE /users/:username
  * lock user                       -> PUT /users/:username
  * (re)set password?               -> PUT /users/:username
  * get notification settings       -> GET /users/:username/settings
  * set notification settings       -> PUT /users/:username/settings
  * get userinfo                    -> GET /users/:username


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
  * HMSET tournament:uuid:noshows player-uuid "Taktisk fravær" player-uuid2 "Fisketur"
  * ZADD tournament:uuid:winnings 800 player-uuid
  * ZADD tournament:uuid:winnings -200 player-uuid
  * SADD tournament:uuid:bettingpools player-uuid
  * SET tournament:uuid:bettingpool:player-uuid `full json`

Operations:

  What                                 API endpoint

  * create new season               -> POST /seasons
  * get seasons                     -> GET /seasons
  * get season tournaments?         -> GET /season/:year/tournaments
  * get season standings            -> GET /season/:year/standings
  * get tournaments                 -> GET /tournaments
  * create new tournament           -> POST /tournaments
  * delete tournament               -> DELETE /tournaments/:uuid
  * get basic tournament info       -> GET /tournaments/:uuid
  * update basic tournament info    -> PUT /tournaments/:uuid
  * set tournament result           -> PUT /tournaments/:uuid/result
  * get tournament result           -> GET /tournaments/:uuid/result
  * get tournament players          -> GET /tournaments/:uuid/players
  * add tournament noshow           -> POST /tournaments/:uuid/noshows
  * set noshow reason               -> PUT /tournaments/:uuid/noshows/:playeruuid
  * remove tournament noshow        -> DELETE /tournaments/:uuid/noshows/:playeruuid
  * get tournament noshows          -> GET /tournaments/:uuid/noshows

  * set betting pool entry          -> PUT /tournaments/:uuid/bettingpool/:playeruuid
  * delete betting pool entry       -> DELETE /tournaments/:uuid/bettingpool/:playeruuid
  * get betting pool entries        -> GET /tournaments/:uuid/bettingpool
  * get betting pool results        -> GET /tournaments/:uuid/bettingpool/results

  * get standings                   -> GET /tournaments/standings
    (season/date range/all time, leader/placement/points/heads up etc..)
  
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

Possible Redis storage format:

  * SADD locations uuid
  * HMSET location:uuid:basicinfo host player-uuid url http://gulesider.no/Xy97V3 name Heimdal description "Arenaen stod ferdig ..."
  * SADD location:uuid:pictures `jpeg-data`

Operations:

  What                                 API endpoint

  * add new location                -> POST /locations
  * update location info            -> PUT /locations/:uuid
  * remove location                 -> DELETE /locations/:uuid

### Catering service

Local Data:

  * uuid
  * caterer
  * tournament
  * meal
  * votes

Possible Redis storage format:

  * SADD caterings uuid
  * HMSET catering:uuid:basicinfo caterer player-uuid tournament tournament-uuid meal "Kokte grillpølser"
  * SADD catering:uuid:votes player-uuid    ZADD?
  * SET catering:uuid:votes:player-uuid 4    ^^?

Operations:

  What                                 API endpoint

  * add new catering                -> POST /caterings
  * update catering info            -> PUT /caterings/:uuid
  * remove catering                 -> DELETE /caterings/:uuid
  * vote on catering                -> POST /caterings/:uuid/votes
  * update vote                     -> PUT /caterings/:uuid/votes/:playeruuid
  * remove vote                     -> DELETE /caterings/:uuid/votes/:playeruuid
  
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
