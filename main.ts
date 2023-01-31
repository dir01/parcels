import { DB } from "https://deno.land/x/sqlite@v3.7.0/mod.ts";
import PostalService, { PostalApi } from "./postalService.ts";
import HttpServer from "./httpServer.ts";
import SQLitePostalApiResponseStorage from "./storage.ts";
import CainiaoApi from "./apis/cainiao/cainiao.ts";

const dbPath = Deno.env.get("DB_PATH") || "./db/sqlite.db";
const storage = new SQLitePostalApiResponseStorage({ db: new DB(dbPath) });

const cainiaoApi: PostalApi = new CainiaoApi();

const postalService = new PostalService({
  storage,
  hoursBetweenChecks: 24,
  expiryDays: 30,
  postalApiMap: {
    cainiao: cainiaoApi,
  },
});

new HttpServer({ postalService, port: 8080 }).serve();
