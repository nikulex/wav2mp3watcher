mkdir build
cd commandline
go build -o ../build/mp3mirror
cd ../trayapp
cp -r MP3mirror ../build/MP3mirror.app
go build -o ../build/MP3mirror.app/Contents/MacOS/MP3mirror
cd ..
