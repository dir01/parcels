/**
 * TrackingInfo represents final result of the service:
 * parsed and normalized representation of parcel tracking info.
 */
export type TrackingInfo = {
  trackingNumber: string;
  apiName: string;
  fetchedAt: Date;
  originCountry: string;
  destinationCountry: string;
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
  trackingNumber: string;
  apiName: string;
  fetchedAt: Date;
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

  async getTrackingInfo(trackingNumber: string): Promise<TrackingInfo[]> {
    const apiNames = Object.keys(this.postalApiMap); // TODO: smart filtering

    const promises = apiNames.map((apiName) =>
      this.getTrackingInfoFromApi(trackingNumber, apiName)
    );

    const results = await Promise.all(promises);
    return results.filter((r) => r !== null) as TrackingInfo[];
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

    if (storedResp && fetchedResp.responseBody !== storedResp.responseBody) {
      await this.storage.append({
        trackingNumber,
        apiName,
        response: fetchedResp,
      });
    }

    return api.parse(fetchedResp);
  }

  /*
  Contrary to popular belief, tracking numbers are not unique.
  In fact, they are often reused.
  To solve this, we discard any results that are older than a certain number of days.
  */
  private isRelatedToOldPackage(stored: PostalApiResponse): boolean {
    const diff = this.now().getTime() - stored.fetchedAt.getTime();
    const diffInDays = diff / 1000 / 60 / 60 / 24;
    return diffInDays > this.expiryDays;
  }

  /**
   * If we have recently checked the tracking number, we don't need to check it again.
   * This is to avoid unnecessary calls to the API.
   */
  private wasRecentlyFetched(info: PostalApiResponse): boolean {
    const diff = this.now().getTime() - info.fetchedAt.getTime();
    const diffInHours = diff / 1000 / 60 / 60;
    return diffInHours > this.hoursBetweenChecks;
  }
}
