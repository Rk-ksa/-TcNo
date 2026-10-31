const port = Number(process.argv[2] || process.env.WAILS_VITE_PORT || 9245);
const deadline = Date.now() + 30_000;
const url = `http://localhost:${port}/`;

while (Date.now() < deadline) {
  try {
    const response = await fetch(url, { signal: AbortSignal.timeout(1_000) });
    if (response.ok) {
      process.exit(0);
    }
  } catch {
    // Vite has not bound its port yet; retry until the startup deadline.
  }
  await new Promise((resolve) => setTimeout(resolve, 250));
}

console.error(`Vite did not become ready at ${url} within 30 seconds.`);
process.exit(1);
