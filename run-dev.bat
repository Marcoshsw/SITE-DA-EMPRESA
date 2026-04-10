@echo off
setlocal
set PORT=8080
set GIN_MODE=release
set FRONTEND_DIR=%~dp0frontend
set CORS_ORIGINS=*
cd /d "%~dp0backend"
go run .
