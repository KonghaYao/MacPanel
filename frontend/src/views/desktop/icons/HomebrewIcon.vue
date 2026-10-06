<template>
    <IconShell name="homebrew" :colors="colors">
        <path d="M22 18h20l-2 30H24L22 18Z" fill="#fff5e6" opacity="0.95" />
        <rect x="20" y="26" width="24" height="5" rx="1.5" fill="#fff" />
        <ellipse cx="32" cy="18" rx="12" ry="4" fill="#fff" />
        <path
            d="M26 14c1.5 2 3 2 3 4M32 12c1.5 2 3 2 3 4"
            stroke="#fff5e6"
            stroke-width="2.5"
            stroke-linecap="round"
            fill="none"
        />
        <path d="M38 32h6a2.5 2.5 0 0 1 0 5h-6" stroke="#8b6914" stroke-width="3" stroke-linecap="round" fill="none" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.homebrew.colors;
</script>
