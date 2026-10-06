<template>
    <svg viewBox="0 0 64 64" class="mac-icon" aria-hidden="true">
        <defs>
            <clipPath :id="clipId">
                <path :d="SQUIRCLE_PATH" />
            </clipPath>
            <linearGradient :id="bgId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" :stop-color="colors[0]" />
                <stop offset="48%" :stop-color="colors[1]" />
                <stop offset="100%" :stop-color="colors[2]" />
            </linearGradient>
            <linearGradient :id="shineId" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.5" />
                <stop offset="38%" stop-color="#fff" stop-opacity="0.1" />
                <stop offset="100%" stop-color="#fff" stop-opacity="0" />
            </linearGradient>
            <radialGradient :id="vignetteId" cx="0.5" cy="0.58" r="0.72">
                <stop offset="55%" stop-color="#000" stop-opacity="0" />
                <stop offset="100%" stop-color="#000" stop-opacity="0.2" />
            </radialGradient>
            <filter :id="shadowId" x="-15%" y="-8%" width="130%" height="125%">
                <feDropShadow dx="0" dy="1.5" stdDeviation="1.4" flood-color="#000" flood-opacity="0.32" />
            </filter>
        </defs>
        <g :clip-path="`url(#${clipId})`" :filter="`url(#${shadowId})`">
            <rect width="64" height="64" :fill="`url(#${bgId})`" />
            <slot />
            <rect width="64" height="64" :fill="`url(#${vignetteId})`" />
            <rect width="64" height="64" :fill="`url(#${shineId})`" />
        </g>
    </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { SQUIRCLE_PATH } from '@/views/desktop/icons/squircle';

const props = defineProps<{
    name: string;
    colors: [string, string, string];
}>();

const clipId = computed(() => `${props.name}-clip`);
const bgId = computed(() => `${props.name}-bg`);
const shineId = computed(() => `${props.name}-shine`);
const vignetteId = computed(() => `${props.name}-vig`);
const shadowId = computed(() => `${props.name}-sh`);
</script>
