module github.com/DW-inc/FileServer

go 1.19

require (
	github.com/gofiber/fiber/v2 v2.40.1
	github.com/joho/godotenv v1.4.0
)

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

// GORM///////////////////////
require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.1
)

replace (
	github.com/go-sql-driver/mysql => ./libs/go-sql-driver/mysql@v1.6.0
	github.com/jinzhu/inflection => ./libs/jinzhu/inflection@v1.0.0
	github.com/jinzhu/now => ./libs/jinzhu/now@v1.1.4
	gorm.io/driver/mysql => ./libs/driver/mysql@v1.3.3
	gorm.io/gorm => ./libs/gorm@v1.23.1
)

// GORM///////////////////////

replace github.com/joho/ => ./libs/godotenv@v1.4.0
