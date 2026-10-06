<template>
    <IconShell name="images" :colors="colors">
        <circle cx="32" cy="32" r="14" fill="#ff375f" opacity="0.92" />
        <circle cx="22" cy="26" r="11" fill="#ff9500" opacity="0.9" />
        <circle cx="42" cy="26" r="11" fill="#bf5af2" opacity="0.9" />
        <circle cx="22" cy="38" r="11" fill="#ffd60a" opacity="0.88" />
        <circle cx="42" cy="38" r="11" fill="#32d74b" opacity="0.88" />
        <circle cx="32" cy="32" r="5" fill="#fff" opacity="0.95" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.images.colors;
</script>
