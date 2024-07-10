.PHONY: install local android ios

install:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	cp gnark ./android/ZPrize/app/src/main/jniLibs/armeabi-v7a/lib_gnark_.so
	cp gnark ./android/ZPrize/app/src/main/jniLibs/arm64-v8a/lib_gnark_.so
	rm gnark

local:
	go run main.go full
	
android:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	adb push gnark ecdsa.r1cs ecdsa.vkey ecdsa.zkey /data/local/tmp/
	adb shell "cd /data/local/tmp && ./gnark"
	mkdir -p android/app/src/main/resources/lib/arm64
	cp gnark ecdsa.r1cs ecdsa.vkey ecdsa.zkey android/app/src/main/resources/lib/arm64

ios:
	gomobile bind --target ios -o ./ios/Ecdsa.xcframework ./ecdsa
	# run in xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

clean:
	rm gnark
	rm *.txt