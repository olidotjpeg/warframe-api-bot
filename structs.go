package main

type WarframeData struct {
	Item   string
	Chance int
	Place  string
	Rarity string
}

type Arbitration struct {
	Id          string `json:"id"`
	Activation  string `json:"activation"`
	Expiry      string `json:"expiry"`
	StartString string `json:"startString"`
	Active      bool   `json:"active"`
	Node        string `json:"node"`
	Enemy       string `json:"enemy"`
	EnemyKey    string `json:"enemyKey"`
	Type        string `json:"type"`
	TypeKey     string `json:"typeKey"`
	Archwing    bool   `json:"archwing"`
	Sharkwing   bool   `json:"sharkwing"`
}

type DarvoDeals struct {
	Sold          int    `json:"sold"`
	Item          string `json:"item"`
	Total         int    `json:"total"`
	Eta           string `json:"eta"`
	OriginalPrice int    `json:"originalPrice"`
	SalePrice     int    `json:"salePrice"`
	Discount      int    `json:"discount"`
	Expiry        string `json:"expiry"`
	Id            string `json:"id"`
}

type SortieState struct {
	Id          string
	Activation  string
	Expiry      string
	StartString string
	Active      bool
	RewardPool  string
	Variants    []Variant
	Boss        string
	Faction     string
	FactionKey  string
	Expired     bool
	Eta         string
}

type Variant struct {
	Node                string
	Boss                string
	MissionType         string
	Planet              string
	Modifier            string
	ModifierDescription string
}

type VoidItem struct {
	Item    string
	Ducat   int
	Credits int
}

type VoidTrader struct {
	Id          string
	Activation  string
	StartString string
	Expiry      string
	Active      bool
	Character   string
	Location    string
	Inventory   []VoidItem
	PsId        string
	EndString   string
}
