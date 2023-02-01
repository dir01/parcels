DENO_FLAGS = --allow-net --allow-env --allow-read --allow-write --no-remote --import-map=vendor/import_map.json

run:
	deno run $(DENO_FLAGS) main.ts

vendor:
	deno vendor main.ts *.test.ts https://deno.land/x/nessie@2.0.10/cli.ts --force

bundle:
	deno bundle main.ts bundle.js --no-remote --import-map=vendor/import_map.json

run-bundle:
	deno run $(DENO_FLAGS) bundle.js

lint:
	deno lint

format:
	deno fmt

test:
	deno test $(DENO_FLAGS) --trace-ops

test-update:
	deno test $(DENO_FLAGS) -- --update

test-refetch:
	FORCE_REFETCH=true make test

test-coverage:
	deno test $(DENO_FLAGS) --coverage=coverage
	deno coverage coverage --lcov > coverage.lcov
	genhtml -o coverage-html coverage.lcov  # brew install lcov

new-migration:
	deno run -A ./vendor/deno.land/x/nessie@2.0.10/cli.ts -c ./db/nessie.config.ts make:migration $(shell bash -c 'read -p "Enter migration name: " name; echo $$name')

migrate:
	deno run -A ./vendor/deno.land/x/nessie@2.0.10/cli.ts -c ./db/nessie.config.ts migrate

migrate-rollback:
	deno run -A ./vendor/deno.land/x/nessie@2.0.10/cli.ts -c ./db/nessie.config.ts rollback

clean:
	rm -rf coverage coverage.lcov coverage-html bundle.js

.PHONY: run vendor bundle run-bundle lint format test test-update test-refetch test-coverage new-migration migrate migrate-rollback clean
