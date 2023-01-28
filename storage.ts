import { DB } from "https://deno.land/x/sqlite@v3.7.0/mod.ts";
import {
  PostalApiResponse,
  PostalApiResponseStorage,
} from "./postalService.ts";

export default class SQLitePostalApiResponseStorage
  implements PostalApiResponseStorage {
  private db: DB;

  constructor(opts: { db: DB }) {
    this.db = opts.db;
  }

  getLast({
    trackingNumber,
    apiName,
  }: {
    trackingNumber: string;
    apiName: string;
  }): Promise<PostalApiResponse | null> {
    const rows = this.db.queryEntries<{
      api_name: string;
      tracking_number: string;
      fetched_at: number;
      response_body: string;
      status: PostalApiResponse["status"];
      error?: string;
    }>(
      `
      SELECT
        api_name,
        tracking_number,
        fetched_at,
        response_body,
        status,
        error
      FROM postal_api_responses
      WHERE
        api_name = :apiName
        AND tracking_number = :trackingNumber
        ORDER BY fetched_at DESC
        LIMIT 1
        `,
      { apiName, trackingNumber },
    );
    if (rows.length === 0) {
      return Promise.resolve(null);
    }
    if (rows.length > 1) {
      throw new Error("Unexpected number of rows, this should never happen");
    }
    const row = rows[0];
    return Promise.resolve({
      apiName: row.api_name,
      trackingNumber: row.tracking_number,
      fetchedAt: new Date(row.fetched_at),
      responseBody: row.response_body,
      status: row.status,
      error: row.error,
    });
  }

  append({
    trackingNumber,
    apiName,
    response,
  }: {
    trackingNumber: string;
    apiName: string;
    response: PostalApiResponse;
  }): Promise<void> {
    const query = this.db.prepareQuery(`
      INSERT INTO postal_api_responses (
        api_name,
        tracking_number,
        fetched_at,
        response_body,
        status,
        error
      ) VALUES (
        :apiName,
        :trackingNumber,
        :fetchedAt,
        :responseBody,
        :status,
        :error
      )`);
    query.execute({
      apiName,
      trackingNumber,
      fetchedAt: response.fetchedAt.getTime(),
      responseBody: response.responseBody,
      status: response.status,
      error: response.error,
    });
    return Promise.resolve();
  }
}
