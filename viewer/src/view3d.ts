import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { RoomEnvironment } from "three/examples/jsm/environments/RoomEnvironment.js";
import type { Kind, MeshDef } from "./format";

/** Материальные параметры по типу: металл чуть глянцевее кости. */
const SURFACE: Record<Kind, { roughness: number; metalness: number }> = {
  anatomy: { roughness: 0.75, metalness: 0 },
  bone: { roughness: 0.8, metalness: 0 },
  implant: { roughness: 0.35, metalness: 0.6 },
  abutment: { roughness: 0.3, metalness: 0.6 },
  sleeve: { roughness: 0.4, metalness: 0.2 },
  waxup: { roughness: 0.6, metalness: 0 },
};

export interface Entry {
  def: MeshDef;
  mesh: THREE.Mesh<THREE.BufferGeometry, THREE.MeshStandardMaterial>;
}

/** Направления камеры от цели. Система координат: Z вверх, −Y вперёд, +X – лево пациента. */
export const VIEWS = {
  front: { label: "Спереди", dir: new THREE.Vector3(0, -1, 0) },
  back: { label: "Сзади", dir: new THREE.Vector3(0, 1, 0) },
  left: { label: "Слева", dir: new THREE.Vector3(1, 0, 0) },
  right: { label: "Справа", dir: new THREE.Vector3(-1, 0, 0) },
  // Небольшой сдвиг по Y: вид строго вдоль up вырожден для OrbitControls,
  // а знак сдвига делает перед челюсти верхом экрана.
  top: { label: "Сверху", dir: new THREE.Vector3(0, 0.001, 1).normalize() },
  bottom: { label: "Снизу", dir: new THREE.Vector3(0, -0.001, -1).normalize() },
} as const;
export type ViewName = keyof typeof VIEWS;

export class View3D {
  readonly renderer: THREE.WebGLRenderer;
  readonly scene = new THREE.Scene();
  readonly camera = new THREE.PerspectiveCamera(30, 1, 0.1, 1000);
  readonly controls: OrbitControls;
  readonly entries: Entry[] = [];
  private pending = false;

  constructor(private host: HTMLElement) {
    this.renderer = new THREE.WebGLRenderer({ antialias: true });
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    this.renderer.setClearColor(0x1b1c21);
    host.append(this.renderer.domElement);

    const pmrem = new THREE.PMREMGenerator(this.renderer);
    this.scene.environment = pmrem.fromScene(new RoomEnvironment(), 0.04).texture;
    this.scene.environmentIntensity = 0.6;
    pmrem.dispose();

    // Фонарь у камеры: освещённая сторона всегда та, на которую смотрим.
    const head = new THREE.DirectionalLight(0xffffff, 1.6);
    head.position.set(0.3, 0.4, 0);
    const headTarget = new THREE.Object3D();
    headTarget.position.set(0, 0, -1);
    this.camera.add(head, headTarget);
    head.target = headTarget;
    this.scene.add(this.camera);

    this.camera.up.set(0, 0, 1);
    this.controls = new OrbitControls(this.camera, this.renderer.domElement);
    this.controls.zoomToCursor = true;
    this.controls.addEventListener("change", () => this.render());

    this.renderer.domElement.addEventListener("dblclick", (e) => this.pickTarget(e));
    new ResizeObserver(() => this.resize()).observe(host);
    this.resize();
  }

  add(def: MeshDef, geometry: THREE.BufferGeometry): Entry {
    const material = new THREE.MeshStandardMaterial({
      color: new THREE.Color(def.color),
      side: THREE.DoubleSide, // у STL бывает вывернутый обход
      flatShading: def.shading === "flat",
      ...SURFACE[def.kind],
    });
    const mesh = new THREE.Mesh(geometry, material);
    mesh.name = def.id;
    if (def.transform) {
      mesh.matrixAutoUpdate = false;
      mesh.matrix.set(...(def.transform as Parameters<THREE.Matrix4["set"]>));
    }
    mesh.visible = def.visible;
    const entry = { def, mesh };
    this.setOpacity(entry, def.opacity);
    this.scene.add(mesh);
    this.entries.push(entry);
    return entry;
  }

  setOpacity(e: Entry, v: number): void {
    const m = e.mesh.material;
    m.opacity = v;
    const transparent = v < 0.999;
    if (m.transparent !== transparent) {
      m.transparent = transparent;
      // Полупрозрачное не пишет глубину: иначе закрывает импланты внутри кости.
      m.depthWrite = !transparent;
      m.needsUpdate = true;
    }
    e.mesh.renderOrder = transparent ? 1 : 0;
    this.render();
  }

  setVisible(e: Entry, v: boolean): void {
    e.mesh.visible = v;
    this.render();
  }

  setColor(e: Entry, hex: string): void {
    e.mesh.material.color.set(hex);
    this.render();
  }

  /** Вписать в кадр видимые объекты (или переданные), сохранив направление взгляда. */
  fit(only?: Entry[], dir?: THREE.Vector3): void {
    const list = (only ?? this.entries).filter((e) => e.mesh.visible || only);
    const box = new THREE.Box3();
    for (const e of list.length ? list : this.entries) box.expandByObject(e.mesh);
    if (box.isEmpty()) return;
    const sphere = box.getBoundingSphere(new THREE.Sphere());
    const r = Math.max(sphere.radius, 1);

    const vFov = THREE.MathUtils.degToRad(this.camera.fov);
    const hFov = 2 * Math.atan(Math.tan(vFov / 2) * this.camera.aspect);
    const dist = (r / Math.sin(Math.min(vFov, hFov) / 2)) * 1.15;

    const d = (dir ?? this.camera.position.clone().sub(this.controls.target)).clone().normalize();
    if (d.lengthSq() === 0) d.copy(VIEWS.front.dir);
    this.controls.target.copy(sphere.center);
    this.camera.position.copy(sphere.center).addScaledVector(d, dist);
    // near от размера сцены, а не от дистанции: при зуме вплотную ничего не обрезается.
    this.camera.near = r * 0.005;
    this.camera.far = dist + r * 20;
    this.camera.updateProjectionMatrix();
    this.controls.update();
  }

  setView(name: ViewName): void {
    this.fit(undefined, VIEWS[name].dir);
  }

  render(): void {
    if (this.pending) return;
    this.pending = true;
    requestAnimationFrame(() => {
      this.pending = false;
      this.renderer.render(this.scene, this.camera);
    });
  }

  private resize(): void {
    const w = this.host.clientWidth;
    const h = this.host.clientHeight;
    if (!w || !h) return;
    this.renderer.setSize(w, h, false);
    this.camera.aspect = w / h;
    this.camera.updateProjectionMatrix();
    this.render();
  }

  /** Двойной клик по поверхности – туда переносится центр вращения. */
  private pickTarget(e: MouseEvent): void {
    const rect = this.renderer.domElement.getBoundingClientRect();
    const ndc = new THREE.Vector2(((e.clientX - rect.left) / rect.width) * 2 - 1, -((e.clientY - rect.top) / rect.height) * 2 + 1);
    const ray = new THREE.Raycaster();
    ray.setFromCamera(ndc, this.camera);
    const visible = this.entries.filter((x) => x.mesh.visible).map((x) => x.mesh);
    const hit = ray.intersectObjects(visible, false)[0];
    if (!hit) return;
    // Сдвигаем камеру вместе с целью, чтобы картинка не прыгала по повороту.
    const delta = hit.point.clone().sub(this.controls.target);
    this.controls.target.add(delta);
    this.camera.position.add(delta);
    this.controls.update();
  }
}
