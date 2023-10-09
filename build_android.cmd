set ANDROID_NDK_HOME=E:/androidsdk/nsdk/ndk/24.0.8215888
set CURRENT_DIR=%~dp0
set CGO_CFLAGS=-I"%CURRENT_DIR%/include/android"
set CGO_ENABLED=1
set GOOS=android

set GOARCH=arm
set CC=%ANDROID_NDK_HOME%/toolchains/llvm/prebuilt/windows-x86_64/bin/armv7a-linux-androideabi21-clang.cmd
go build -buildmode=c-shared -ldflags="-w -s" -v -x -o output/android/armeabi-v7a/libAXPBEncJNI.so

echo "Build armeabi-v7a success"

set GOARCH=arm64
set CC=%ANDROID_NDK_HOME%/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android21-clang.cmd
go build -buildmode=c-shared -ldflags="-w -s" -v -x -o output/android/arm64-v8a/libAXPBEncJNI.so

echo "Build arm64-v8a success"