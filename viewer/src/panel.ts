import { KIND_LABELS, KINDS, type Kind } from "./format";
import { VIEWS, type Entry, type View3D, type ViewName } from "./view3d";

function el<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  props: Partial<HTMLElementTagNameMap[K]> & { className?: string } = {},
  ...children: (Node | string)[]
): HTMLElementTagNameMap[K] {
  const e = document.createElement(tag);
  Object.assign(e, props);
  e.append(...children);
  return e;
}

/** Строка меша: видимость, цвет, имя, прозрачность, «показать». */
function meshRow(view: View3D, e: Entry, onChange: () => void): HTMLElement {
  const vis = el("input", { type: "checkbox", checked: e.mesh.visible, title: "Показать/скрыть" });
  const color = el("input", { type: "color", value: "#" + e.mesh.material.color.getHexString(), title: "Цвет" });
  const name = el("span", { className: "name", textContent: e.def.name, title: e.def.name });
  const opacity = el("input", {
    type: "range",
    min: "0",
    max: "100",
    value: String(Math.round(e.mesh.material.opacity * 100)),
    title: "Прозрачность",
  });
  const focus = el("button", { className: "icon", textContent: "◎", title: "Показать крупно" });

  vis.addEventListener("change", () => {
    view.setVisible(e, vis.checked);
    onChange();
  });
  name.addEventListener("click", () => vis.click());
  color.addEventListener("input", () => view.setColor(e, color.value));
  opacity.addEventListener("input", () => view.setOpacity(e, Number(opacity.value) / 100));
  focus.addEventListener("click", () => {
    if (!e.mesh.visible) vis.click();
    view.fit([e]);
  });

  const row = el("div", { className: "row" }, vis, color, name, opacity, focus);
  // Синхронизация, когда видимость меняет групповой чекбокс.
  (row as HTMLElement & { sync?: () => void }).sync = () => {
    vis.checked = e.mesh.visible;
    row.classList.toggle("off", !e.mesh.visible);
  };
  row.classList.toggle("off", !e.mesh.visible);
  return row;
}

function group(view: View3D, kind: Kind, entries: Entry[]): HTMLElement {
  const all = el("input", { type: "checkbox", title: "Все в группе" });
  const head = el(
    "div",
    { className: "group-head" },
    all,
    el("span", { className: "name", textContent: KIND_LABELS[kind] }),
    el("span", { className: "count", textContent: String(entries.length) }),
  );
  const rows: (HTMLElement & { sync?: () => void })[] = [];
  const syncHead = () => {
    const n = entries.filter((e) => e.mesh.visible).length;
    all.checked = n === entries.length;
    all.indeterminate = n > 0 && n < entries.length;
  };
  for (const e of entries) rows.push(meshRow(view, e, syncHead));
  all.addEventListener("change", () => {
    for (const e of entries) view.setVisible(e, all.checked);
    rows.forEach((r) => r.sync?.());
    syncHead();
  });
  head.querySelector(".name")!.addEventListener("click", () => all.click());
  syncHead();
  return el("section", { className: "group" }, head, ...rows);
}

export function buildPanel(view: View3D, title: string): HTMLElement {
  const byKind = new Map<Kind, Entry[]>();
  for (const e of view.entries) {
    const list = byKind.get(e.def.kind) ?? [];
    list.push(e);
    byKind.set(e.def.kind, list);
  }
  const groups = KINDS.filter((k) => byKind.has(k)).map((k) => group(view, k, byKind.get(k)!));

  return el(
    "aside",
    { className: "panel" },
    el("h1", { textContent: title, title }),
    el("div", { className: "groups" }, ...groups),
    el(
      "p",
      { className: "help" },
      "ЛКМ – вращение, ПКМ – сдвиг, колесо – масштаб. Двойной клик по модели – центр вращения. F – вписать.",
    ),
  );
}

export function buildToolbar(view: View3D): HTMLElement {
  const bar = el("div", { className: "toolbar" });
  const fit = el("button", { textContent: "Вписать", title: "F" });
  fit.addEventListener("click", () => view.fit());
  bar.append(fit);
  for (const [name, v] of Object.entries(VIEWS)) {
    const b = el("button", { textContent: v.label });
    b.addEventListener("click", () => view.setView(name as ViewName));
    bar.append(b);
  }
  return bar;
}
