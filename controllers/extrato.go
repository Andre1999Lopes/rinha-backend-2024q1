package controllers

import (
	"errors"
	"rinha/database"
	"rinha/dtos"
	"rinha/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Extrato(c *gin.Context) {
	var usuario models.Usuario
	var ultimasTransacoes = []dtos.UltimasTransacoes{}
	var idParam = c.Params.ByName("id")
	var id, err = strconv.Atoi(idParam)

	if err != nil || id < 0 {
		c.JSON(422, gin.H{"message": "O ID deve ser um número inteiro positivo"})
		return
	}

	if err := database.DB.
		Model(&models.Usuario{}).
		Preload("Transacoes", func(db *gorm.DB) *gorm.DB {
			return db.
				Model(&models.Transacao{}).
				Where("transacoes.usuario_id = ?", id).
				Order("transacoes.realizada_em DESC").
				Limit(10)
		}).
		First(&usuario, id).
		Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"message": "Usuário não encontrado"})
		return
	}

	for _, transacao := range usuario.Transacoes {
		ultimasTransacoes = append(ultimasTransacoes, dtos.UltimasTransacoes{
			Valor:       transacao.Valor,
			Descricao:   transacao.Descricao,
			Tipo:        transacao.Tipo,
			RealizadaEm: transacao.RealizadaEm,
		})
	}

	var extrato = &dtos.ExtratoOutput{
		Saldo: dtos.Saldo{
			Total:       usuario.Saldo,
			DataExtrato: time.Now(),
			Limite:      usuario.Limite,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	c.JSON(200, extrato)
}
