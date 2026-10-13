# dentview

Генератор автономных HTML-вьюеров для дентальной имплантологии: STL (+ позже DICOM) → один `.html`,
который открывается офлайн в браузере.

```
task build        # вьюер (esbuild) + bin/dentview.exe
task demo         # bin/sample.html из testdata/sample
task viewer:dev   # dev-сервер вьюера на синтетической сцене, http://localhost:5173
task check        # vet + тесты + typecheck
```

Формат: [docs/FORMAT.md](docs/FORMAT.md). Конфиг: [testdata/sample/config.json](testdata/sample/config.json),
схема – [docs/config.schema.json](docs/config.schema.json).

Реальные данные пациентов – только в `testdata/private/` (в .gitignore).
