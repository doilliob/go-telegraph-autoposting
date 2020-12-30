# go-telegraph-autoposting

 Program finds folders with name like `2020-12-30 - Some title`, uploads JPEG/PNG images (and resizes them if needed) in the folder to https://telegra.ph and makes post with title from folder name `Some title`. Program creates posts for folders with current date only.

# Configuration in `configuration.yaml`

* `telegraph_account_name` - Telegraph account name for new posts
* `maximum_image_size_mb` - image size limit (bigger files will be resized)
* `debug` - not working yet

# Install and compile:

1. Please install `dep` https://github.com/golang/dep
2. Run `dep ensure`
3. Make build folder: `mkdir build`
4. Compile project: `go build -o build/autoposting.exe`
5. Copy configuration file `cp configuration.yaml build/`
6. Change dir to build and run programm: `cd build/ && ./autoposting.exe`
