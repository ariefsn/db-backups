// Generates the OpenAPI client from the running backend's swagger spec.
//
// Why this script instead of calling `openapi --input <url>` directly:
//   1. The bundled @apidevtools/json-schema-ref-parser fails to resolve a
//      remote http(s) `--input` under some Node/Bun versions
//      ("Unable to resolve $ref pointer http://localhost:8080/..."). Fetching
//      the spec to a local temp file first sidesteps that entirely.
//   2. After generation we must reset OpenAPI.BASE to '/api' — the frontend
//      talks to the backend through the SvelteKit proxy at `/api`
//      (see src/hooks.server.ts), but the codegen derives BASE from the
//      backend's swagger host (http://localhost:8080).
//
// Usage: node ./scripts/generate-api.mjs [specUrl]
import { spawnSync } from 'node:child_process';
import { mkdtempSync, writeFileSync, readFileSync, rmSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { dirname, resolve, join } from 'node:path';
import { tmpdir } from 'node:os';

const here = dirname(fileURLToPath(import.meta.url));
const webRoot = resolve(here, '..');
const SPEC_URL = process.argv[2] || 'http://localhost:8080/swagger/doc.json';
const OUTPUT_DIR = resolve(webRoot, 'src/lib/api');
const OPENAPI_BIN = resolve(webRoot, 'node_modules/.bin/openapi');
const DESIRED_BASE = '/api';

const tmp = mkdtempSync(join(tmpdir(), 'dbb-swagger-'));
const specFile = join(tmp, 'doc.json');

try {
	// 1. Fetch the spec to a local file.
	console.log(`[generate-api] Fetching spec from ${SPEC_URL}`);
	const res = await fetch(SPEC_URL).catch((err) => {
		throw new Error(
			`Failed to reach ${SPEC_URL}. Is the backend running on :8080? (${err.message})`
		);
	});
	if (!res.ok) throw new Error(`Spec request returned HTTP ${res.status}`);
	writeFileSync(specFile, await res.text());

	// 2. Run the codegen against the local file.
	const gen = spawnSync(
		OPENAPI_BIN,
		['--input', specFile, '--output', OUTPUT_DIR, '--client', 'fetch'],
		{ stdio: 'inherit' }
	);
	if (gen.status !== 0) throw new Error(`openapi codegen exited with code ${gen.status}`);

	// 3. Patch OpenAPI.BASE for the /api proxy.
	const openApiFile = join(OUTPUT_DIR, 'core/OpenAPI.ts');
	const src = readFileSync(openApiFile, 'utf8');
	const patched = src.replace(/(\bBASE:\s*)'[^']*'/, `$1'${DESIRED_BASE}'`);
	if (patched === src) {
		console.warn(`[generate-api] BASE assignment not found in ${openApiFile}; left unchanged.`);
	} else {
		writeFileSync(openApiFile, patched);
		console.log(`[generate-api] Set OpenAPI.BASE = '${DESIRED_BASE}'`);
	}

	console.log('[generate-api] Done.');
} finally {
	rmSync(tmp, { recursive: true, force: true });
}
