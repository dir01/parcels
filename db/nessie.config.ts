import {
  ClientSQLite3,
  NessieConfig,
} from "https://raw.githubusercontent.com/dir01/deno-nessie/main/mod.ts";

const dbPath = Deno.env.get("DB_PATH") || "./db/sqlite.db";
const client = new ClientSQLite3(dbPath);

const config: NessieConfig = {
  client,
  migrationFolders: ["./db/migrations"],
  seedFolders: ["./db/seeds"],
};

export default config;
