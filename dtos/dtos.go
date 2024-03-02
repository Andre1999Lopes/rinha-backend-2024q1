package dtos

import "time"

type TransacaoInput struct {
	Valor     int    `json:"valor" binding:"required"`
	Tipo      string `json:"tipo" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
}

type TransacaoOutput struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type ExtratoOutput struct {
	Saldo             Saldo               `json:"saldo"`
	UltimasTransacoes []UltimasTransacoes `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type UltimasTransacoes struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}
