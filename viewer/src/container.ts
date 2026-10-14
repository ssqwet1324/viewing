// Чтение блобов: base64 → gunzip → контейнер (magic | version | hdrLen | JSON | body).

export interface Container<H> {
  version: number;
  header: H;
  buffer: ArrayBuffer;
  bodyOffset: number; // смещение тела в buffer, кратно 8
}

function decodeBase64(b64: string): Uint8Array<ArrayBuffer> {
  const native = (Uint8Array as unknown as { fromBase64?: (s: string) => Uint8Array<ArrayBuffer> }).fromBase64;
  if (native) return native(b64);
  const bin = atob(b64);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out;
}

/** Читает и распаковывает блоб из <script id=...>. Элемент удаляется, чтобы не держать строку в памяти. */
export async function readGzipBlob(id: string): Promise<ArrayBuffer | null> {
  const el = document.getElementById(id);
  if (!el) return null;
  const bin = decodeBase64((el.textContent ?? "").trim());
  el.remove();
  const stream = new Blob([bin]).stream().pipeThrough(new DecompressionStream("gzip"));
  return new Response(stream).arrayBuffer();
}

export function parseContainer<H>(buffer: ArrayBuffer, magic: string, maxVersion: number): Container<H> {
  if (buffer.byteLength < 12) throw new Error("блоб обрезан");
  const dv = new DataView(buffer);
  const got = String.fromCharCode(dv.getUint8(0), dv.getUint8(1), dv.getUint8(2), dv.getUint8(3));
  if (got !== magic) throw new Error(`ожидался ${magic}, получен ${JSON.stringify(got)}`);
  const version = dv.getUint32(4, true);
  if (version < 1 || version > maxVersion) {
    throw new Error(`версия ${magic} ${version} не поддерживается (максимум ${maxVersion}), обновите генератор/вьюер`);
  }
  const hdrLen = dv.getUint32(8, true);
  if (12 + hdrLen > buffer.byteLength) throw new Error("заголовок обрезан");
  const header = JSON.parse(new TextDecoder().decode(new Uint8Array(buffer, 12, hdrLen))) as H;
  return { version, header, buffer, bodyOffset: 12 + hdrLen };
}
