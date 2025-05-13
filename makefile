# Makefile simples
.PHONY: clean build clean-all

# Build básico
build:
	@go build -o robo.exe ./main.go

# Limpa apenas arquivos de build
clean:
	@rm -f robo.exe
	@rm -f *.exe
	@rm -rf bin/

# Limpa tudo, incluindo pastas desnecessárias
clean-all:
	@echo "Removing development files..."
	@rm -rf config/
	@rm -rf vendor/
	@rm -rf .github/
	@rm -f .gitignore
	@rm -f README.md
	@rm -f *.md
	@rm -f robo.exe
	@rm -f *.exe
	@echo "Cleanup complete!"