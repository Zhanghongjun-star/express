@echo off
cd /d "%~dp0"
set "target=%~1"
if "%target%"=="" goto help
if "%target%"=="api" goto api
if "%target%"=="config" goto config
if "%target%"=="build" goto build
if "%target%"=="generate" goto generate
if "%target%"=="init" goto init
if "%target%"=="all" goto all
goto help

:api
buf generate --template buf.gen.yaml
goto :eof

:config
buf generate --template buf.gen.config.yaml
goto :eof

:build
for /f "tokens=*" %%i in ('git describe --tags --always 2^>nul') do set "version=%%i"
if not defined version set "version=dev"
if not exist bin mkdir bin
go build -ldflags "-X main.Version=%version%" -o ./bin/ ./...
goto :eof

:generate
go generate ./...
go mod tidy
goto :eof

:init
go install github.com/google/wire/cmd/wire@latest
go install github.com/bufbuild/buf/cmd/buf@latest
goto :eof

:all
call :api
call :config
call :generate
goto :eof

:help
echo Usage: make [target]
echo   api       generate api proto (buf generate --template buf.gen.yaml)
echo   config    generate internal config proto
echo   build     compile into ./bin
echo   generate  go generate ./... ^&^& go mod tidy
echo   all       api + config + generate
echo   init      install wire / buf
echo   help      show this help
goto :eof
