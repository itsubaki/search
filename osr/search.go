package osr

type SearchResult[T any] struct {
	Took     int         `json:"took"`
	TimedOut bool        `json:"timed_out"`
	Shards   Shards      `json:"_shards"`
	Hits     HitsBody[T] `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsBody[T any] struct {
	Total    TotalInfo      `json:"total"`
	MaxScore float64        `json:"max_score"`
	Hits     []HitDetail[T] `json:"hits"`
}

type TotalInfo struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitDetail[T any] struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}
