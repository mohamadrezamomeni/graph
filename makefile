ifeq ($(OS),Windows_NT)
    EXE = .exe
    SEP = \\
    MKDIR = if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
    RMDIR = if exist $(BUILD_DIR) rmdir /S /Q $(BUILD_DIR)
else
    EXE =
    SEP = /
    MKDIR = mkdir -p $(BUILD_DIR)
    RMDIR = rm -rf $(BUILD_DIR)
endif

BUILD_DIR = bin
PBSERVER = $(BUILD_DIR)$(SEP)pbserver$(EXE)
PBCLIENT = $(BUILD_DIR)$(SEP)pbclient$(EXE)

.PHONY: all build pbserver pbclient clean

all: build

build: pbserver pbclient

pbserver:
	@echo Building pbserver...
	@$(MKDIR)
	go build -o $(PBSERVER) ./cmd/server/pbserver.go

pbclient:
	@echo Building pbclient...
	@$(MKDIR)
	go build -o $(PBCLIENT) ./cmd/client/pbclient.go



clean:
	@echo Cleaning...
	@$(RMDIR)