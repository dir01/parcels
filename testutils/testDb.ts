import { DB } from "https://deno.land/x/sqlite@v3.7.0/mod.ts";

export function prepareTestDb() {
  const db = new DB(":memory:");
  db.execute(
    `
    CREATE TABLE postal_api_responses (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      api_name TEXT NOT NULL,
      tracking_number TEXT NOT NULL,
      first_fetched_at INTEGER NOT NULL,
      last_fetched_at INTEGER NOT NULL,
      response_body TEXT NOT NULL,
      status TEXT NOT NULL,
      error TEXT
    );
    `,
  );
  return db;
}
