<template>
    <IconShell name="settings" :colors="colors">
        <defs>
            <linearGradient id="gear-metal" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#e8e8ec" />
                <stop offset="50%" stop-color="#98989d" />
                <stop offset="100%" stop-color="#48484a" />
            </linearGradient>
            <radialGradient id="gear-hub" cx="0.5" cy="0.5" r="0.5">
                <stop offset="0%" stop-color="#d8d8dc" />
                <stop offset="100%" stop-color="#636366" />
            </radialGradient>
        </defs>
        <path
            d="M32 10v6.5M32 47.5V54M10 32h6.5M47.5 32H54M16.2 16.2l4.6 4.6M43.2 43.2l4.6 4.6M47.8 16.2l-4.6 4.6M20.8 43.2l-4.6 4.6"
            stroke="url(#gear-metal)"
            stroke-width="5.5"
            stroke-linecap="round"
        />
        <circle cx="32" cy="32" r="9" fill="url(#gear-hub)" stroke="#48484a" stroke-width="0.6" />
        <circle cx="32" cy="32" r="5" fill="#636366" />
        <circle cx="32" cy="32" r="5" fill="#fff" opacity="0.2" />
        <circle cx="32" cy="32" r="2" fill="#98989d" />
        <path
            d="M32 10v6.5M32 47.5V54M10 32h6.5M47.5 32H54"
            stroke="#fff"
            stroke-width="1"
            stroke-linecap="round"
            opacity="0.25"
        />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.settings.colors;
</script>
