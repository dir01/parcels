import { Database } from "https://deno.land/x/sqlite3@0.7.3/mod.ts";
import PostalService, { PostalApi } from "./postalService.ts";
import HttpServer from "./httpServer.ts";
import SQLitePostalApiResponseStorage from "./storage.ts";
import CainiaoApi from "./apis/cainiao/cainiao.ts";

const dbPath = Deno.env.get("DB_PATH") || "./db/database.sqlite";
let db: Database;
let attmpts = 0;
while (true) {
  try {
    db = new Database(dbPath);
    break;
  } catch (e) {
    if (attmpts > 5) {
      console.error("Failed to connect to database, exiting...");
      Deno.exit(1);
    }
    console.error(e);
    console.log("Retrying in 5 seconds...");
    await new Promise((resolve) => setTimeout(resolve, 5000));
    attmpts++;
  }
}

const storage = new SQLitePostalApiResponseStorage({ db });

const cainiaoApi: PostalApi = new CainiaoApi();

const postalService = new PostalService({
  storage,
  hoursBetweenChecks: 24,
  expiryDays: 30,
  postalApiMap: {
    cainiao: cainiaoApi,
  },
});

postalService.startPolling();
new HttpServer({ postalService, port: 8080 }).serve();
