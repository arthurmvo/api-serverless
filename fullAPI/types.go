package main

type Drop struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Link        string  `json:"link"`
	IsNightmare bool    `json:"isNightmare"`
	NpcPrice    float64 `json:"npcPrice"`
	Chance      float64 `json:"chance"`
	IsRare      bool    `json:"isRare"`
}

type Clan struct {
	ID      string   `json:"id"`
	Clan    string   `json:"clan"`
	Icon    string   `json:"icon"`
	Element string   `json:"element"`
	Lurer   string   `json:"lurer"`
	Hunts   []string `json:"hunts"`
}

type Hunt struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	Desc         string `json:"desc"`
	WhoHunts     []int  `json:"whoHunts"`
	Level        int    `json:"level"`
	Pokemons     []int  `json:"pokemons"`
	IsNightmare  bool   `json:"isNightmare"`
	IsAngryShiny bool   `json:"isAngryShiny"`
	AngryShiny   int    `json:"angryShiny"`
	Boss         []int  `json:"boss"`
}

type Pokemon struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Link  string `json:"link"`
	Drops []int  `json:"drops"`
}

type Player struct {
	UID        int `json:"uid"`
	Characters []struct {
		Name string `json:"name"`
		Clan string `json:"clan"`
	} `json:"characters"`
	Transactions struct {
		Hunts []struct {
			Type        string          `json:"type"`
			HuntID      int             `json:"hunt_id"`
			Description string          `json:"description"`
			Character   string          `json:"character"`
			Time        int             `json:"time"`
			Market      float64         `json:"market"`
			Result      float64         `json:"result"`
			Drops       [][]interface{} `json:"drops"`
			Mobs        []int           `json:"mobs"`
			Extra       [][]interface{} `json:"extra"`
			Date        string          `json:"date"`
			Id          string          `json:"id"`
		} `json:"hunts"`
		Services []struct {
			Type        string  `json:"type"`
			Category    string  `json:"category"`
			Character   string  `json:"character"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Result      float64 `json:"result"`
			Date        string  `json:"date"`
			Id          string  `json:"id"`
		} `json:"services"`
		Purchases []struct {
			Type        string  `json:"type"`
			Category    string  `json:"category"`
			Character   string  `json:"character"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Result      float64 `json:"result"`
			Date        string  `json:"date"`
			Id          string  `json:"id"`
		} `json:"purchases"`
		Rifts []struct {
			Type string `json:"type"`
		} `json:"rifts"`
		Boss []struct {
			Type      string          `json:"type"`
			Result    float64         `json:"result"`
			Date      string          `json:"date"`
			Id        int             `json:"id"`
			Bosses    []string        `json:"bosses"`
			Drops     [][]interface{} `json:"drops"`
			Character string          `json:"character"`
		} `json:"boss"`
	} `json:"transactions"`
}
