<template>
    <svg viewBox="0 0 64 64" class="mac-icon" aria-hidden="true">
        <defs>
            <clipPath :id="clipId">
                <path :d="SQUIRCLE_PATH" />
            </clipPath>
            <linearGradient :id="bgId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" :stop-color="colors[0]" />
                <stop offset="52%" :stop-color="colors[1]" />
                <stop offset="100%" :stop-color="colors[2]" />
            </linearGradient>
            <radialGradient :id="specId" cx="0.5" cy="0.08" r="0.62" fx="0.5" fy="0">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.62" />
                <stop offset="42%" stop-color="#fff" stop-opacity="0.14" />
                <stop offset="100%" stop-color="#fff" stop-opacity="0" />
            </radialGradient>
            <linearGradient :id="shadeId" x1="0.5" y1="0.55" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#000" stop-opacity="0" />
                <stop offset="55%" stop-color="#000" stop-opacity="0.06" />
                <stop offset="100%" stop-color="#000" stop-opacity="0.28" />
            </linearGradient>
            <linearGradient :id="rimId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.42" />
                <stop offset="12%" stop-color="#fff" stop-opacity="0.06" />
                <stop offset="88%" stop-color="#000" stop-opacity="0.04" />
                <stop offset="100%" stop-color="#000" stop-opacity="0.2" />
            </linearGradient>
            <filter :id="shadowId" x="-22%" y="-12%" width="144%" height="138%" color-interpolation-filters="sRGB">
                <feDropShadow dx="0" dy="2.2" stdDeviation="2.1" flood-color="#000" flood-opacity="0.38" />
            </filter>
        </defs>
        <g :filter="`url(#${shadowId})`">
            <g :clip-path="`url(#${clipId})`">
                <path :d="SQUIRCLE_PATH" :fill="`url(#${bgId})`" />
                <slot />
                <path :d="SQUIRCLE_PATH" :fill="`url(#${shadeId})`" />
                <path :d="SQUIRCLE_PATH" :fill="`url(#${specId})`" />
            </g>
            <path :d="SQUIRCLE_PATH" fill="none" :stroke="`url(#${rimId})`" stroke-width="0.65" />
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
const specId = computed(() => `${props.name}-spec`);
const shadeId = computed(() => `${props.name}-shade`);
const rimId = computed(() => `${props.name}-rim`);
const shadowId = computed(() => `${props.name}-shadow`);
</script>

<style scoped>
.mac-icon {
    width: 100%;
    height: 100%;
    display: block;
    overflow: visible;
}
</style>
