// Точка входа бандла: загрузчик из HTML запускает этот файл, когда блобы уже в DOM.
import { startViewer } from "./app";
import { readGzipBlob } from "./container";
import { ID_SCENE } from "./format";

async function boot(): Promise<void> {
  const scene = await readGzipBlob(ID_SCENE);
  if (!scene) throw new Error("в файле нет сцены");
  startViewer(scene);
  document.getElementById("boot")?.remove();
}

boot().catch((e: unknown) => {
  console.error(e);
  const msg = e instanceof Error ? e.message : String(e);
  const box = document.getElementById("boot") ?? document.body.appendChild(document.createElement("div"));
  box.textContent = "Не удалось открыть: " + msg;
});
