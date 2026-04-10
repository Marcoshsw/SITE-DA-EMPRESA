# Site Grupo Mano

Site institucional do Grupo Mano com backend em Go e frontend em HTML/CSS/JS.

## Estrutura
- `backend/`: servidor Go, API, contato, healthcheck e entrega do frontend
- `frontend/`: interface visual do site
- `frontend/assets/js/app.js`: animações e comportamento da página
- `Dockerfile`: build da aplicação em container
- `docker-compose.yml`: execução local com container
- `docker-compose.prod.yml`: execução em produção com Caddy
- `Caddyfile`: proxy reverso e HTTPS

## Como rodar localmente
```bash
cd backend
go run .
```

Ou, na raiz do projeto:
- Windows: `run-dev.bat`
- Linux/macOS: `./run-dev.sh`

Acesse:
- `http://localhost:8080`
- `http://localhost:8080/healthz`

## Variáveis de ambiente
- `PORT`: porta do servidor
- `GIN_MODE`: modo do Gin, normalmente `release`
- `FRONTEND_DIR`: caminho do diretório do frontend
- `CORS_ORIGINS`: origens liberadas para a API

## Endpoint de contato
- `POST /api/contact`

## Deploy
Veja os detalhes em [DEPLOY.md](DEPLOY.md).
