<template>
    <IconShell name="toolbox" :colors="colors">
        <path
            d="M41 19l-9 9a6 6 0 1 1-8.5 8.5l-9 9 4 4 9-9a6 6 0 1 1 8.5-8.5l9-9-4-4Z"
            fill="#fff"
        />
        <circle cx="26.5" cy="37.5" r="5" fill="none" stroke="#fff" stroke-width="3.8" />
        <rect x="38" y="17" width="6" height="3" rx="1.2" fill="#fff" opacity="0.7" transform="rotate(45 41 18.5)" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.toolbox.colors;
</script>
