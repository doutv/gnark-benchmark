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

ios:
	gomobile bind --target ios ./ecdsa
	# run in xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

clean:
	rm gnark
	rm *.txt