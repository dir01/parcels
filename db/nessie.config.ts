import {
  ClientSQLite,
  NessieConfig,
} from "https://deno.land/x/nessie@2.0.10/mod.ts";

const client = new ClientSQLite("./db/sqlite.db");

const config: NessieConfig = {
  client,
  migrationFolders: ["./db/migrations"],
  seedFolders: ["./db/seeds"],
};

export default config;
