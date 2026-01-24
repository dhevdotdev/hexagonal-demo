@echo off
REM Setup script for Hexagonal Architecture Demo

echo ========================================
echo Hexagonal Architecture Demo - Setup
echo ========================================
echo.

echo Step 1: Creating database...
psql -U dhev -h localhost -c "CREATE DATABASE bookstore;" 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Database 'bookstore' created successfully!
) else (
    echo Database 'bookstore' may already exist, continuing...
)
echo.

echo Step 2: Running migrations...
psql -U dhev -h localhost -d bookstore -f migrations\001_create_books_table.sql
if %ERRORLEVEL% EQU 0 (
    echo Migrations applied successfully!
) else (
    echo Failed to apply migrations. Please check PostgreSQL connection.
    exit /b 1
)
echo.

echo Step 3: Building applications...
cd cmd\api
go build -o ..\..\bin\api.exe main.go
cd ..\..

cd cmd\cli
go build -o ..\..\bin\cli.exe main.go
cd ..\..

echo Executables built in bin\ directory
echo.

echo ========================================
echo Setup completed successfully!
echo ========================================
echo.
echo You can now run:
echo   bin\api.exe      - Start the HTTP API server
echo   bin\cli.exe      - Use the CLI tool
echo.
echo Example CLI commands:
echo   bin\cli.exe create "Clean Code" "Robert C. Martin" "978-0132350884"
echo   bin\cli.exe list
echo.
