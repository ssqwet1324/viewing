// Dev-страница: процедурная сцена → тот же бинарный формат → startViewer.
// Не нужен ни Go, ни реальные данные.
import * as THREE from "three";
import { startViewer } from "../src/app";
import { SceneEncoder } from "./encode";

const enc = new SceneEncoder("Dev: синтетическая сцена");

// Кость: половина тора, перед челюсти в сторону −Y.
const bone = new THREE.TorusGeometry(25, 6, 32, 128, Math.PI).rotateZ(Math.PI).translate(0, 5, 8);
enc.addMesh({ id: "bone", name: "Кость", kind: "bone", color: "#e8e0d0", opacity: 0.7, visible: true, shading: "smooth", geom: enc.addGeometry(bone) });

// Импланты, абатменты и втулки в локальных координатах + transform. Один geom на все копии.
const implant = new THREE.CylinderGeometry(1.9, 1.5, 11, 40, 20).rotateX(Math.PI / 2).translate(0, 0, -5.5);
const abutment = new THREE.CylinderGeometry(1.4, 1.7, 4, 24).rotateX(Math.PI / 2).translate(0, 0, 2);
const sleeve = new THREE.CylinderGeometry(3.4, 3.4, 5, 32, 1, true).rotateX(Math.PI / 2).translate(0, 0, 8.5);
const gi = enc.addGeometry(implant);
const ga = enc.addGeometry(abutment);
const gs = enc.addGeometry(sleeve);

for (const [tooth, a, tilt] of [["15", 0.2, 0.2], ["12", 0.42, 0.1], ["22", 0.58, -0.1], ["25", 0.8, -0.2]] as const) {
  const ang = Math.PI * a;
  const m = new THREE.Matrix4()
    .makeTranslation(25 * Math.cos(ang), -25 * Math.sin(ang) + 5, 14)
    .multiply(new THREE.Matrix4().makeRotationX(tilt));
  // three хранит column-major, формат – row-major.
  const t = m.clone().transpose().toArray();
  enc.addMesh({ id: `i${tooth}`, name: `Имплант ${tooth}`, kind: "implant", color: "#9aa6b4", opacity: 1, visible: true, shading: "crease", crease_deg: 30, transform: t, geom: gi });
  enc.addMesh({ id: `a${tooth}`, name: `Абатмент ${tooth}`, kind: "abutment", color: "#c9a96e", opacity: 1, visible: true, shading: "crease", crease_deg: 30, transform: t, geom: ga });
  enc.addMesh({ id: `s${tooth}`, name: `Втулка ${tooth}`, kind: "sleeve", color: "#5f9fd6", opacity: 0.9, visible: true, shading: "crease", crease_deg: 30, transform: t, geom: gs });
}

const wax = new THREE.BoxGeometry(5, 5, 5).translate(0, -20, 17.5);
enc.addMesh({ id: "waxup", name: "Wax-up", kind: "waxup", color: "#f3efe2", opacity: 1, visible: false, shading: "flat", geom: enc.addGeometry(wax) });

startViewer(enc.finish());
document.getElementById("boot")?.remove();
