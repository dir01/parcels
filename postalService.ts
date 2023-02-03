/**
 * TrackingInfo represents final result of the service:
 * parsed and normalized representation of parcel tracking info.
 */
export type TrackingInfo = {
  trackingNumber: string;
  apiName: string;
  lastFetchedAt: Date;
  originCountry: string;
  destinationCountry: string;
  events: TrackingEvent[];
};

export type TrackingEvent = {
  time: Date;
  description: string;
  code:
    | "SHIPMENT_INFO_RECEIVED"
    | "PACKAGING_COMPLETED"
    | "DISPATCHED_FROM_WAREHOUSE"
    | "WMS_CONFIRMED" // WMS = Warehouse Management System
    | "ARRIVED_AT_SORTING_CENTER"
    | "ACCEPTED_BY_CARRIER"
    | "DEPARTED_FROM_SORTING_CENTER"
    | "LH_HO_IN_SUCCESS" // Arrived at departure transport hub
    | "TRANSIT_PORT_REROUTE_CALLBACK" //preMainCode:SINOA00241668IL TODO: What does this mean?
    | "EXPORT_CUSTOMS_CLEARANCE_STARTED"
    | "LH_HO_AIRLINE" // Leaving from departure country/region
    | "IMPORT_CUSTOMS_CLEARANCE_STARTED"
    | "IMPORT_CUSTOMS_CLEARANCE_SUCCESS"
    | "DEPARTED_ORIGIN_COUNTRY"
    | "ARRIVED_AT_LINEHAUL_OFFICE"
    | "ARRIVED_AT_CUSTOMS"
    | "DEPARTED_FROM_CUSTOMS"
    | "EXPORT_CUSTOMS_CLEARANCE_SUCCESS"
    | "UNKNOWN";
};

/**
 * PostalApi represents a single postal service API.
 * It should know 2 things:
 * 1. How to fetch raw response from the postal service's API
 * 2. How to parse raw response into TrackingInfo
 */
export interface PostalApi {
  fetch(
    trackingNumber: string,
    httpClient: typeof fetch,
  ): Promise<PostalApiResponse | null>;

  parse(rawResponse: PostalApiResponse): TrackingInfo | null;
}

/**
 * PostalApiResponse represents raw response from the API.
 * We could have used raw strings and known exceptions,
 * but I wanted to make API implementations contractually obliged
 * to indicate the status of the response while exposing possible expected error types
 */
export type PostalApiResponse = {
  id: number; // db key
  trackingNumber: string;
  apiName: string;
  firstFetchedAt: Date;
  lastFetchedAt: Date;
  responseBody: string;
  status: "success" | "rateLimited" | "notFound" | "error";
  error?: string;
};

/**
 * PostalApiResponseStorage is a storage that contains whole history of PostalAPI responses.
 * It is used to avoid unnecessary calls to the API.
 * We store whole history of responses for further analysis.
 * However, this storage is only concerned with the last response.
 */
export interface PostalApiResponseStorage {
  getLast(opts: {
    trackingNumber: string;
    apiName: string;
  }): Promise<PostalApiResponse | null>;

  append(opts: {
    trackingNumber: string;
    apiName: string;
    response: PostalApiResponse;
  }): Promise<void>;

  updateLastFetchedAt(id: number, lastFetchedAt: Date): Promise<void>;
}

export default class PostalService {
  private postalApiMap: { [key: string]: PostalApi };
  private storage: PostalApiResponseStorage;
  private hoursBetweenChecks: number;
  private expiryDays: number;
  private now: () => Date;

  constructor(opts: {
    storage: PostalApiResponseStorage;
    hoursBetweenChecks: number;
    expiryDays: number;
    postalApiMap: { [key: string]: PostalApi };
    now?: () => Date; // for testing
  }) {
    this.storage = opts.storage;
    this.hoursBetweenChecks = opts.hoursBetweenChecks;
    this.postalApiMap = opts.postalApiMap;
    this.expiryDays = opts.expiryDays;
    this.now = opts.now || (() => new Date());
  }

  startPolling() {
    setInterval(() => {
      this.poll();
    }, 1000 * 60 * 60 * this.hoursBetweenChecks);
  }

  async getTrackingInfo(trackingNumber: string): Promise<TrackingInfo[]> {
    const apiNames = Object.keys(this.postalApiMap); // TODO: smart filtering

    const promises = apiNames.map((apiName) =>
      this.getTrackingInfoFromApi(trackingNumber, apiName)
    );

    const results = await Promise.all(promises);
    return results.filter((r) => r !== null) as TrackingInfo[];
  }

  private poll() {
    console.log("Polling...");
  }

  private async getTrackingInfoFromApi(
    trackingNumber: string,
    apiName: string,
  ): Promise<TrackingInfo | null> {
    const api = this.postalApiMap[apiName];

    const storedResp = await this.storage.getLast({ trackingNumber, apiName });

    if (
      storedResp &&
      storedResp.responseBody &&
      storedResp.status === "success" &&
      this.wasRecentlyFetched(storedResp) &&
      !this.isRelatedToOldPackage(storedResp)
    ) {
      return api.parse(storedResp);
    }

    const fetchedResp = await api.fetch(trackingNumber, fetch);
    if (fetchedResp === null) {
      return null;
    }

    const now = this.now();

    if (!storedResp || storedResp.responseBody !== fetchedResp.responseBody) {
      fetchedResp.firstFetchedAt = now;
      fetchedResp.lastFetchedAt = now;
      await this.storage.append({
        trackingNumber,
        apiName,
        response: fetchedResp,
      });
      return api.parse(fetchedResp);
    } else if (storedResp) {
      await this.storage.updateLastFetchedAt(storedResp.id, now);
      return api.parse({ ...storedResp, lastFetchedAt: now });
    } else {
      return api.parse({ ...fetchedResp, lastFetchedAt: now });
    }
  }

  /*
  Contrary to popular belief, tracking numbers are not unique.
  In fact, they are often reused.
  To solve this, we discard any results that are older than a certain number of days.
  */
  private isRelatedToOldPackage(resp: PostalApiResponse): boolean {
    const diff = this.now().getTime() - resp.lastFetchedAt.getTime();
    const diffInDays = diff / 1000 / 60 / 60 / 24;
    return diffInDays > this.expiryDays;
  }

  /**
   * If we have recently checked the tracking number, we don't need to check it again.
   * This is to avoid unnecessary calls to the API.
   */
  private wasRecentlyFetched(resp: PostalApiResponse): boolean {
    const diff = this.now().getTime() - resp.lastFetchedAt.getTime();
    const diffInHours = diff / 1000 / 60 / 60;
    return diffInHours > this.hoursBetweenChecks;
  }
}
