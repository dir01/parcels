import {
  AbstractMigration,
  ClientSQLite,
  Info,
} from "https://deno.land/x/nessie@2.0.10/mod.ts";

export default class extends AbstractMigration<ClientSQLite> {
  /** Runs on migrate */
  up(_info: Info): Promise<void> {
    this.client.execute(`
      CREATE TABLE postal_api_responses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        api_name TEXT NOT NULL,
        tracking_number TEXT NOT NULL,
        first_fetched_at INTEGER NOT NULL,
        last_fetched_at INTEGER NOT NULL,
        response_body TEXT NOT NULL,
        status TEXT NOT NULL,
        error TEXT
      );`);
    return Promise.resolve();
  }

  /** Runs on rollback */
  down(_info: Info): Promise<void> {
    this.client.execute("DROP TABLE postal_api_responses");
    return Promise.resolve();
  }
}
