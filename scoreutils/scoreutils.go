package scoreutils

type Score struct {
	Score      int64
	Discipline int
	Name       string
}

var TopScores Score

func (s Score) GetDiscipline() string {
	switch s.Discipline {
	case 0:
		return ".22lr Prone"
	case 1:
		return ".22lr Kneeling"
	case 2:
		return ".22lr Offhand Carbine"
	}
	return "misc"
}
