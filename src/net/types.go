package net

type RoleProfile struct {
	// AvatarUrl string // Always empty
	NickName string `json:"nickname"`
	Region   string `json:"region"`
	Level    int    `json:"level"`
}

type AvatarProfile struct {
	Id                      int    `json:"id"`
	Image                   string `json:"image"`
	Name                    string `json:"name"`
	Element                 string `json:"element"`
	Fetter                  int    `json:"fetter"`
	Level                   int    `json:"level"`
	Rarity                  int    `json:"rarity"`
	AvtivedConstellationNum int    `json:"actived_constellation_num"`
	CardImage               string `json:"card_image"`
	IsChosen                bool   `json:"is_chosen"`
}

type Statistics struct {
	ActiveDayNumber      int    `json:"active_day_number"`
	AchievementNumber    int    `json:"achievement_number"`
	AnemoculusNumber     int    `json:"anemoculus_number"`
	GeoculusNumber       int    `json:"geoculus_number"`
	AvatarNumber         int    `json:"avatar_number"`
	WayPointNumber       int    `json:"way_point_number"`
	DomainNumber         int    `json:"domain_number"`
	SpiralAbyss          string `json:"spiral_abyss"`
	PreciousChestNumber  int    `json:"precious_chest_number"`
	LuxuriousChestNumber int    `json:"luxurious_chest_number"`
	ExquisiteChestNumber int    `json:"exquisite_chest_number"`
	CommonChestNumber    int    `json:"common_chest_number"`
	ElectroculusNumber   int    `json:"electroculus_number"`
	MagicChestNumber     int    `json:"magic_chest_number"`
	DendroculusNumber    int    `json:"dendroculus_number"`
}

type OfferingProfile struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Icon  string `json:"icon"`
}

type WorldExplorationProfile struct {
	Level                 int               `json:"level"`
	ExplorationPercentage float64           `json:"exploration_percentage"`
	Icon                  string            `json:"icon"`
	Name                  string            `json:"name"`
	Type                  string            `json:"type"`
	Offerings             []OfferingProfile `json:"offerings"`
	Id                    int               `json:"id"`
	ParentId              int               `json:"parent_id"`
}

type HomeProfile struct {
	Level            int    `json:"level"`
	VisitNumber      int    `json:"visit_num"`
	ComfortNumber    int    `json:"comfort_num"`
	ItemNumber       int    `json:"item_num"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	ComfortLevelName string `json:"comfort_level_name"`
	ComfortLevelIcon string `json:"comfort_level_icon"`
}

type PlayerProfile struct {
	Role              RoleProfile               `json:"role"`
	Avatars           []AvatarProfile           `json:"avatars"`
	Statistics        Statistics                `json:"stats"`
	WorldExplorations []WorldExplorationProfile `json:"world_explorations"`
	Homes             []HomeProfile             `json:"homes"`
}

type PlayerIndexResponse struct {
	ReturnCode int           `json:"retcode"`
	Message    string        `json:"message"`
	Data       PlayerProfile `json:"data"`
}
