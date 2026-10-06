<template>
    <IconShell name="images" :colors="colors">
        <circle cx="32" cy="32" r="16" fill="#ff2d55" />
        <circle cx="20" cy="26" r="12.5" fill="#ff9500" />
        <circle cx="44" cy="26" r="12.5" fill="#bf5af2" />
        <circle cx="20" cy="40" r="12.5" fill="#ffd60a" />
        <circle cx="44" cy="40" r="12.5" fill="#30d158" />
        <circle cx="32" cy="32" r="5.5" fill="#fff" opacity="0.92" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.images.colors;
</script>
