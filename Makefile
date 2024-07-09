.PHONY: install run local clean

install:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	cp gnark ./android/ZPrize/app/src/main/jniLibs/armeabi-v7a/lib_gnark_.so
	cp gnark ./android/ZPrize/app/src/main/jniLibs/arm64-v8a/lib_gnark_.so
	rm gnark

local:
	go build -ldflags="-s -w" -o gnark .
	./gnark
	
run-android:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	adb push gnark /data/local/tmp/
	adb shell "/data/local/tmp/gnark"

run-ios:
	gomobile bind --target ios
	# run in xcode
	open ios/gnark-benchmark/gnark-benchmark.xcodeproj

clean:
	rm gnark
	rm *.txt