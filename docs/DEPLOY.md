# Deploy do Site Grupo Mano

## 1. O que já está pronto
- Backend em Go servindo o frontend
- Endpoint de contato em `POST /api/contact`
- Configuração por ambiente via `PORT`, `GIN_MODE`, `FRONTEND_DIR` e `CORS_ORIGINS`
- Healthcheck em `GET /healthz`
- Deploy em container com [infra/Dockerfile](../infra/Dockerfile), [infra/docker-compose.yml](../infra/docker-compose.yml) e [infra/docker-compose.prod.yml](../infra/docker-compose.prod.yml)
- HTTPS com Caddy via [infra/Caddyfile](../infra/Caddyfile)

## 2. Variáveis de ambiente
### Desenvolvimento
- `PORT=8080`
- `GIN_MODE=release`
- `FRONTEND_DIR=/app/frontend` ou caminho local equivalente
- `CORS_ORIGINS=*`

### Produção
- `SITE_DOMAIN=seu-dominio.com`
- `PORT=8080`
- `GIN_MODE=release`
- `FRONTEND_DIR=/app/frontend`
- `CORS_ORIGINS=https://seu-dominio.com`

## 3. Subida local com container
```bash
docker compose -f infra/docker-compose.yml up --build
```

Acesso local:
- Site: `http://localhost:8080`
- Healthcheck: `http://localhost:8080/healthz`

## 4. Subida em produção
1. Apontar o DNS do domínio para o IP do servidor
2. Definir `SITE_DOMAIN` no ambiente do servidor
3. Subir os serviços:

```bash
docker compose -f infra/docker-compose.prod.yml up -d --build
```

## 5. O que cada camada faz
### `backend/main.go`
- Serve API, frontend e assets estáticos
- Valida e recebe o formulário de contato
- Aplica cache, headers de segurança, logs e healthcheck

### `frontend/index.html`
- Estrutura visual da página
- Carrega o JavaScript externo versionado
- Mantém o formulário, galerias e seções do site

### `frontend/assets/js/app.js`
- Controla animações, scroll suave, galeria, carregamento dos serviços e logo animada
- Faz o `fetch` para `GET /api/services`

### `infra/Dockerfile`
- Compila o backend Go e monta a imagem final do site

### `infra/docker-compose.yml`
- Sobe o site em container local com a mesma configuração de produção base

### `infra/docker-compose.prod.yml`
- Sobe o app Go e o Caddy juntos
- Publica HTTPS nas portas 80 e 443

### `infra/Caddyfile`
- Faz proxy reverso para o app Go
- Ativa HTTPS automático
- Adiciona headers de segurança na borda

## 6. Verificação pós-deploy
- `GET /` responde `200`
- `GET /healthz` responde `200`
- `POST /api/contact` aceita payload válido
- O domínio abre com HTTPS
- Os cards de serviço carregam normalmente

## 7. Observação
- O site já está estruturado em Go no backend. O frontend continua em HTML/CSS/JS, o que é o padrão correto para preservar design e performance.