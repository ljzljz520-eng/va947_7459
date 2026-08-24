import { mkdir, writeFile } from "node:fs/promises";

await mkdir("dist", { recursive: true });
const html = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width"><title>Child Fitness Receiving</title></head>
<body><main><h1>Child Fitness Receiving</h1><p>Review class measurements with privacy-aware display.</p></main></body></html>`;
await writeFile("dist/index.html", html);
