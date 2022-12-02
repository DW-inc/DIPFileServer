module github.com/DW-inc/FileServer

go 1.19

require github.com/gofiber/fiber/v2 v2.40.1

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.41.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)

replace (
	github.com/andybalholm/brotli => ./libs/brotli@v1.0.4
	github.com/gofiber/fiber/v2 => ./libs/gofiber/fiber/v2@v2.37.1
	github.com/klauspost/compress => ./libs/compress@v1.15.0
	github.com/valyala/bytebufferpool => ./libs/bytebufferpool@v1.0.0
	github.com/valyala/fasthttp => ./libs/fasthttp@v1.40.0
	github.com/valyala/tcplisten => ./libs/tcplisten@v1.0.0
	golang.org/x/sys => ./libs/sys@v0.0.0-20220227234510-4e6760a101f9
)
