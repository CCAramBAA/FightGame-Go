@echo off
chcp 65001 >nul 2>nul
setlocal enabledelayedexpansion

:: ============================================================
::  FightGame Dev Environment Launcher (Windows)
::  Starts: Go Backend + Vue Frontend + Vue Admin
::  Prerequisites: MySQL running, Redis optional
:: ============================================================

title FightGame - Starting...

echo.
echo ============================================================
echo          FightGame Dev Environment Launcher
echo ============================================================
echo.

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"

:: ==================== 1. Check MySQL ====================
echo [1/6] Checking MySQL...

:: Common MySQL install paths
set "MYSQL_FOUND="
for %%p in (
    "C:\Program Files\MySQL\MySQL Server 8.0\bin"
    "C:\Program Files\MySQL\MySQL Server 8.4\bin"
    "C:\Program Files\MySQL\MySQL Server 5.7\bin"
    "C:\xampp\mysql\bin"
) do (
    if exist %%p\mysql.exe (
        set "MYSQL_FOUND=%%~p"
        set "PATH=!PATH!;%%~p"
    )
)

:: Read DB config from .env
set "DB_PASS="
set "DB_PORT=3306"
if exist "%ROOT_DIR%.env" (
    for /f "usebackq tokens=1,2 delims==" %%a in ("%ROOT_DIR%.env") do (
        if "%%a"=="DB_PASSWORD" set "DB_PASS=%%b"
        if "%%a"=="DB_PORT"     set "DB_PORT=%%b"
    )
)
if not defined DB_PASS set "DB_PASS=root123"

if defined MYSQL_FOUND (
    :: Try .env password and port
    mysqladmin ping -h 127.0.0.1 -P %DB_PORT% -u root -p%DB_PASS% --silent 2>nul
    if !errorlevel! equ 0 (
        echo   MySQL connected (127.0.0.1:%DB_PORT%)
        goto mysql_ok
    )
    :: Try no-password
    mysqladmin ping -h 127.0.0.1 -P %DB_PORT% -u root --silent 2>nul
    if !errorlevel! equ 0 (
        echo   [WARN] MySQL connected but password may not be "%DB_PASS%". Check .env
        goto mysql_ok
    )
    echo   [WARN] MySQL found but cannot connect (port: %DB_PORT%). Trying to start service...
    net start MySQL80 >nul 2>&1
    if !errorlevel! neq 0 (
        net start MySQL >nul 2>&1
    )
    timeout /t 3 /nobreak >nul
    mysqladmin ping -h 127.0.0.1 -P %DB_PORT% -u root -p%DB_PASS% --silent 2>nul
    if !errorlevel! equ 0 (
        echo   MySQL connected!
        goto mysql_ok
    )
)

:: MySQL not installed or unreachable
echo.
echo   ======================================================
echo   [ERROR] Cannot connect to MySQL!
echo   ======================================================
echo.
echo   This project requires MySQL 8.0.
echo   Install options:
echo.
echo   Option 1: Download MySQL Installer
echo     https://dev.mysql.com/downloads/installer/
echo     After install, create database:
echo       mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS fight_game;"
echo     Set DB_PASSWORD in %ROOT_DIR%.env
echo.
echo   Option 2: Install XAMPP (includes MySQL)
echo     https://www.apachefriends.org/
echo.
echo   Option 3: Use Docker Desktop
echo     https://www.docker.com/products/docker-desktop/
echo     Then run: docker-compose up -d mysql redis
echo.
echo   Then re-run this script.
echo   ======================================================
echo.
pause
exit /b 1

:mysql_ok
:: Ensure database exists
echo   Creating database fight_game if not exists...
mysql -h 127.0.0.1 -P %DB_PORT% -u root -p%DB_PASS% -e "CREATE DATABASE IF NOT EXISTS fight_game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul
echo   MySQL ready!

echo.

:: ==================== 2. Check Redis ====================
echo [2/6] Checking Redis...
where redis-cli >nul 2>&1
if %errorlevel% equ 0 (
    redis-cli -h 127.0.0.1 -p 6379 ping >nul 2>&1
    if %errorlevel% equ 0 (
        echo   Redis connected (127.0.0.1:6379)
    ) else (
        echo   [WARN] Redis not running (some features unavailable)
    )
) else (
    echo   [WARN] Redis not found (some features unavailable)
)

echo.

:: ==================== 3. Check Dev Tools ====================
echo [3/6] Checking dev tools...

where node >nul 2>&1
if %errorlevel% neq 0 (
    echo   [ERROR] Node.js not found. Install Node.js ^>= 18: https://nodejs.org/
    pause
    exit /b 1
)
for /f "tokens=*" %%i in ('node -v') do echo   Node.js:  %%i

where go >nul 2>&1
if %errorlevel% neq 0 (
    echo   [ERROR] Go not found. Install Go ^>= 1.21: https://go.dev/dl/
    pause
    exit /b 1
)
for /f "tokens=*" %%i in ('go version 2^>^&1') do echo   %%i

echo.

:: ==================== 4. Install Dependencies ====================
echo [4/6] Checking dependencies...

if not exist "frontend\node_modules" (
    echo   Installing frontend dependencies...
    cd /d "%ROOT_DIR%frontend"
    call npm install
    cd /d "%ROOT_DIR%"
) else (
    echo   Frontend deps: OK
)

if not exist "admin\node_modules" (
    echo   Installing admin dependencies...
    cd /d "%ROOT_DIR%admin"
    call npm install
    cd /d "%ROOT_DIR%"
) else (
    echo   Admin deps: OK
)

echo   Backend deps: will download on build

echo.

:: ==================== 5. Start Services ====================
echo [5/6] Starting services...

:: Pre-build Go backend
echo   Building backend...
cd /d "%ROOT_DIR%server"
:: Use Chinese Go proxy for faster downloads
set GOPROXY=https://goproxy.cn,https://goproxy.io,direct
set GOPRIVATE=
go build -o server.exe .\cmd\main.go 2>nul
if %errorlevel% neq 0 (
    echo   [ERROR] Backend build failed! Check Go code.
    pause
    exit /b 1
)
echo   Build OK

:: Start Go backend
start "FightGame - Backend" cmd /k "cd /d "%ROOT_DIR%server" && title FightGame Backend :8080 && echo ===== FightGame Backend :8080 ===== && echo Starting... && .\server.exe"
cd /d "%ROOT_DIR%"

echo   Waiting for backend...
set /a COUNT=0
:wait_server
timeout /t 2 /nobreak >nul
curl -s -o NUL http://localhost:8080/api/health 2>nul
if %errorlevel% equ 0 goto server_ready
set /a COUNT+=1
if %COUNT% lss 15 goto wait_server
echo   [WARN] Backend startup timeout. Check backend window.
goto skip_server_wait
:server_ready
echo   Backend ready (http://localhost:8080)
:skip_server_wait

:: Start frontend
start "FightGame - Frontend" cmd /k "cd /d "%ROOT_DIR%frontend" && title FightGame Frontend :3000 && echo ===== FightGame Frontend :3000 ===== && npm run dev"

:: Start admin
start "FightGame - Admin" cmd /k "cd /d "%ROOT_DIR%admin" && title FightGame Admin :3001 && echo ===== FightGame Admin :3001 ===== && npm run dev"

echo   Waiting for frontend build (10s)...
timeout /t 10 /nobreak >nul

echo.

:: ==================== 6. Open Browser ====================
echo [6/6] Opening browser...
start "" http://localhost:3000
timeout /t 1 /nobreak >nul
start "" http://localhost:3001

echo.
echo ============================================================
echo          All services started!
echo ------------------------------------------------------------
echo   Game Frontend:  http://localhost:3000
echo   Admin Panel:    http://localhost:3001
echo   API Server:     http://localhost:8080
echo ------------------------------------------------------------
echo   Test Accounts:
echo     Player 1: player1 / 123456
echo     Player 2: player2 / 123456
echo     Admin:    admin  / admin123
echo ------------------------------------------------------------
echo   Tip: Close individual windows to stop each service.
echo ============================================================
echo.
pause
