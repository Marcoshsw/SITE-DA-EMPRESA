package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func applyStaticCacheHeaders(c *gin.Context, targetFile string) {
	if strings.HasSuffix(strings.ToLower(targetFile), ".html") {
		c.Header("Cache-Control", "no-cache")
		return
	}

	if strings.TrimSpace(c.Query("v")) != "" {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
}

func resolveFrontendDir() (string, error) {
	wd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	candidates := []string{
		filepath.Join(wd, "frontend"),
		filepath.Join(wd, "..", "frontend"),
		filepath.Join(exeDir, "frontend"),
		filepath.Join(exeDir, "..", "frontend"),
	}

	for _, dir := range candidates {
		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			abs, absErr := filepath.Abs(dir)
			if absErr == nil {
				return abs, nil
			}
			return dir, nil
		}
	}

	return "", fmt.Errorf("nao foi possivel localizar frontend/index.html")
}

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	frontendDir, err := resolveFrontendDir()
	if err != nil {
		panic(err)
	}

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

	// Serve a página inicial e arquivos estáticos do frontend sem depender de caminho absoluto.
	r.GET("/", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.File(filepath.Join(frontendDir, "index.html"))
	})

	r.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/api") {
			c.AbortWithStatus(404)
			return
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(requestPath, "/"))
		if cleanPath == "." {
			c.Header("Cache-Control", "no-cache")
			c.File(filepath.Join(frontendDir, "index.html"))
			return
		}

		targetFile := filepath.Join(frontendDir, cleanPath)
		targetAbs, absErr := filepath.Abs(targetFile)
		if absErr != nil || (!strings.HasPrefix(targetAbs, frontendDir+string(os.PathSeparator)) && targetAbs != frontendDir) {
			c.AbortWithStatus(404)
			return
		}

		if info, statErr := os.Stat(targetFile); statErr == nil && !info.IsDir() {
			applyStaticCacheHeaders(c, targetFile)
			c.File(targetFile)
			return
		}

		c.Header("Cache-Control", "no-cache")
		c.File(filepath.Join(frontendDir, "index.html"))
	})

	r.Run(":8080")
}
