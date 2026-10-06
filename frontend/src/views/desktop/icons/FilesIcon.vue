<template>
    <IconShell name="files" :colors="colors">
        <path
            d="M32 13c-11 0-19 7.8-19 17.2v3.6C13 43.2 21 51 32 51s19-7.8 19-17.2v-3.6C51 20.8 43 13 32 13Z"
            fill="#fff"
        />
        <path d="M32 13c-11 0-19 7.8-19 17.2v3.6C13 43.2 21 51 32 51s19-7.8 19-17.2v-3.6C51 20.8 43 13 32 13Z" fill="#fff" opacity="0.15" />
        <ellipse cx="25.5" cy="28.5" rx="2.6" ry="3.2" fill="#2c2c2e" />
        <ellipse cx="38.5" cy="28.5" rx="2.6" ry="3.2" fill="#2c2c2e" />
        <path d="M24 38.5c3 4.2 11 4.2 16 0" stroke="#2c2c2e" stroke-width="2.6" stroke-linecap="round" fill="none" />
        <path d="M24 22c2.5-1.5 5.5-2 8-2s5.5.5 8 2" stroke="#147ce5" stroke-width="1.8" stroke-linecap="round" fill="none" opacity="0.35" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.files.colors;
</script>
