@echo off
REM ==================== FightGame SSL 证书生成脚本 ====================
REM 需要安装 OpenSSL 或使用 Git Bash / WSL

echo.
echo Generating self-signed SSL certificate for development...
echo.

openssl req -x509 -nodes -days 365 -newkey rsa:2048 ^
    -keyout docker/nginx/ssl/server.key ^
    -out docker/nginx/ssl/server.crt ^
    -subj "/C=CN/ST=Beijing/L=Beijing/O=FightGame/CN=localhost"

if %errorlevel% equ 0 (
    echo.
    echo [OK] SSL certificate generated successfully!
    echo   - docker/nginx/ssl/server.key
    echo   - docker/nginx/ssl/server.crt
) else (
    echo.
    echo [ERROR] Failed to generate SSL certificate.
    echo Make sure OpenSSL is installed.
    echo Download: https://slproweb.com/products/Win32OpenSSL.html
)

pause
