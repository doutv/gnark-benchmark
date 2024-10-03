# Benchmark Gnark on Mobile and Desktop

Dependecies:
- Golang, Gomobile
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
# Local
make local
make local-plonk

# iOS Xcode app
make ios

# Android Studio app
make android
# Android Run binary executable in adb shell
make android-groth16
make android-plonk
```

export LD_LIBRARY_PATH=/opt/icicle/