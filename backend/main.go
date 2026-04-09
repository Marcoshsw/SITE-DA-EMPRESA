package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	r := gin.Default()

	// Configuração do CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	services := []Service{
		{ID: 1, Name: "Mecânica Pesada", Description: "Manutenção completa de veículos pesados com diagnóstico e reparo especializado."},
		{ID: 2, Name: "Câmbio (Auto/Mec)", Description: "Reparo e ajuste de transmissões automáticas e mecânicas de linha pesada."},
		{ID: 3, Name: "Diferencial", Description: "Revisão e regulagem precisa do diferencial para máximo desempenho e tração."},
		{ID: 4, Name: "Motor", Description: "Retífica, manutenção e acerto completo de motores diesel."},
		{ID: 5, Name: "Automação", Description: "Programação e atualização de módulos eletrônicos para melhor performance."},
		{ID: 6, Name: "Parametrização (Motor/Câmbio)", Description: "Calibração eletrônica e ajustes técnicos para motor e câmbio."},
		{ID: 7, Name: "Arla 32", Description: "Manutenção e configuração do sistema de pós-tratamento de emissões."},
		{ID: 8, Name: "Elétrica (Geral)", Description: "Diagnóstico e reparo completo do sistema elétrico veicular."},
		{ID: 9, Name: "Rastreamento/Diagnóstico", Description: "Leitura de falhas e análise eletrônica com equipamentos de última geração."},
		{ID: 10, Name: "Common Rail", Description: "Manutenção técnica avançada em sistemas de injeção Common Rail."},
		{ID: 11, Name: "Bomba Injetora", Description: "Reparo e calibração de bombas injetoras convencionais e eletrônicas."},
		{ID: 12, Name: "Serviço de Molas", Description: "Suspensão completa, arqueamento e troca de feixes de molas."},
		{ID: 13, Name: "Borracharia", Description: "Serviços completos de pneus, alinhamento e balanceamento pesado."},
		{ID: 14, Name: "Solda", Description: "Soldagem especializada em chassis e componentes estruturais."},
		{ID: 15, Name: "Torno", Description: "Usinagem de peças com precisão milimétrica em tornearia própria."},
		{ID: 16, Name: "Unidade Injetora", Description: "Reparo e teste em bancada de unidades injetoras diesel."},
	}

	r.GET("/api/services", func(c *gin.Context) {
		c.JSON(200, services)
	})

	r.GET("/api/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Grupo Mano",
			"address": "Rua das Oficinas, 123 - Industrial",
			"phone":   "(XX) XXXX-XXXX",
			"colors":  []string{"vermelho", "branco", "cinza"},
		})
	})

	// Serve a página inicial sem criar wildcard que conflita com /api
	r.GET("/", func(c *gin.Context) {
		c.File("c:\\Users\\Marcoss\\Documents\\Site Grupo Mano\\frontend\\index.html")
	})

	r.Run(":8080")
}
