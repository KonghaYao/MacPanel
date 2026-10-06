<template>
    <IconShell name="firewall" :colors="colors">
        <defs>
            <linearGradient id="shield-body" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#ffffff" />
                <stop offset="100%" stop-color="#ffd0cc" />
            </linearGradient>
        </defs>
        <path d="M32 10L50 18v14.5c0 11.5-8.5 19.5-18 22-9.5-2.5-18-10.5-18-22V18L32 10Z" fill="url(#shield-body)" />
        <path d="M32 10L50 18v14.5c0 11.5-8.5 19.5-18 22-9.5-2.5-18-10.5-18-22V18L32 10Z" fill="#fff" opacity="0.25" />
        <path
            d="M32 10L50 18v14.5c0 11.5-8.5 19.5-18 22-9.5-2.5-18-10.5-18-22V18L32 10Z"
            fill="none"
            stroke="#cc0000"
            stroke-width="0.6"
            opacity="0.3"
        />
        <path
            d="M22 31l7.5 7.5 14-16"
            stroke="#cc0000"
            stroke-width="4"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
        />
    </IconShell>
</template>

<script setup lang="ts">
import IconShell from '@/views/desktop/icons/IconShell.vue';
import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

const colors = ICON_PALETTES.firewall.colors;
</script>
