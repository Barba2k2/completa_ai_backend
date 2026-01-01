package models

type Section struct {
	Code   string    `json:"code"`
	Name   string    `json:"name"`
	Total  int       `json:"total"`
	Type   string    `json:"type"`
	Items  []Sticker `json:"items,omitempty"`
}

type Sticker struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	SectionCode string `json:"section_code"`
	IsShiny     bool   `json:"is_shiny"`
}

var AlbumSections = []Section{
	{Code: "FWC", Name: "FIFA World Cup", Total: 10, Type: "special"},
	{Code: "STD", Name: "Estádios", Total: 16, Type: "special"},
	{Code: "CTY", Name: "Cidades-Sede", Total: 16, Type: "special"},
	{Code: "LEG", Name: "Legends", Total: 10, Type: "special"},
	{Code: "USA", Name: "Estados Unidos", Total: 20, Type: "team"},
	{Code: "MEX", Name: "México", Total: 20, Type: "team"},
	{Code: "CAN", Name: "Canadá", Total: 20, Type: "team"},
	{Code: "ARG", Name: "Argentina", Total: 20, Type: "team"},
	{Code: "BRA", Name: "Brasil", Total: 20, Type: "team"},
	{Code: "URU", Name: "Uruguai", Total: 20, Type: "team"},
	{Code: "COL", Name: "Colômbia", Total: 20, Type: "team"},
	{Code: "ECU", Name: "Equador", Total: 20, Type: "team"},
	{Code: "CHI", Name: "Chile", Total: 20, Type: "team"},
	{Code: "PAR", Name: "Paraguai", Total: 20, Type: "team"},
	{Code: "PER", Name: "Peru", Total: 20, Type: "team"},
	{Code: "VEN", Name: "Venezuela", Total: 20, Type: "team"},
	{Code: "ENG", Name: "Inglaterra", Total: 20, Type: "team"},
	{Code: "FRA", Name: "França", Total: 20, Type: "team"},
	{Code: "GER", Name: "Alemanha", Total: 20, Type: "team"},
	{Code: "ESP", Name: "Espanha", Total: 20, Type: "team"},
	{Code: "ITA", Name: "Itália", Total: 20, Type: "team"},
	{Code: "POR", Name: "Portugal", Total: 20, Type: "team"},
	{Code: "NED", Name: "Holanda", Total: 20, Type: "team"},
	{Code: "BEL", Name: "Bélgica", Total: 20, Type: "team"},
	{Code: "CRO", Name: "Croácia", Total: 20, Type: "team"},
	{Code: "SRB", Name: "Sérvia", Total: 20, Type: "team"},
	{Code: "SUI", Name: "Suíça", Total: 20, Type: "team"},
	{Code: "POL", Name: "Polônia", Total: 20, Type: "team"},
	{Code: "DEN", Name: "Dinamarca", Total: 20, Type: "team"},
	{Code: "AUT", Name: "Áustria", Total: 20, Type: "team"},
	{Code: "WAL", Name: "País de Gales", Total: 20, Type: "team"},
	{Code: "SCO", Name: "Escócia", Total: 20, Type: "team"},
	{Code: "UKR", Name: "Ucrânia", Total: 20, Type: "team"},
	{Code: "JPN", Name: "Japão", Total: 20, Type: "team"},
	{Code: "KOR", Name: "Coreia do Sul", Total: 20, Type: "team"},
	{Code: "AUS", Name: "Austrália", Total: 20, Type: "team"},
	{Code: "IRN", Name: "Irã", Total: 20, Type: "team"},
	{Code: "KSA", Name: "Arábia Saudita", Total: 20, Type: "team"},
	{Code: "QAT", Name: "Catar", Total: 20, Type: "team"},
	{Code: "MAR", Name: "Marrocos", Total: 20, Type: "team"},
	{Code: "SEN", Name: "Senegal", Total: 20, Type: "team"},
	{Code: "GHA", Name: "Gana", Total: 20, Type: "team"},
	{Code: "NGA", Name: "Nigéria", Total: 20, Type: "team"},
	{Code: "CMR", Name: "Camarões", Total: 20, Type: "team"},
	{Code: "TUN", Name: "Tunísia", Total: 20, Type: "team"},
	{Code: "EGY", Name: "Egito", Total: 20, Type: "team"},
	{Code: "ALG", Name: "Argélia", Total: 20, Type: "team"},
	{Code: "CIV", Name: "Costa do Marfim", Total: 20, Type: "team"},
	{Code: "RSA", Name: "África do Sul", Total: 20, Type: "team"},
	{Code: "CRC", Name: "Costa Rica", Total: 20, Type: "team"},
	{Code: "PAN", Name: "Panamá", Total: 20, Type: "team"},
	{Code: "JAM", Name: "Jamaica", Total: 20, Type: "team"},
}

func GetTotalStickers() int {
	total := 0
	for _, section := range AlbumSections {
		total += section.Total
	}
	return total
}

func GetSectionByCode(code string) *Section {
	for _, section := range AlbumSections {
		if section.Code == code {
			return &section
		}
	}
	return nil
}
