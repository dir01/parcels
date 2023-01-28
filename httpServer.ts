import PostalService from "./postalService.ts";
import { serve } from "https://deno.land/std@0.174.0/http/server.ts";

export default class HttpServer {
  private postalService: PostalService;
  private port: number;

  constructor(opts: { postalService: PostalService; port: number }) {
    this.postalService = opts.postalService;
    this.port = opts.port;
  }

  serve() {
    return serve(this.handleRequest, { port: this.port });
  }

  async handleRequest(req: Request): Promise<Response> {
    const body = await req.json();
    const trackingNumber = body.trackingNumber;
    const res = await this.postalService.getTrackingInfo(trackingNumber);
    return new Response(JSON.stringify(res), {
      status: 200,
      headers: { "content-type": "application/json" },
    });
  }
}
