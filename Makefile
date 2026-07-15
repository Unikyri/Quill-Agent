.PHONY: e2e

E2E_COMPOSE = docker compose -p quill-e2e -f docker-compose.yml -f e2e/docker-compose.e2e.yml

# Runs against the assembled Docker stack, not mocks. Qwen is deliberately a
# required input: this is a local pre-merge proof, not a unit-CI target.
e2e:
	@set -e; set -a; if [ -f .env ]; then . ./.env; fi; set +a; \
		test -n "$$QWEN_API_KEY" || (echo "QWEN_API_KEY is required for the live E2E suite"; exit 1); \
		trap '$(E2E_COMPOSE) down -v --remove-orphans' EXIT; \
		$(E2E_COMPOSE) up -d --build; \
		attempt=0; until curl -fsS http://127.0.0.1:18080/api/v1/health >/dev/null; do \
			attempt=$$((attempt + 1)); if [ $$attempt -ge 60 ]; then echo "Timed out waiting for E2E backend health"; exit 1; fi; sleep 2; \
		done; \
		attempt=0; until curl -fsS http://127.0.0.1:13001/ >/dev/null; do \
			attempt=$$((attempt + 1)); if [ $$attempt -ge 60 ]; then echo "Timed out waiting for E2E frontend"; exit 1; fi; sleep 2; \
		done; \
	cd frontend && QWEN_API_KEY="$$QWEN_API_KEY" PLAYWRIGHT_BASE_URL=http://127.0.0.1:13001 TMPDIR=/tmp npx playwright test --config=../e2e/playwright.config.ts
