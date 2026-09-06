# Digital Stamp Frontend

The customer and staff web experience for Digital Stamp Card, built with Next.js 16, React 19, TypeScript, and Tailwind CSS.

## Local development

1. Copy `.env.example` to `.env.local`.
2. Install dependencies with `bun install`.
3. Start the frontend with `bun run dev`.
4. Open <http://localhost:3000>.

The frontend expects the Go API at `NEXT_PUBLIC_API_URL`, which defaults to `http://localhost:8080/api`.

## Commands

```bash
bun run dev
bun run lint
bun run build
```

The core API client is in `lib/api.ts`, and shared loyalty types are in `types/domain.ts`.
