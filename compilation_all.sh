echo "Start compilation for windows 64bits"

CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_ENABLED=1 go build -o win_Spatial_Explorer.exe .


#echo "Start compilation for Mac 64bits"
#GOOS=darwin go build -o Mac_Spatial_Explorer.bin .


echo "Start compilation for linux 64bits"
go build -o linux_Spatial_Explorer.bin .
