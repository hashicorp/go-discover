default: copywriteheaders

.PHONY: deps
deps:
	@go install github.com/hashicorp/copywrite@b3e6599f43beff698f471c6f46888045453fa030 # v0.25.3

.PHONY: copywriteheaders
copywriteheaders:
	@echo "==> Running copywrite headers plan..."
	@copywrite headers --plan
	@echo "==> Done"
