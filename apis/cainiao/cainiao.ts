import {
  PostalApi,
  PostalApiResponse,
  TrackingEvent,
  TrackingInfo,
} from "../../postalService.ts";

type CainiaoDetail = {
  time: number;
  desc: string; // "preMainCode:SINOA00241668IL",
  standerdDesc: string; // "Carrier update",
  descTitle: string; // "Carrier note:",
  actionCode: // cat testdata/https___global_cainiao_com_global_detail_json_mailNos_AE010698498_lang_en_US.json | jq -r '.body.module[0].detailList[] | "| \"" + .actionCode + "\" //" + .standerdDesc'
    | "CC_IM_SUCCESS" //Import customs clearance complete
    | "CC_HO_OUT_SUCCESS" //Departed from customs
    | "CC_HO_IN_SUCCESS" //Arrived at customs
    | "CUSTOMS_ARRIVED_IN_AREA_CALLBACK" //Carrier update
    | "LH_ARRIVE" //Arrived at linehual office
    | "LH_DEPART" //Departed from departure country/region
    | "CC_IM_START" //Import customs clearance started
    | "LH_HO_AIRLINE" //Leaving from departure country/region
    | "CC_EX_SUCCESS" //Export customs clearance complete
    | "CC_EX_START" //Export customs clearance started
    | "TRANSIT_PORT_REROUTE_CALLBACK" //Carrier update
    | "LH_HO_IN_SUCCESS" //Arrived at departure transport hub
    | "SC_OUTBOUND_SUCCESS" //[Shatian Town] Departed from sorting center
    | "PU_PICKUP_SUCCESS" //Received by logistics company
    | "SC_INBOUND_SUCCESS" //[Shatian Town] Processing at sorting center
    | "WMS_CONFIRMED" //Carrier update
    | "GWMS_OUTBOUND" //Package shipped out from warehouse
    | "GWMS_PACKAGE" //Package ready for shipping from warehouse
    | "GWMS_ACCEPT"; //Shipment information received by warehouse electronically};
};

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
    if (rawResponse.status !== "success") {
      return null;
    }

    try {
      return this.transform(rawResponse);
    } catch (e) {
      console.error(e);
      return null;
    }
  }

  private transform(rawResponse: PostalApiResponse): TrackingInfo | null {
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
      destinationCountry: module.destCountry,
      events: (module.detailList || [])
        .map(this.transformDetail)
        .filter((e: TrackingEvent | null) => e !== null),
    };
  }

  private transformDetail(detail: CainiaoDetail): TrackingEvent | null {
    const code = CODES_MAP[detail.actionCode] || "UNKNOWN";
    return {
      time: new Date(detail.time),
      description: detail.standerdDesc,
      code: code as TrackingEvent["code"],
    };
  }
}

const CODES_MAP: {
  [actionCode in CainiaoDetail["actionCode"]]: TrackingEvent["code"];
} = {
  GWMS_ACCEPT: "SHIPMENT_INFO_RECEIVED",
  GWMS_PACKAGE: "PACKAGING_COMPLETED",
  GWMS_OUTBOUND: "DISPATCHED_FROM_WAREHOUSE",
  WMS_CONFIRMED: "WMS_CONFIRMED", // FIXME: What does this mean?
  SC_INBOUND_SUCCESS: "ARRIVED_AT_SORTING_CENTER",
  PU_PICKUP_SUCCESS: "ACCEPTED_BY_CARRIER",
  SC_OUTBOUND_SUCCESS: "DEPARTED_FROM_SORTING_CENTER",
  LH_HO_IN_SUCCESS: "LH_HO_IN_SUCCESS", // FIXME: What does this mean?
  TRANSIT_PORT_REROUTE_CALLBACK: "TRANSIT_PORT_REROUTE_CALLBACK", // FIXME: What does this mean?
  CC_EX_START: "EXPORT_CUSTOMS_CLEARANCE_STARTED",
  CC_EX_SUCCESS: "EXPORT_CUSTOMS_CLEARANCE_SUCCESS",
  LH_HO_AIRLINE: "LH_HO_AIRLINE", // FIXME
  CC_IM_START: "IMPORT_CUSTOMS_CLEARANCE_STARTED",
  LH_DEPART: "DEPARTED_ORIGIN_COUNTRY",
  LH_ARRIVE: "ARRIVED_AT_LINEHAUL_OFFICE",
  CC_HO_IN_SUCCESS: "ARRIVED_AT_CUSTOMS",
  CC_HO_OUT_SUCCESS: "DEPARTED_FROM_CUSTOMS",
  CC_IM_SUCCESS: "IMPORT_CUSTOMS_CLEARANCE_SUCCESS",
  CUSTOMS_ARRIVED_IN_AREA_CALLBACK: "ARRIVED_AT_CUSTOMS",
};
