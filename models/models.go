package models

import "time"

type Usuario struct {
	Id         int `gorm:"primaryKey;autoIncrement"`
	Nome       string
	Limite     int
	Saldo      int
	Transacoes []Transacao `gorm:"foreignKey:UsuarioId"`
}

type Transacao struct {
	Id          int `gorm:"primaryKey;autoIncrement"`
	UsuarioId   int
	Valor       int
	Tipo        string `gorm:"size:1"`
	Descricao   string `gorm:"size:10"`
	RealizadaEm time.Time
}

type Tabler interface {
	TableName() string
}

func (Transacao) TableName() string {
	return "transacoes"
}
