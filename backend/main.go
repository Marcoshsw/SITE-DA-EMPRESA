package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Service string `json:"service"`
	Message string `json:"message"`
}

type appConfig struct {
	Port              string
	CorsOrigins       []string
	FrontendDir       string
	GinMode           string
	AllowWildcardCORS bool
}

func loadAppConfig() appConfig {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	ginMode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}

	frontendDir := strings.TrimSpace(os.Getenv("FRONTEND_DIR"))
	if frontendDir == "" {
		resolvedDir, err := resolveFrontendDir()
		if err != nil {
			panic(err)
		}
		frontendDir = resolvedDir
	} else if absDir, err := filepath.Abs(frontendDir); err == nil {
		frontendDir = absDir
	}

	originsRaw := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	origins := []string{"*"}
	allowWildcard := true
	if originsRaw != "" {
		parts := strings.Split(originsRaw, ",")
		origins = make([]string, 0, len(parts))
		for _, part := range parts {
			origin := strings.TrimSpace(part)
			if origin != "" {
				origins = append(origins, origin)
			}
		}
		if len(origins) == 0 {
			origins = []string{"*"}
		}
		allowWildcard = false
	}

	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	return appConfig{
		Port:              port,
		CorsOrigins:       origins,
		FrontendDir:       frontendDir,
		GinMode:           ginMode,
		AllowWildcardCORS: allowWildcard,
	}
}

func requestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %3d | %13v | %15s | %-7s %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			param.Latency.Truncate(time.Millisecond),
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	})
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
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
	config := loadAppConfig()
	gin.SetMode(config.GinMode)

	r := gin.New()
	r.Use(requestLogger())
	r.Use(gin.Recovery())
	r.Use(securityHeaders())

	// Configuração do CORS
	corsConfig := cors.Config{
		AllowOrigins:     config.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	if config.AllowWildcardCORS {
		corsConfig.AllowWildcard = true
	}
	r.Use(cors.New(corsConfig))

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

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"service":  "Grupo Mano",
			"uptime":   time.Now().UTC().Format(time.RFC3339),
			"frontend": "available",
		})
	})

	r.POST("/api/contact", func(c *gin.Context) {
		var req ContactRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "Dados invalidos para envio."})
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)
		req.Service = strings.TrimSpace(req.Service)
		req.Message = strings.TrimSpace(req.Message)

		if len(req.Name) < 3 {
			c.JSON(400, gin.H{"message": "Informe um nome valido."})
			return
		}

		if _, err := mail.ParseAddress(req.Email); err != nil {
			c.JSON(400, gin.H{"message": "Informe um e-mail valido."})
			return
		}

		if req.Service == "" {
			c.JSON(400, gin.H{"message": "Selecione o servico desejado."})
			return
		}

		if len(req.Message) < 10 {
			c.JSON(400, gin.H{"message": "A mensagem deve ter pelo menos 10 caracteres."})
			return
		}

		log.Printf("Novo contato | nome=%s email=%s servico=%s mensagem=%q", req.Name, req.Email, req.Service, req.Message)
		c.JSON(200, gin.H{"message": "Solicitacao enviada com sucesso. Nossa equipe retornara em breve."})
	})

	// Serve a página inicial e arquivos estáticos do frontend sem depender de caminho absoluto.
	r.GET("/", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.File(filepath.Join(config.FrontendDir, "index.html"))
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
			c.File(filepath.Join(config.FrontendDir, "index.html"))
			return
		}

		targetFile := filepath.Join(config.FrontendDir, cleanPath)
		targetAbs, absErr := filepath.Abs(targetFile)
		if absErr != nil || (!strings.HasPrefix(targetAbs, config.FrontendDir+string(os.PathSeparator)) && targetAbs != config.FrontendDir) {
			c.AbortWithStatus(404)
			return
		}

		if info, statErr := os.Stat(targetFile); statErr == nil && !info.IsDir() {
			applyStaticCacheHeaders(c, targetFile)
			c.File(targetFile)
			return
		}

		c.Header("Cache-Control", "no-cache")
		c.File(filepath.Join(config.FrontendDir, "index.html"))
	})

	r.Run(":" + config.Port)
}
