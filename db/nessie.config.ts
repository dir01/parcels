import {
  ClientSQLite,
  NessieConfig,
} from "https://deno.land/x/nessie@2.0.10/mod.ts";

const dbPath = Deno.env.get("DB_PATH") || "./db/sqlite.db";
const client = new ClientSQLite(dbPath);

const config: NessieConfig = {
  client,
  migrationFolders: ["./db/migrations"],
  seedFolders: ["./db/seeds"],
};

export default config;
