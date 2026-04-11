import { serve, file } from 'bun';
import { join } from 'path';

const port = process.env.FRONTEND_PORT || 8080;
const baseDir = './dist';

serve({
  port: port as number,
  async fetch(req) {
    const url = new URL(req.url);
    let path = url.pathname;

    if (path === '/') path = '/index.html';

    let targetFile = file(join(baseDir, path));

    if (!(await targetFile.exists())) {
      targetFile = file(join(baseDir, 'index.html'));
    }

    return new Response(targetFile);
  },
});

console.log(`Bun server running on port ${port}`);
