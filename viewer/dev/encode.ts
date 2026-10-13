// Эталонный TS-энкодер сцены для dev-режима. Делает то же, что Go-генератор (docs/FORMAT.md),
// и полезен как вторая реализация спеки при отладке расхождений.
import type * as THREE from "three";
import type { GeomDef, MeshDef, SceneHeader } from "../src/format";

const align = (n: number, a: number) => Math.ceil(n / a) * a;

export class SceneEncoder {
  private chunks: { off: number; bytes: Uint8Array }[] = [];
  private size = 0;
  readonly header: SceneHeader;

  constructor(title: string) {
    this.header = { title, units: "mm", meshes: [], geoms: [] };
  }

  private put(bytes: Uint8Array): number {
    const off = align(this.size, 4);
    this.chunks.push({ off, bytes });
    this.size = off + bytes.byteLength;
    return off;
  }

  addGeometry(g: THREE.BufferGeometry): number {
    const p = g.getAttribute("position");
    const n = p.count;
    const lo = [Infinity, Infinity, Infinity];
    const hi = [-Infinity, -Infinity, -Infinity];
    for (let i = 0; i < n; i++) {
      for (let a = 0; a < 3; a++) {
        const v = p.getComponent(i, a);
        lo[a] = Math.min(lo[a], v);
        hi[a] = Math.max(hi[a], v);
      }
    }
    const scale = lo.map((l, a) => (hi[a] > l ? (hi[a] - l) / 65535 : 1)) as [number, number, number];
    const q = new Uint16Array(n * 3);
    for (let i = 0; i < n; i++) {
      for (let a = 0; a < 3; a++) {
        q[i * 3 + a] = Math.min(65535, Math.max(0, Math.round((p.getComponent(i, a) - lo[a]) / scale[a])));
      }
    }
    const src = g.index ? Array.from(g.index.array) : Array.from({ length: n }, (_, i) => i);
    const idx = n <= 65536 ? Uint16Array.from(src) : Uint32Array.from(src);

    const def: GeomDef = {
      origin: lo as [number, number, number],
      scale,
      vertices: n,
      triangles: src.length / 3,
      index: n <= 65536 ? "u16" : "u32",
      pos_off: this.put(new Uint8Array(q.buffer)),
      pos_len: q.byteLength,
      idx_off: 0,
      idx_len: idx.byteLength,
    };
    def.idx_off = this.put(new Uint8Array(idx.buffer));
    this.header.geoms.push(def);
    return this.header.geoms.length - 1;
  }

  addMesh(m: MeshDef): void {
    this.header.meshes.push(m);
  }

  finish(): ArrayBuffer {
    const json = new TextEncoder().encode(JSON.stringify(this.header));
    const hdrLen = align(12 + json.byteLength, 8) - 12;
    const out = new Uint8Array(12 + hdrLen + this.size);
    out.set(new TextEncoder().encode("DVSC"), 0);
    const dv = new DataView(out.buffer);
    dv.setUint32(4, 1, true);
    dv.setUint32(8, hdrLen, true);
    out.fill(0x20, 12, 12 + hdrLen);
    out.set(json, 12);
    for (const c of this.chunks) out.set(c.bytes, 12 + hdrLen + c.off);
    return out.buffer;
  }
}
