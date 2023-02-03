import { Database } from "https://deno.land/x/sqlite3@0.7.3/mod.ts";
import {
  PostalApiResponse,
  PostalApiResponseStorage,
} from "./postalService.ts";

export default class SQLitePostalApiResponseStorage
  implements PostalApiResponseStorage
{
  private db: Database;

  constructor(opts: { db: Database }) {
    this.db = opts.db;
  }

  getLast({
    trackingNumber,
    apiName,
  }: {
    trackingNumber: string;
    apiName: string;
  }): Promise<PostalApiResponse | null> {
    const rows = this.db
      .prepare(
        `
      SELECT
        id,
        api_name,
        tracking_number,
        first_fetched_at,
        last_fetched_at,
        response_body,
        status,
        error
      FROM postal_api_responses
      WHERE
        api_name = :apiName
        AND tracking_number = :trackingNumber
        ORDER BY last_fetched_at DESC
        LIMIT 1
        `
      )
      .all<{
        id: number;
        api_name: string;
        tracking_number: string;
        first_fetched_at: number;
        last_fetched_at: number;
        response_body: string;
        status: PostalApiResponse["status"];
        error?: string;
      }>({ apiName, trackingNumber });
    if (rows.length === 0) {
      return Promise.resolve(null);
    }
    if (rows.length > 1) {
      throw new Error("Unexpected number of rows, this should never happen");
    }
    const row = rows[0];
    return Promise.resolve({
      id: row.id,
      apiName: row.api_name,
      trackingNumber: row.tracking_number,
      firstFetchedAt: new Date(row.first_fetched_at),
      lastFetchedAt: new Date(row.last_fetched_at),
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
    const query = this.db
      .prepare(
        `
      INSERT INTO postal_api_responses (
        api_name,
        tracking_number,
        first_fetched_at,
        last_fetched_at,
        response_body,
        status,
        error
      ) VALUES (
        :apiName,
        :trackingNumber,
        :firstFetchedAt,
        :lastFetchedAt,
        :responseBody,
        :status,
        :error
      )`
      )
      .run({
        apiName,
        trackingNumber,
        firstFetchedAt: response.firstFetchedAt.getTime(),
        lastFetchedAt: response.lastFetchedAt.getTime(),
        responseBody: response.responseBody,
        status: response.status,
        error: response.error ?? null,
      });

    return Promise.resolve();
  }

  updateLastFetchedAt(id: number, lastFetchedAt: Date): Promise<void> {
    this.db
      .prepare(
        `UPDATE postal_api_responses
        SET
          last_fetched_at = :lastFetchedAt,
        WHERE
          id = :id`
      )
      .run({ id, lastFetchedAt: lastFetchedAt.getTime() });
    return Promise.resolve();
  }
}
