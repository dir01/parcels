DENO_FLAGS = --unstable --no-remote --import-map=vendor/import_map.json
DENO_PERMISSIONS = --allow-net --allow-env --allow-read --allow-write --allow-ffi

run:
	deno run $(DENO_FLAGS) $(DENO_PERMISSIONS) main.ts

vendor:
	deno vendor main.ts *.test.ts db/nessie.config.ts https://raw.githubusercontent.com/dir01/deno-nessie/main/cli.ts https://deno.land/x/nessie@2.0.10/mod.ts --force

bundle:
	deno bundle main.ts bundle.js --unstable

run-bundle:
	deno run $(DENO_FLAGS) $(DENO_PERMISSIONS) bundle.js

lint:
	deno lint

format:
	deno fmt


RUN_TEST = deno test -A --unstable --trace-ops

test:
	$(RUN_TEST)

test-update:
	$(RUN_TEST) -- --update

test-refetch:
	FORCE_REFETCH=true make test

test-coverage:
	$(RUN_TEST) --coverage=coverage
	deno coverage coverage --lcov > coverage.lcov
	genhtml -o coverage-html coverage.lcov  # brew install lcov


MIGRATOR = deno run -A --unstable --no-remote --import-map=vendor/import_map.json https://raw.githubusercontent.com/dir01/deno-nessie/main/cli.ts -c ./db/nessie.config.ts

new-migration:
	$(MIGRATOR) make:migration $(shell bash -c 'read -p "Enter migration name: " name; echo $$name')

migrate:
	$(MIGRATOR) migrate

migrate-rollback:
	$(MIGRATOR) rollback

clean:
	rm -rf coverage coverage.lcov coverage-html bundle.js

.PHONY: run vendor bundle run-bundle lint format test test-update test-refetch test-coverage new-migration migrate migrate-rollback clean
