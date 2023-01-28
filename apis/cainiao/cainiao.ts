import {
  PostalApi,
  PostalApiResponse,
  TrackingInfo,
} from "../../postalService.ts";

export default class API implements PostalApi {
  async fetch(
    trackingNumber: string,
    doFetch: typeof fetch,
  ): Promise<PostalApiResponse> {
    const resp = await doFetch(
      `https://global.cainiao.com/global/detail.json?mailNos=${trackingNumber}&lang=en-US`,
    );
    const responseBody = await resp.text();
    return {
      apiName: "cainiao",
      trackingNumber,
      fetchedAt: new Date(),
      responseBody,
      status: "success",
    };
  }

  parse(rawResponse: PostalApiResponse): TrackingInfo | null {
    try {
      const json = JSON.parse(rawResponse.responseBody);
      if (json.success !== true) {
        return null;
      }
      if (json.module.length > 1) {
        console.error("Unexpected number of modules in response");
      }
      if (json.module.length === 0) {
        return null;
      }
      const module = json.module[0];

      return {
        apiName: "cainiao",
        trackingNumber: rawResponse.trackingNumber,
        fetchedAt: rawResponse.fetchedAt,
        originCountry: module.originCountry,
        destinationCountry: module.destinationCountry,
      };
    } catch {
      return null;
    }
  }
}
