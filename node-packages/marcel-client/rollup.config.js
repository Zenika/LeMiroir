/* eslint-env node */
import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import postcss from 'rollup-plugin-postcss'
import preprocess from 'svelte-preprocess';
import css from 'rollup-plugin-css-only';
import { terser } from 'rollup-plugin-terser';
import livereload from 'rollup-plugin-livereload';
import replace from '@rollup/plugin-replace'

const production = !process.env.ROLLUP_WATCH;

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/build/bundle.js'
	},
	plugins: [
		replace({
			IS_DEV_MODE: !production
		}),
		svelte({
			compilerOptions: { dev: !production },
			preprocess: preprocess()
		}),

		css({ output: 'bundle.css' }),
		postcss(),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins.
		resolve({
			browser: true,
			dedupe: ['svelte']
		}),
		commonjs(),

		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),
		!production && livereload('public'),


		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
};

function serve() {
	let started = false;

	return {
		writeBundle() {
			if (!started) {
				started = true;

				// eslint-disable-next-line no-undef
				require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
					stdio: ['ignore', 'inherit', 'inherit'],
					shell: true
				});
			}
		}
	};
}
