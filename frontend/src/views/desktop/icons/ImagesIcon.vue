<template>
    <IconShell name="images" :colors="colors">
        <defs>
            <radialGradient id="petal-red" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#ff6680" />
                <stop offset="100%" stop-color="#ff2d55" />
            </radialGradient>
            <radialGradient id="petal-orange" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#ffb340" />
                <stop offset="100%" stop-color="#ff9500" />
            </radialGradient>
            <radialGradient id="petal-yellow" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#ffe040" />
                <stop offset="100%" stop-color="#ffd60a" />
            </radialGradient>
            <radialGradient id="petal-green" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#5ae878" />
                <stop offset="100%" stop-color="#30d158" />
            </radialGradient>
            <radialGradient id="petal-blue" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#5ac8fa" />
                <stop offset="100%" stop-color="#007aff" />
            </radialGradient>
            <radialGradient id="petal-purple" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#d080ff" />
                <stop offset="100%" stop-color="#bf5af2" />
            </radialGradient>
        </defs>
        <circle cx="32" cy="32" r="17" fill="url(#petal-red)" />
        <circle cx="19" cy="26" r="13" fill="url(#petal-orange)" />
        <circle cx="45" cy="26" r="13" fill="url(#petal-purple)" />
        <circle cx="19" cy="40" r="13" fill="url(#petal-yellow)" />
        <circle cx="45" cy="40" r="13" fill="url(#petal-green)" />
        <circle cx="32" cy="32" r="13" fill="url(#petal-blue)" />
        <circle cx="32" cy="32" r="5" fill="#fff" opacity="0.95" />
        <circle cx="32" cy="32" r="5" fill="none" stroke="#fff" stroke-width="0.5" opacity="0.4" />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.images.colors;
</script>
