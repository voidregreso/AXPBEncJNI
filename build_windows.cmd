set CURRENT_DIR=%~dp0
set CGO_CFLAGS=-I%CURRENT_DIR%include\win32
set GOOS=windows
go build -buildmode=c-shared -ldflags="-w -s" -v -x -o output/windows/x64/AXPBEncJNI.dll