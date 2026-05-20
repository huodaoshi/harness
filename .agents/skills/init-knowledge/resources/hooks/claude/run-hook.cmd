: << 'CMDBLOCK'
@echo off
REM Polyglot wrapper: Windows runs this cmd block; Unix bash skips to shell via heredoc.
REM Windows: find bash.exe below, then run extensionless hook script in this directory.
REM Unix: lines from here through CMDBLOCK are ignored by bash.
REM Hook scripts use no extension (e.g. session-start) to avoid Windows .sh auto-wrapping.
REM Usage: run-hook.cmd <hook-name-without-extension> [args...]

if "%~1"=="" (
    echo run-hook.cmd: missing hook script name >&2
    exit /b 1
)

set "HOOK_DIR=%~dp0"

REM Git for Windows default install paths
if exist "C:\Program Files\Git\bin\bash.exe" (
    "C:\Program Files\Git\bin\bash.exe" "%HOOK_DIR%%~1" %2 %3 %4 %5 %6 %7 %8 %9
    exit /b %ERRORLEVEL%
)
if exist "C:\Program Files (x86)\Git\bin\bash.exe" (
    "C:\Program Files (x86)\Git\bin\bash.exe" "%HOOK_DIR%%~1" %2 %3 %4 %5 %6 %7 %8 %9
    exit /b %ERRORLEVEL%
)

REM bash on PATH (Git Bash, MSYS2, Cygwin, etc.)
where bash >nul 2>nul
if %ERRORLEVEL% equ 0 (
    bash "%HOOK_DIR%%~1" %2 %3 %4 %5 %6 %7 %8 %9
    exit /b %ERRORLEVEL%
)

REM No bash: exit 0 so hooks fail-open (SessionStart injection skipped)
exit /b 0
CMDBLOCK

# Unix: run named hook script in this directory via bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_NAME="$1"
shift
exec bash "${SCRIPT_DIR}/${SCRIPT_NAME}" "$@"
