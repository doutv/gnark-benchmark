# Benchmark Gnark on Mobile and Desktop

## Desktop
Run Gnark solidity verifier test

```bash
# 1. Install https://github.com/Consensys/gnark-solidity-checker
# 2. Run go test with tags
cd p256
go test -timeout 10m -tags solccheck,prover_checks -test.v
```

## Mobile
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