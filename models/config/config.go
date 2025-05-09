package config

import "encoding/json"

func UnmarshalConfig(data []byte) (Config, error) {
	var r Config
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Config) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Config struct {
	Lojas []Loja `json:"lojas"`
}

type Loja struct {
	NomeFantasia              string      `json:"nomeFantasia"`
	Licenca                   string      `json:"licenca"`
	Cnpj                      string      `json:"cnpj"`
	NumeroLoja                int         `json:"numeroLoja"`
	CodigoSistema             float64     `json:"codigoSistema"`
	Automacao                 string      `json:"automacao"`
	Localizacao               Localizacao `json:"localizacao"`
	Horarios                  Horarios    `json:"horarios"`
	BancoFrenteLojaPrimario   BancoArio   `json:"bancoFrenteLojaPrimario"`
	BancoFrenteLojaSecundario BancoArio   `json:"bancoFrenteLojaSecundario"`
	BancoRetaguardaPrimario   BancoArio   `json:"bancoRetaguardaPrimario"`
	BancoRetaguardaSecundario BancoArio   `json:"bancoRetaguardaSecundario"`
}

type BancoArio struct {
	Host      string `json:"host"`
	NomeBanco string `json:"nomeBanco"`
	Usuario   string `json:"usuario"`
	Senha     string `json:"senha"`
	Porta     int    `json:"porta"`
}

type Horarios struct {
	Uteis   DiasUteis `json:"uteis"`
	Domingo Domingo   `json:"domingo"`
}

type DiasUteis struct {
	Inicio string `json:"inicio"`
	Fim    string `json:"fim"`
}

type Domingo struct {
	Inicio string `json:"inicio"`
	Fim    string `json:"fim"`
}

type Localizacao struct {
	Estado       string `json:"estado"`
	CodigoEstado int64  `json:"codigoEstado"`
	Cidade       string `json:"cidade"`
	CodigoCidade int64  `json:"codigoCidade"`
	Bairro       string `json:"bairro"`
	TimeZone     string `json:"timeZone"`
}
