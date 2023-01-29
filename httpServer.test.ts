import { assertEquals } from "https://deno.land/std@0.173.0/testing/asserts.ts";
import { assertSnapshot } from "https://deno.land/std@0.168.0/testing/snapshot.ts";
import { mockFetch, unMockFetch } from "https://deno.land/x/metch@0.1.0/mod.ts";
import HttpServer from "./httpServer.ts";
import PostalService from "./postalService.ts";
import SQLitePostalApiResponseStorage from "./storage.ts";
import CainiaoApi from "./apis/cainiao/cainiao.ts";
import { fetchCached } from "./testutils/fetchCached.ts";
import { prepareTestDb } from "./testutils/testDb.ts";

Deno.test("HTTP server", async (t) => {
  const postalService = new PostalService({
    storage: new SQLitePostalApiResponseStorage({ db: prepareTestDb() }),
    hoursBetweenChecks: 0,
    expiryDays: 0,
    postalApiMap: {
      cainiao: new CainiaoApi(),
    },
    now: () => new Date("2020-01-01T00:00:00.000Z"),
  });
  const server = new HttpServer({ postalService, port: 9000 });

  async function sendRequest(trackingNumber: string) {
    const req = new Request("http://whatever/", {
      method: "POST",
      body: JSON.stringify({ trackingNumber }),
    });
    return await server.handleRequest(req);
  }

  await t.step("cainiao", async () => {
    await withMockedFetch(
      "https://global.cainiao.com/global/detail.json?mailNos=AE010698498&lang=en-US",
      async () => {
        const resp = await sendRequest("AE010698498");
        const resBody = await resp.json();

        assertEquals(resp.status, 200);
        assertEquals(resBody.length, 1);
        assertEquals(resBody[0].trackingNumber, "AE010698498");
        await assertSnapshot(t, resBody);
      }
    );
  });
});

async function withMockedFetch(url: string, f: () => Promise<void>) {
  await mockFetch(new Request(url), await fetchCached(url));
  try {
    await f();
  } finally {
    unMockFetch();
  }
}
