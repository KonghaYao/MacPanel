<template>
    <IconShell name="files" :colors="colors">
        <path
            d="M32 14c-10.5 0-18 7.5-18 16.5v3C14 42.5 21.5 50 32 50s18-7.5 18-16.5v-3C50 21.5 42.5 14 32 14Z"
            fill="#fff"
        />
        <ellipse cx="26.5" cy="28" rx="2.4" ry="3" fill="#1a1a1a" />
        <ellipse cx="37.5" cy="28" rx="2.4" ry="3" fill="#1a1a1a" />
        <path d="M25 37c2.8 4 10.2 4 14 0" stroke="#1a1a1a" stroke-width="2.4" stroke-linecap="round" fill="none" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.files.colors;
</script>
