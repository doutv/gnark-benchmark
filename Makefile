.PHONY: install local android ios android-binary

local:
	go run main.go full
	
android:
	gomobile bind --target android -androidapi 21 -o ./android/app/libs/ecdsa.aar ./ecdsa
	open android -a Android\ Studio

android-binary:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	adb push gnark ecdsa.r1cs ecdsa.vkey ecdsa.zkey /data/local/tmp/
	adb shell "cd /data/local/tmp && ./gnark"

ios:
	gomobile bind --target ios -o ./ios/Ecdsa.xcframework ./ecdsa
	# run in xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

clean:
	rm gnark