package tournaments

import (
	"github.com/m4rw3r/uuid"
	"sort"
	"time"
)

type YellowPeriod struct {
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Player uuid.UUID `json:"uuid"`
	Active bool      `json:"active"`
}

type MonthStats struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
	Best  uuid.UUID  `json:"best"`
	Worst uuid.UUID  `json:"worst"`
}

type PeriodStats struct {
	YellowPeriods []YellowPeriod `json:"yellowPeriods"`
	MonthStats    []*MonthStats  `json:"monthStats"`
}

type SeasonTitles struct {
	Season   int `json:"season"`
	Champion struct {
		Uuid     uuid.UUID `json:"uuid"`
		Winnings int       `json:"winnings"`
	} `json:"champion"`
	AvgPlaceWinner struct {
		Uuid     uuid.UUID `json:"uuid"`
		AvgPlace float64   `json:"avgPlace"`
	} `json:"avgPlaceWinner"`
	PointsWinner struct {
		Uuid   uuid.UUID `json:"uuid"`
		Points int       `json:"points"`
	} `json:"pointsWinner"`
	MostYellowDays struct {
		Uuid uuid.UUID `json:"uuid"`
		Days int       `json:"days"`
	} `json:"mostYellowDays"`
	PlayerOfTheYear struct {
		Uuid   uuid.UUID `json:"uuid"`
		Months int       `json:"months"`
	} `json:"playerOfTheYear"`
	LoserOfTheYear struct {
		Uuid   uuid.UUID `json:"uuid"`
		Months int       `json:"months"`
	} `json:"loserOfTheYear"`
}

func BestPlayer(tournaments Tournaments) (PlayerStandings, bool) {
	standings := NewStandings(tournaments)
	standings.ByBestPlayer()

	if standings[0].Results.Equals(standings[1].Results) {
		// It's still a tie
		return standings, true
	}
	return standings, false
}

func WorstPlayer(tournaments Tournaments) (PlayerStandings, bool) {
	standings := NewStandings(tournaments)
	standings.ByWorstPlayer()
	if standings[0].Results.Equals(standings[1].Results) {
		// It's still a tie
		return standings, true
	}
	return standings, false
}

func YellowPeriods(tournaments Tournaments) []YellowPeriod {
	var periods []YellowPeriod
	var currentPeriod *YellowPeriod
	var season, seasonIndex int

	sort.Stable(tournaments)
	for i := range tournaments {
		if !tournaments[i].Played {
			continue
		}
		// Leader is based on results for the season, so start from "scratch" on new seasons
		if tournaments[i].Info.Season != season {
			season = tournaments[i].Info.Season
			seasonIndex = i
		}
		standings := NewStandings(tournaments[seasonIndex : i+1])
		standings.ByWinnings(season < 2013)
		if currentPeriod == nil {
			currentPeriod = &YellowPeriod{
				From:   tournaments[i].Info.Scheduled,
				To:     tournaments[i].Info.Scheduled,
				Player: standings[0].Player,
				Active: true,
			}
		} else if currentPeriod.Player == standings[0].Player {
			currentPeriod.To = tournaments[i].Info.Scheduled
		} else {
			currentPeriod.Active = false
			currentPeriod.To = tournaments[i].Info.Scheduled
			periods = append(periods, *currentPeriod)
			currentPeriod = &YellowPeriod{
				From:   tournaments[i].Info.Scheduled,
				To:     tournaments[i].Info.Scheduled,
				Player: standings[0].Player,
				Active: true,
			}
		}
	}
	periods = append(periods, *currentPeriod)
	return periods
}

func YellowDaysInSeason(yellowPeriods []YellowPeriod, season int) map[uuid.UUID]int {
	seasonStartDate := time.Date(season, 1, 1, 0, 0, 0, 0, time.Local)
	seasonEndDate := seasonStartDate.AddDate(1, 0, 0)

	daysByPlayer := make(map[uuid.UUID]int)
	for _, p := range yellowPeriods {

		// Skip entries that don't concert the given season
		if p.To.Before(seasonStartDate) || p.From.After(seasonEndDate) {
			continue
		}

		// If entry starts last season, adjust it so it starts at season start date
		if p.From.Before(seasonStartDate) {
			p.From = seasonStartDate
		}

		// If entry ends next season or is active, adjust it so it ends at season end date
		if p.To.After(seasonEndDate) || p.Active {
			p.To = seasonEndDate
		}

		days := int(p.To.Sub(p.From).Hours() / 24)
		daysByPlayer[p.Player] += days
	}
	return daysByPlayer
}

func mostYellowDaysInSeason(yellowPeriods []YellowPeriod, season int) (uuid.UUID, int) {
	daysByPlayer := YellowDaysInSeason(yellowPeriods, season)
	max := 0
	var maxPlayer uuid.UUID
	for p, d := range daysByPlayer {
		if d > max {
			max = d
			maxPlayer = p
		}
	}

	return maxPlayer, max
}

func loserOfTheYear(monthPeriods []*MonthStats, season int) ([]uuid.UUID, int) {
	lastPlacesByPlayer := make(map[uuid.UUID]int)
	playersByPlace := make(map[int][]uuid.UUID)

	for _, mp := range monthPeriods {
		if mp.Year != season {
			continue
		}
		lastPlacesByPlayer[mp.Worst] += 1
	}

	max := 0
	for player, count := range lastPlacesByPlayer {
		playersByPlace[count] = append(playersByPlace[count], player)
		if count > max {
			max = count
		}
	}

	return playersByPlace[max], max
}

func playerOfTheYear(monthPeriods []*MonthStats, season int) ([]uuid.UUID, int) {
	topPlacesByPlayer := make(map[uuid.UUID]int)
	playersByPlace := make(map[int][]uuid.UUID)

	for _, mp := range monthPeriods {
		if mp.Year != season {
			continue
		}
		topPlacesByPlayer[mp.Best] += 1
	}

	max := 0
	for player, count := range topPlacesByPlayer {
		playersByPlace[count] = append(playersByPlace[count], player)
		if count > max {
			max = count
		}
	}

	return playersByPlace[max], max
}

func Titles(seasons []int) []*SeasonTitles {

	var titleList []*SeasonTitles

	seasonStats := SeasonStats(seasons)
	for _, season := range seasons {
		titles := &SeasonTitles{Season: season}
		t, _ := TournamentsBySeason(season)
		seasonStandings := NewStandings(t)

		seasonStandings.ByWinnings(season < 2013)

		for i := range seasonStandings {
			if seasonStandings[i].Enough {
				p, w := seasonStandings[i].Player, seasonStandings[i].Winnings
				titles.Champion.Uuid = p
				titles.Champion.Winnings = w
				break
			}
		}

		seasonStandings.ByAvgPlace()
		for i := range seasonStandings {
			if seasonStandings[i].Enough {
				p, ap := seasonStandings[i].Player, seasonStandings[i].AvgPlace
				titles.AvgPlaceWinner.Uuid = p
				titles.AvgPlaceWinner.AvgPlace = ap
				break
			}
		}

		players, c := playerOfTheYear(seasonStats.MonthStats, season)
		if len(players) == 1 {
			titles.PlayerOfTheYear.Uuid = players[0]
		} else {
		POTYLoop:
			for _, p := range players {
				for _, s := range seasonStandings {
					if s.Player == p && s.Enough {
						titles.PlayerOfTheYear.Uuid = p
						break POTYLoop
					}
				}
			}
		}
		titles.PlayerOfTheYear.Months = c

		players, c = loserOfTheYear(seasonStats.MonthStats, season)
		if len(players) == 1 {
			titles.LoserOfTheYear.Uuid = players[0]
		} else {
		LOTYLoop:
			for _, p := range players {
				for i := len(seasonStandings) - 1; i >= 0; i-- {
					if seasonStandings[i].Player == p && seasonStandings[i].Enough {
						titles.LoserOfTheYear.Uuid = p
						break LOTYLoop
					}
				}
			}
		}

		titles.LoserOfTheYear.Months = c

		seasonStandings.ByPoints()
		for i := range seasonStandings {
			if seasonStandings[i].Enough {
				p, pnts := seasonStandings[i].Player, seasonStandings[i].Points
				titles.PointsWinner.Uuid = p
				titles.PointsWinner.Points = pnts
				break
			}
		}

		// FIXME: This is missing tie breaks..
		p, d := mostYellowDaysInSeason(seasonStats.YellowPeriods, season)
		titles.MostYellowDays.Uuid = p
		titles.MostYellowDays.Days = d

		titleList = append(titleList, titles)
	}
	return titleList
}

func SeasonStats(seasons []int) *PeriodStats {
	stats := new(PeriodStats)
	all, err := AllTournaments()
	if err != nil {
		// TODO
	}

	var t Tournaments
	for _, event := range all {
		for _, s := range seasons {
			if event.Info.Season == s {
				t = append(t, event)
			}
		}
	}

	yellows := YellowPeriods(t)
	stats.YellowPeriods = yellows

	for _, season := range seasons {
		byMonth := t.GroupByMonths(season)
		var sortedMonths []int
		for k := range byMonth {
			sortedMonths = append(sortedMonths, int(k))
		}
		sort.Ints(sortedMonths)
		for _,i := range sortedMonths {
			v := byMonth[time.Month(i)]
			monthStats := new(MonthStats)

			played := v.Played()

			if len(played) == 0 {
				continue
			}

			sort.Stable(v)
			monthStats.Year = season
			monthStats.Month = time.Month(i)

			best, tie := BestPlayer(v)
			bestplayer := best[0].Player
			if tie {
				var tiebreakTournaments Tournaments
				for j := 1; j <= i; j++ {
					tiebreakTournaments = append(tiebreakTournaments, byMonth[time.Month(j)]...)
				}
				tiebest, tie := BestPlayer(tiebreakTournaments)
				if tie {
					println("    Warning: Tied for best player for month", i, "in year", season)
				}
				for p := range tiebest {
					if tiebest[p].Player == best[0].Player || tiebest[p].Player == best[1].Player {
						bestplayer = tiebest[p].Player
						break
					}
				}
			}
			monthStats.Best = bestplayer

			worst, tie := WorstPlayer(v)
			worstplayer := worst[0].Player
			if tie {
				var tiebreakTournaments Tournaments
				for j := 1; j <= int(i); j++ {
					tiebreakTournaments = append(tiebreakTournaments, byMonth[time.Month(j)]...)
				}

				tieworst, tie := WorstPlayer(tiebreakTournaments)
				if tie {
					println("    Warning: Tied for worst player for month", i, "in year", season)
				}
				for p := range tieworst {
					if tieworst[p].Player == worst[0].Player || tieworst[p].Player == worst[1].Player {
						worstplayer = tieworst[p].Player
						break
					}
				}
			}
			monthStats.Worst = worstplayer

			stats.MonthStats = append(stats.MonthStats, monthStats)
		}
	}
	return stats
}
