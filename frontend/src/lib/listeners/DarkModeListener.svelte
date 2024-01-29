<script lang="ts">
	import { onDarkModeChange } from "$lib/stores/DarkModeStore";
	import { onDestroy } from "svelte";


    const unsubscribe = onDarkModeChange((enable: boolean) => {
        if (enable) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    })
    onDestroy(unsubscribe)
</script>
