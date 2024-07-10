# Benchmark Gnark on Mobile and Desktop

Dependecies:
- Golang
- iOS: Xcode
- Android: Android Studio, Android SDK, Android NDK

Android related environment variables on MacOS:
```bash
# SDK
export ANDROID_HOME="$HOME/Library/Android/sdk"
# NDK
export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/25.2.9519653
# adb
export PATH=$PATH:$HOME/Library/Android/sdk/platform-tools
```

```bash
make local
make android
make android-binary
make ios
```