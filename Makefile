.PHONY: install local local-plonk android ios android-groth16 android-plonk clean

local:
	go run main.go

local-plonk:
	go run main.go plonk
	
android:
	gomobile bind --target android -androidapi 21 -o ./android/app/libs/ecdsa.aar ./ecdsa
	# Open Android Studio
	open android -a Android\ Studio

ios:
	gomobile bind --target ios -o ./ios/Ecdsa.xcframework ./ecdsa
	# Open Xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

android-groth16:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark .
	adb push gnark ecdsa.r1cs ecdsa.vkey ecdsa.zkey /data/local/tmp/
	adb shell "cd /data/local/tmp && ./gnark"

android-plonk:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark .
	adb push gnark ecdsa.plonk.r1cs ecdsa.plonk.vkey ecdsa.plonk.zkey /data/local/tmp/
	adb shell "cd /data/local/tmp && ./gnark plonk"

clean:
	rm gnark