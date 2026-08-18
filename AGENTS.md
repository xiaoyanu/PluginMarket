# Repository Guidelines

## Project Structure & Module Organization

PluginMarket is split into a Go backend and a Nuxt frontend. Backend code lives in `server/`: `main.go` starts the Gin service, `router/` wires routes, `controller/` handles HTTP workflows, `repository/` owns data access, `model/` defines GORM models, and `database/` handles migration and seeding. Frontend code lives in `web/app/`: `pages/` maps to Nuxt routes, `components/` contains Vue UI pieces, `composables/` contains shared API and state logic, `stores/` uses Pinia, and `utils/` holds helpers. Static frontend assets are under `web/public/` and `web/app/assets/`; README screenshots are in `resource/`.

## Build, Test, and Development Commands

Backend:

- `cd server && go mod download`: install Go dependencies.
- `cd server && go run .`: run the API service, normally on `127.0.0.1:8855`.
- `cd server && go test ./...`: run all Go tests.
- `cd server && go fmt ./...`: format Go packages before committing.

Frontend:

- `cd web && pnpm install`: install Nuxt dependencies.
- `cd web && pnpm dev`: start the Nuxt dev server, normally on `localhost:3000`.
- `cd web && pnpm build`: create a production build.
- `cd web && pnpm preview`: preview the production build locally.

## Coding Style & Naming Conventions

Use standard Go formatting and package-level organization for backend changes. Keep controllers focused on request handling and move reusable database logic into repositories. Frontend files use Vue single-file components with PascalCase component names such as `PluginCard.vue`; composables follow `useX.ts`, for example `useApiFetch.ts`. Prefer TypeScript helpers in `web/app/utils/` when logic is shared across pages.

## Testing Guidelines

There is no broad test suite checked in yet. Add focused Go tests as `*_test.go` beside the package under test, then run `go test ./...`. For frontend utility tests, follow the README pattern and use Node's native runner, for example `node --test app/utils/*.test.ts`. Always run `pnpm build` when changing Nuxt pages, layouts, or configuration.

## Commit & Pull Request Guidelines

Recent commits use concise Chinese prefixes such as `更新：` and `修复：`, with short bullet-style details when needed. Keep commits focused, for example `修复：插件详情页下载地址显示`. Pull requests should describe the change, affected backend or frontend areas, test results, and any database, upload storage, environment, or `server/config.yaml` changes. Include screenshots for visible UI changes.

## Security & Configuration Tips

Do not commit production secrets, SMTP credentials, JWT secrets, database passwords, generated Nuxt output, or uploaded files. `server/uploads/`, `web/.output/`, `web/.nuxt/`, and related build artifacts are intentionally ignored.
