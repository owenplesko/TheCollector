package types

type Match struct {
	Id          string    `json:"id"`
	Date        int64     `json:"Date"`
	DataVersion string    `json:"DataVersion"`
	Info        MatchInfo `json:"Info"`
}
