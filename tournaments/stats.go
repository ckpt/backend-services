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

func BestPlayer(tournaments Tournaments) uuid.UUID {
	standings := NewStandings(tournaments)
	standings.ByAvgPlace()

	// TODO: More than 2 ppl could share best AvgPlace
	if standings[0].AvgPlace == standings[1].AvgPlace {
		sort.Sort(standings[0].Results)
		sort.Sort(standings[1].Results)

		placesA := standings[0].Results
		placesB := standings[1].Results

		for i := 0; i < standings[0].NumTotal; i++ {
			if len(placesA) >= i+1 && len(placesB) >= i+1 {
				if placesA[i].Place < placesB[i].Place {
					return standings[0].Player
				}
				if placesB[i].Place < placesA[i].Place {
					return standings[1].Player
				}
			} else {
				if len(placesA) < len(placesB) {
					return standings[0].Player
				} else {
					return standings[1].Player
				}
			}
		}
	}
	return standings[0].Player
}

func WorstPlayer(tournaments Tournaments) uuid.UUID {
	standings := NewStandings(tournaments)
	standings.ByAvgPlace()

	m, n := len(standings)-1, len(standings)-2
	// TODO: More than 2 ppl could share worst AvgPlace
	if standings[m].AvgPlace == standings[n].AvgPlace {
		println("Found two players with same worst avg:")
		println("Player ", standings[m].Player.String(), "with ", standings[m].AvgPlace)
		println("Player ", standings[n].Player.String(), "with ", standings[n].AvgPlace)
		sort.Sort(sort.Reverse(standings[m].Results))
		sort.Sort(sort.Reverse(standings[n].Results))

		placesA := standings[m].Results
		placesB := standings[n].Results

		for i := 0; i < standings[m].NumTotal; i++ {
			if len(placesA) >= i+1 && len(placesB) >= i+1 {
				if placesA[i].Place > placesB[i].Place {
					return standings[m].Player
				}
				if placesB[i].Place > placesA[i].Place {
					return standings[n].Player
				}
			} else {
				if len(placesA) < len(placesB) {
					return standings[m].Player
				} else {
					return standings[n].Player
				}
			}
		}
	}
	return standings[m].Player
}

func YellowPeriods(tournaments Tournaments) []YellowPeriod {
	var periods []YellowPeriod
	var currentPeriod *YellowPeriod
	var season, seasonIndex int

	sort.Sort(tournaments)
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

func mostYellowDaysInSeason(yellowPeriods []YellowPeriod, season int) (uuid.UUID, int) {
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

func loserOfTheYear(monthPeriods []*MonthStats, season int) (uuid.UUID, int) {
	lastPlacesByPlayer := make(map[uuid.UUID]int)

	for _, mp := range monthPeriods {
		if mp.Year != season {
			continue
		}
		lastPlacesByPlayer[mp.Worst] += 1
	}

	max := 0
	var maxPlayer uuid.UUID
	for player, count := range lastPlacesByPlayer {
		if count > max {
			max = count
			maxPlayer = player
		}
	}

	return maxPlayer, max
}

func playerOfTheYear(monthPeriods []*MonthStats, season int) (uuid.UUID, int) {
	topPlacesByPlayer := make(map[uuid.UUID]int)

	for _, mp := range monthPeriods {
		if mp.Year != season {
			continue
		}
		topPlacesByPlayer[mp.Best] += 1
	}

	max := 0
	var maxPlayer uuid.UUID
	for player, count := range topPlacesByPlayer {
		if count > max {
			max = count
			maxPlayer = player
		}
	}

	return maxPlayer, max
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

		// FIXME: Need enough tournaments to get titles..
		seasonStandings.ByAvgPlace()
		for i := range seasonStandings {
			if seasonStandings[i].Enough {
				p, ap := seasonStandings[i].Player, seasonStandings[i].AvgPlace
				titles.AvgPlaceWinner.Uuid = p
				titles.AvgPlaceWinner.AvgPlace = ap
				break
			}
		}

		seasonStandings.ByPoints()
		for i := range seasonStandings {
			if seasonStandings[i].Enough {
				p, pnts := seasonStandings[i].Player, seasonStandings[i].Points
				titles.PointsWinner.Uuid = p
				titles.PointsWinner.Points = pnts
				break
			}
		}

		// FIXME: These are missing tie breaks..
		p, d := mostYellowDaysInSeason(seasonStats.YellowPeriods, season)
		titles.MostYellowDays.Uuid = p
		titles.MostYellowDays.Days = d

		p, c := playerOfTheYear(seasonStats.MonthStats, season)
		titles.PlayerOfTheYear.Uuid = p
		titles.PlayerOfTheYear.Months = c

		p, c = loserOfTheYear(seasonStats.MonthStats, season)
		titles.LoserOfTheYear.Uuid = p
		titles.LoserOfTheYear.Months = c

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
		for k, v := range byMonth {
			monthStats := new(MonthStats)

			played := v.Played()

			if len(played) == 0 {
				continue
			}

			sort.Sort(v)
			monthStats.Year = season
			monthStats.Month = k

			monthStats.Best = BestPlayer(v)
			println("Generating worst player for month", k, "in year", season)
			println("    Based on ", len(v), "tournaments")
			monthStats.Worst = WorstPlayer(v)
			println("    Worst player:", monthStats.Worst.String())

			stats.MonthStats = append(stats.MonthStats, monthStats)
		}
	}
	return stats
}
