// Типы заголовка сцены. Зеркало internal/scene/scene.go, спека – docs/FORMAT.md.

export const ID_SCENE = "dv-scene";
export const ID_VOLUME = "dv-volume";

export const MAGIC_SCENE = "DVSC";
export const SCENE_VERSION = 1;

export type Kind = "anatomy" | "bone" | "implant" | "abutment" | "sleeve" | "waxup";
export type Shading = "smooth" | "flat" | "crease";

export const KINDS: Kind[] = ["anatomy", "bone", "implant", "abutment", "sleeve", "waxup"];

export const KIND_LABELS: Record<Kind, string> = {
  anatomy: "Анатомия",
  bone: "Кость",
  implant: "Импланты",
  abutment: "Абатменты",
  sleeve: "Втулки",
  waxup: "Wax-up",
};

export interface MeshDef {
  id: string;
  name: string;
  kind: Kind;
  color: string;
  opacity: number;
  visible: boolean;
  shading: Shading;
  crease_deg?: number;
  transform?: number[]; // 16, row-major
  geom: number;
}

export interface GeomDef {
  origin: [number, number, number];
  scale: [number, number, number];
  vertices: number;
  triangles: number;
  index: "u16" | "u32";
  pos_off: number;
  pos_len: number;
  idx_off: number;
  idx_len: number;
}

export interface SceneHeader {
  title: string;
  units: string;
  meshes: MeshDef[];
  geoms: GeomDef[];
}
