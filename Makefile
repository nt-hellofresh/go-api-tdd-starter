COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

dev-up:
	docker compose up -d

dev-down:
	docker compose down

test:
	@echo "Running tests and collecting coverage profiles..."
	@rm -f $(COVERAGE_FILE)
	@echo "mode: atomic" > $(COVERAGE_FILE)
	@for pkg in $$(go list ./internal/...); do \
		echo "Testing package $$pkg"; \
		go test -covermode=atomic -coverprofile=profile.out $$pkg; \
		if [ -f profile.out ]; then \
			tail -n +2 profile.out >> $(COVERAGE_FILE); \
			rm profile.out; \
		fi \
	done

coverage:
	@go tool cover -func=$(COVERAGE_FILE)

clean:
	rm -rf $(COVERAGE_FILE) $(COVERAGE_HTML)

