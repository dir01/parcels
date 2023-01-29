run:
	deno run --allow-net --allow-env --allow-read --allow-write main.ts

bundle:
	deno bundle main.ts bundle.js

run-bundle:
	deno run --allow-net --allow-env --allow-read --allow-write bundle.js

lint:
	deno lint

format:
	deno fmt

test:
	deno test --allow-net --allow-env --allow-read --allow-write --trace-ops

test-update:
	deno test --allow-net --allow-env --allow-read --allow-write -- --update

test-refetch:
	FORCE_REFETCH=true deno test --allow-net --allow-env --allow-read --allow-write

test-coverage:
	deno test --allow-net --allow-env --coverage=coverage
	deno coverage coverage --lcov > coverage.lcov
	genhtml -o coverage-html coverage.lcov

new-migration:
	deno run -A https://deno.land/x/nessie/cli.ts -c ./db/nessie.config.ts make:migration $(shell bash -c 'read -p "Enter migration name: " name; echo $$name')

migrate:
	deno run -A https://deno.land/x/nessie/cli.ts -c ./db/nessie.config.ts migrate

migrate-rollback:
	deno run -A https://deno.land/x/nessie/cli.ts -c ./db/nessie.config.ts rollback

clean:
	rm -rf coverage coverage.lcov coverage-html bundle.js

