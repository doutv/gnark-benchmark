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
# Groth16
make local
# Xcode app
make ios
# Android Studio app
make android
# Run binary executable in adb shell 
make android-groth16

# Plonk
make local-plonk
make android-plonk
```