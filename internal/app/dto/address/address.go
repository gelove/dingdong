package address

type Result struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Valid    []*Item `json:"valid_address"`
	Invalid  []*Item `json:"invalid_address"`
	MaxCount int     `json:"max_address_count"`
	CanAdd   bool    `json:"can_add_address"`
}

type Item struct {
	Id          string      `json:"id"`
	Gender      int         `json:"gender"`
	Mobile      string      `json:"mobile"`
	Location    Location    `json:"location"`
	Label       string      `json:"label"`
	UserName    string      `json:"user_name"`
	AddrDetail  string      `json:"addr_detail"`
	StationId   string      `json:"station_id"`
	StationName string      `json:"station_name"`
	IsDefault   bool        `json:"is_default"`
	CityNumber  string      `json:"city_number"`
	InfoStatus  int         `json:"info_status"`
	StationInfo StationInfo `json:"station_info"`
	VillageId   string      `json:"village_id"`
}

type Location struct {
	TypeCode string    `json:"typecode"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
	Id       string    `json:"id"`
}

type StationInfo struct {
	Id           string `json:"id"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	BusinessTime string `json:"business_time"`
	CityName     string `json:"city_name"`
	CityNumber   string `json:"city_number"`
}

type Info struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	StationId  string  `json:"station_id"`
	CityNumber string  `json:"city_number"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	UserName   string  `json:"user_name"`
	Mobile     string  `json:"mobile"`
	Address    string  `json:"address"`
	AddrDetail string  `json:"addr_detail"`
}
