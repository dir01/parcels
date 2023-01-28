export async function fetchCached(
  url: string | URL | Request,
  init?: RequestInit | undefined,
): Promise<Response> {
  const force = Deno.env.get("FORCE_FETCH");
  if (force || !(await existsCached(url, init))) {
    const resp = await fetch(url, init);
    await cacheResponse(url, init, resp);
    return resp;
  }

  return await getCachedResponse(url, init);
}

function getFilename(
  url: string | URL | Request,
  init: RequestInit | undefined,
) {
  const urlStr = typeof url === "string" ? url : url.toString();
  const initStr = init ? JSON.stringify(init) : "";
  let filename = urlStr + "_" + initStr;
  filename = filename.replace(/[^a-zA-Z0-9]/g, "_");
  filename = filename.replace(/_$/g, "");
  return `./testdata/${filename}.json`;
}

async function existsCached(
  url: string | URL | Request,
  init: RequestInit | undefined,
) {
  const filename = getFilename(url, init);
  try {
    const stat = await Deno.stat(filename);
    return stat.isFile;
  } catch {
    return false;
  }
}

async function cacheResponse(
  url: string | URL | Request,
  init: RequestInit | undefined,
  resp: Response,
) {
  const filename = getFilename(url, init);
  await Deno.writeFile(filename, await serializeResponse(resp));
}

async function getCachedResponse(
  url: string | URL | Request,
  init: RequestInit | undefined,
) {
  const filename = getFilename(url, init);
  return deserializeResponse(await Deno.readFile(filename));
}

async function serializeResponse(resp: Response): Promise<Uint8Array> {
  let body = await resp.text();
  try {
    body = JSON.parse(body);
  } catch {
    // ignore
  }

  const s = JSON.stringify(
    {
      status: resp.status,
      statusText: resp.statusText,
      headers: [...resp.headers],
      body: body,
    },
    null,
    2,
  );
  return new TextEncoder().encode(s);
}

function deserializeResponse(resp: Uint8Array): Response {
  const s = new TextDecoder().decode(resp);
  const { status, statusText, headers, body } = JSON.parse(s);
  let serializedBody = body;
  if (typeof body === "object") {
    serializedBody = JSON.stringify(body, null, 2);
  }
  return new Response(serializedBody, {
    status,
    statusText,
    headers: new Headers(headers),
  });
}
