BASE_DIR := bin
OS := darwin linux windows
ARCH := amd64 arm64

os_dir = $(foreach os,$(OS),${BASE_DIR}/$(os)/$(arch))
all_dirs := $(foreach arch,$(ARCH),$(os_dir))
mkcmd := mkdir -p $(foreach dir,$(all_dirs),$(dir))

os_build = $(foreach os,$(OS),GOOS=$(os) GOARCH=$(arch) go build -o $(BASE_DIR)/$(os)/$(arch)/ancientPlotter . &&)
all_build = $(foreach arch,$(ARCH),$(os_build)) echo Done

directories:
	@$(mkcmd)

clean:
	@rm -rf bin

all: clean directories
	@echo "Building for $(ARCH) on $(OS)"
	@$(all_build)

daggerbuild:
	dagger call build --src=. 