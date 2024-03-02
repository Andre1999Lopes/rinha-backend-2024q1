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

// Regras: Uma transação de débito nunca pode deixar o saldo do cliente menor que seu limite disponível. Por exemplo, um cliente com limite de 1000 (R$ 10) nunca deverá ter o saldo menor que -1000 (R$ -10). Nesse caso, um saldo de -1001 ou menor significa inconsistência na Rinha de Backend!

// Se uma requisição para débito for deixar o saldo inconsistente, a API deve retornar HTTP Status Code 422 sem completar a transação! O corpo da resposta nesse caso não será testado e você pode escolher como o representar.

// Se o atributo [id] da URL for de uma identificação não existente de cliente, a API deve retornar HTTP Status Code 404. O corpo da resposta nesse caso não será testado e você pode escolher como o representar. Se a API retornar algo como HTTP 200 informando que o cliente não foi encontrado no corpo da resposta ou HTTP 204 sem corpo, ficarei extremamente deprimido e a Rinha será cancelada para sempre.

func Transacoes(c *gin.Context) {
	var transacao dtos.TransacaoInput
	var idParam string = c.Params.ByName("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id < 0 {
		c.JSON(422, gin.H{"mensagem": "O ID precisa ser um número inteiro positivo"})
		return
	}

	if err := c.ShouldBindJSON(&transacao); err != nil {
		c.JSON(422, gin.H{"mensagem": "Campos incorretos no corpo da requisição"})
		return
	}

	var usuario models.Usuario

	if err := database.DB.
		First(&usuario, id).
		Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"mensagem": "Usuário não encontrado"})
		return
	}

	switch transacao.Tipo {
	case "d":
		{
			usuario.Saldo -= transacao.Valor

			if usuario.Saldo < usuario.Limite*(-1) {
				c.JSON(422, gin.H{"mensagem": "Limite indisponível"})
				return
			}

			if err := database.DB.Create(&models.Transacao{
				UsuarioId:   id,
				Valor:       transacao.Valor,
				Tipo:        transacao.Tipo,
				Descricao:   transacao.Descricao,
				RealizadaEm: time.Now(),
			}).Error; err != nil {
				c.JSON(422, gin.H{"mensagem": err})
				return
			}
			database.DB.Model(&models.Usuario{Id: id}).Update("saldo", usuario.Saldo)
			var resposta = dtos.TransacaoOutput{
				Limite: usuario.Limite,
				Saldo:  usuario.Saldo,
			}
			c.JSON(200, resposta)
		}
	case "c":
		{
			usuario.Saldo += transacao.Valor
			if err := database.DB.Create(&models.Transacao{
				UsuarioId:   id,
				Valor:       transacao.Valor,
				Tipo:        transacao.Tipo,
				Descricao:   transacao.Descricao,
				RealizadaEm: time.Now(),
			}).Error; err != nil {
				c.JSON(422, gin.H{"mensagem": err.Error()})
				return
			}
			database.DB.Model(&models.Usuario{Id: id}).Update("saldo", usuario.Saldo)
			var resposta = dtos.TransacaoOutput{
				Limite: usuario.Limite,
				Saldo:  usuario.Saldo,
			}
			c.JSON(200, resposta)
		}

	default:
		{
			c.JSON(422, gin.H{"mensagem": "Valor do tipo não é válido"})
		}
	}
}
