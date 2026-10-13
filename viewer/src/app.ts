import css from "./style.css";
import { parseContainer } from "./container";
import { MAGIC_SCENE, SCENE_VERSION, type SceneHeader } from "./format";
import { GeometryCache } from "./geometry";
import { buildPanel, buildToolbar } from "./panel";
import { View3D } from "./view3d";

/** Собирает интерфейс из распакованного блоба сцены. */
export function startViewer(sceneBuffer: ArrayBuffer): View3D {
  const c = parseContainer<SceneHeader>(sceneBuffer, MAGIC_SCENE, SCENE_VERSION);
  const h = c.header;
  document.title = h.title;

  document.head.append(Object.assign(document.createElement("style"), { textContent: css }));
  const app = document.createElement("div");
  app.id = "app";
  const viewHost = document.createElement("main");
  viewHost.className = "view";
  document.body.append(app);

  const view = new View3D(viewHost);
  const geoms = new GeometryCache(h, c);
  for (const def of h.meshes) view.add(def, geoms.get(def));
  geoms.releaseRaw();

  const toggle = Object.assign(document.createElement("button"), { className: "toggle", textContent: "☰", title: "Панель" });
  toggle.addEventListener("click", () => app.classList.toggle("collapsed"));
  viewHost.append(buildToolbar(view), toggle);
  app.append(buildPanel(view, h.title), viewHost);

  window.addEventListener("keydown", (e) => {
    if (e.target instanceof HTMLInputElement && e.target.type !== "checkbox" && e.target.type !== "range") return;
    if (e.code === "KeyF") view.fit();
  });

  view.setView("front");
  return view;
}
