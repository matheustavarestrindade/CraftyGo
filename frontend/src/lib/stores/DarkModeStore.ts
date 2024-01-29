import { browser } from '$app/environment';
import type { Unsubscriber } from 'svelte/motion';
import { writable } from 'svelte/store';

// Save on local storage if browser
let defaultDarkMode = false;
if (browser) {
	if (localStorage.getItem('darkMode') === null) {
		const prefersDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
		localStorage.setItem('darkMode', prefersDarkMode ? 'true' : 'false');
	} else {
		defaultDarkMode = localStorage.getItem('darkMode') === 'true';
	}
}

const darkMode = writable<boolean>(defaultDarkMode);

darkMode.subscribe((value) => {
	localStorage.setItem('darkMode', String(value));
});

export const toggleDarkMode = () => {
	darkMode.update((value) => !value);
};
export const setDarkMode = (value: boolean) => {
	darkMode.set(value);
};

export const onDarkModeChange = (callback: (value: boolean) => void): Unsubscriber => {
	return darkMode.subscribe(callback);
};
