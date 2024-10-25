import { defineConfig } from 'vite';

export default defineConfig({
    root: 'src',
    server: {
        port: 3000,
    },
    build: {
        outDir: '../dist',
        emptyOutDir: true,
        minify: true,
        copyPublicDir: true,
        manifest: true,
        sourcemap: true,
        lib: {
            entry: 'index.html',
            name: 'HTMX',
            formats: ['es'],
            fileName: 'main',
        },
    },
});
