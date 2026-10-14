import * as THREE from "three";
import { toCreasedNormals } from "three/examples/jsm/utils/BufferGeometryUtils.js";
import type { Container } from "./container";
import type { GeomDef, MeshDef, SceneHeader } from "./format";

/** Декодирует geom из тела в индексированную BufferGeometry без нормалей. */
export function decodeGeom(g: GeomDef, c: Container<unknown>): THREE.BufferGeometry {
  const base = c.bodyOffset;
  const end = base + Math.max(g.pos_off + g.pos_len, g.idx_off + g.idx_len);
  if (end > c.buffer.byteLength) throw new Error("geom выходит за пределы тела");
  if (g.pos_len !== g.vertices * 6) throw new Error("pos_len не совпадает с vertices");

  const q = new Uint16Array(c.buffer, base + g.pos_off, g.vertices * 3);
  const pos = new Float32Array(q.length);
  const [ox, oy, oz] = g.origin;
  const [sx, sy, sz] = g.scale;
  for (let i = 0; i < q.length; i += 3) {
    pos[i] = ox + q[i] * sx;
    pos[i + 1] = oy + q[i + 1] * sy;
    pos[i + 2] = oz + q[i + 2] * sz;
  }

  const n = g.triangles * 3;
  // slice: копия, чтобы geometry не держала весь распакованный буфер сцены.
  const idx =
    g.index === "u16"
      ? new Uint16Array(c.buffer, base + g.idx_off, n).slice()
      : new Uint32Array(c.buffer, base + g.idx_off, n).slice();

  const geo = new THREE.BufferGeometry();
  geo.setAttribute("position", new THREE.BufferAttribute(pos, 3));
  geo.setIndex(new THREE.BufferAttribute(idx, 1));
  return geo;
}

/**
 * Геометрии с нормалями под нужный shading. Несколько мешей с одним geom и одинаковым
 * shading получают один и тот же объект (экономия памяти на повторяющихся втулках).
 */
export class GeometryCache {
  private raw = new Map<number, THREE.BufferGeometry>();
  private shaded = new Map<string, THREE.BufferGeometry>();

  constructor(
    private header: SceneHeader,
    private container: Container<unknown>,
  ) {}

  get(m: MeshDef): THREE.BufferGeometry {
    const angle = m.shading === "crease" ? (m.crease_deg ?? 30) : 0;
    const key = `${m.geom}|${m.shading}|${angle}`;
    let geo = this.shaded.get(key);
    if (geo) return geo;

    const raw = this.rawGeom(m.geom);
    switch (m.shading) {
      case "crease":
        // Не индексированная копия: вершины на острых рёбрах расщепляются.
        geo = toCreasedNormals(raw, THREE.MathUtils.degToRad(angle));
        break;
      case "flat":
        geo = raw; // нормали считает шейдер (material.flatShading)
        break;
      default:
        geo = raw.clone();
        geo.computeVertexNormals();
    }
    geo.computeBoundingSphere();
    geo.computeBoundingBox();
    this.shaded.set(key, geo);
    return geo;
  }

  /** Отпустить сырые геометрии после сборки сцены (те, что используют меши, останутся живы). */
  releaseRaw(): void {
    this.raw.clear();
  }

  private rawGeom(i: number): THREE.BufferGeometry {
    const def = this.header.geoms[i];
    if (!def) throw new Error(`geom ${i} не найден`);
    let g = this.raw.get(i);
    if (!g) {
      g = decodeGeom(def, this.container);
      this.raw.set(i, g);
    }
    return g;
  }
}
