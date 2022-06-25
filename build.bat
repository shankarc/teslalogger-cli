@REM builds the teslalogger-cli commands
@echo off
setlocal enableextensions
md bin
endlocal
for %%a in (*.go) do (
 @echo building %%a
 go build -o bin/ %%a
)