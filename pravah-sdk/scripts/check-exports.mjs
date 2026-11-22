import { createRequire } from 'node:module';
import { readFile } from 'node:fs/promises';
import { runInNewContext } from 'node:vm';
import React from 'react';

const require = createRequire(import.meta.url);
const commonjs = require('../dist/index.js');
const esm = await import('../dist/index.mjs');

const expected = ['default', 'CONNECTION_STATES', 'usePravah', 'usePravahListener', 'usePravahState'];
for (const [format, moduleExports] of [['CommonJS', commonjs], ['ESM', esm]]) {
    for (const name of expected) {
        if (!(name in moduleExports)) {
            throw new Error(`${format} export is missing ${name}`);
        }
    }
    if (typeof moduleExports.default !== 'function') {
        throw new Error(`${format} default export is not a Pravah constructor`);
    }
}

const browserBundle = await readFile(new URL('../dist/index.umd.js', import.meta.url), 'utf8');
const sandbox = { react: React, console, setTimeout, clearTimeout, setInterval, clearInterval };
sandbox.self = sandbox;
runInNewContext(browserBundle, sandbox);
for (const name of expected) {
    if (!(name in sandbox.pravahSdk)) {
        throw new Error(`UMD export is missing ${name}`);
    }
}

console.log('CommonJS, ESM, and UMD exports verified.');
