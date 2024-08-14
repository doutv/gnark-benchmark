.PHONY: install local local-plonk android ios android-groth16 android-plonk clean

local:
	go run main.go

local-plonk:
	go run main.go plonk
	
android:
	gomobile bind --target android -androidapi 21 -o ./android/app/libs/gnark.aar ./ecdsa ./eddsa ./dummy1200k
	# Open Android Studio
	open android -a Android\ Studio

ios:
	gomobile bind --target ios -o ./ios/Gnark.xcframework ./ecdsa ./eddsa
	# Open Xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

android-binary:
	GOARCH=arm64 go build -ldflags="-s -w" -o gnark .
	adb push gnark *.r1cs *.vkey *.zkey /data/local/tmp/

android-groth16: android-binary
	adb shell "cd /data/local/tmp && ./gnark"

android-plonk: android-binary
	adb shell "cd /data/local/tmp && ./gnark plonk"

clean:
	rm gnark *.proof *.r1cs *.vkey *.zkey
