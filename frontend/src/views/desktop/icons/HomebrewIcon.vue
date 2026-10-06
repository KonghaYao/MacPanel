<template>
    <IconShell name="homebrew" :colors="colors">
        <path d="M21 17h22l-2.2 32H23.2L21 17Z" fill="#fff8ee" />
        <rect x="19" y="26" width="26" height="5.5" rx="1.5" fill="#fff" />
        <ellipse cx="32" cy="17" rx="13" ry="4.5" fill="#fff" />
        <path
            d="M25 13.5c1.8 2.2 3.2 2.2 3.5 4.2M32 11.5c1.8 2.2 3.2 2.2 3.5 4.2"
            stroke="#fff8ee"
            stroke-width="2.8"
            stroke-linecap="round"
            fill="none"
        />
        <path d="M39 33h7a3 3 0 0 1 0 6h-7" stroke="#a67c2a" stroke-width="3.2" stroke-linecap="round" fill="none" />
        <ellipse cx="32" cy="36" rx="8" ry="10" fill="#f5deb3" opacity="0.5" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.homebrew.colors;
</script>
