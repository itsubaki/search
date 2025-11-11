package es

type Index struct {
	Settings Settings `json:"settings"`
	Mappings Mappings `json:"mappings"`
}

type Settings struct {
	Analysis Analysis `json:"analysis"`
}

type Analysis struct {
	CharFilter map[string]CharFilter `json:"char_filter"`
	Tokenizer  map[string]Tokenizer  `json:"tokenizer"`
	Filter     map[string]Filter     `json:"filter"`
	Analyzer   map[string]Analyzer   `json:"analyzer"`
}

type CharFilter struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type Tokenizer struct {
	Type                 string   `json:"type"`
	Mode                 string   `json:"mode,omitempty"`
	DiscardCompoundToken bool     `json:"discard_compound_token,omitempty"`
	UserDictionaryRules  []string `json:"user_dictionary_rules,omitempty"`
	MinGram              int      `json:"min_gram,omitempty"`
	MaxGram              int      `json:"max_gram,omitempty"`
	TokenChars           []string `json:"token_chars,omitempty"`
}

type Filter struct {
	Type     string   `json:"type"`
	Lenient  bool     `json:"lenient"`
	Synonyms []string `json:"synonyms"`
}

type Analyzer struct {
	Type       string   `json:"type"`
	CharFilter []string `json:"char_filter"`
	Tokenizer  string   `json:"tokenizer"`
	Filter     []string `json:"filter"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type           string           `json:"type"`
	SearchAnalyzer string           `json:"search_analyzer"`
	Analyzer       string           `json:"analyzer"`
	Fields         map[string]Field `json:"fields"`
}

type Field struct {
	Type           string `json:"type"`
	SearchAnalyzer string `json:"search_analyzer"`
	Analyzer       string `json:"analyzer"`
}
