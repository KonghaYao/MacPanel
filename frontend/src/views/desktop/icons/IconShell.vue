<template>
    <svg viewBox="0 0 64 64" class="ios-icon" aria-hidden="true">
        <defs>
            <clipPath :id="clipId">
                <path :d="IOS_ICON_PATH" />
            </clipPath>
            <linearGradient :id="bgId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" :stop-color="colors[0]" />
                <stop offset="48%" :stop-color="colors[1]" />
                <stop offset="100%" :stop-color="colors[2]" />
            </linearGradient>
            <linearGradient :id="glossId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.72" />
                <stop offset="18%" stop-color="#fff" stop-opacity="0.48" />
                <stop offset="38%" stop-color="#fff" stop-opacity="0.12" />
                <stop offset="52%" stop-color="#fff" stop-opacity="0" />
            </linearGradient>
            <linearGradient :id="gelBandId" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.55" />
                <stop offset="100%" stop-color="#fff" stop-opacity="0" />
            </linearGradient>
            <linearGradient :id="innerShadowId" x1="0.5" y1="0.55" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#000" stop-opacity="0" />
                <stop offset="45%" stop-color="#000" stop-opacity="0.08" />
                <stop offset="100%" stop-color="#000" stop-opacity="0.38" />
            </linearGradient>
            <linearGradient :id="bevelId" x1="0.5" y1="0" x2="0.5" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.55" />
                <stop offset="8%" stop-color="#fff" stop-opacity="0.12" />
                <stop offset="92%" stop-color="#000" stop-opacity="0.06" />
                <stop offset="100%" stop-color="#000" stop-opacity="0.32" />
            </linearGradient>
            <filter :id="shadowId" x="-28%" y="-18%" width="156%" height="150%" color-interpolation-filters="sRGB">
                <feDropShadow dx="0" dy="3" stdDeviation="2.8" flood-color="#000" flood-opacity="0.52" />
            </filter>
        </defs>
        <g :filter="`url(#${shadowId})`">
            <g :clip-path="`url(#${clipId})`">
                <path :d="IOS_ICON_PATH" :fill="`url(#${bgId})`" />
                <slot />
                <rect x="0" y="0" width="64" height="22" :fill="`url(#${gelBandId})`" />
                <path :d="IOS_ICON_PATH" :fill="`url(#${innerShadowId})`" />
                <path :d="IOS_ICON_PATH" :fill="`url(#${glossId})`" />
            </g>
            <path :d="IOS_ICON_PATH" fill="none" :stroke="`url(#${bevelId})`" stroke-width="0.85" />
        </g>
    </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { IOS_ICON_PATH } from '@/views/desktop/icons/ios-icon';

const props = defineProps<{
    name: string;
    colors: [string, string, string];
}>();

const clipId = computed(() => `${props.name}-clip`);
const bgId = computed(() => `${props.name}-bg`);
const glossId = computed(() => `${props.name}-gloss`);
const gelBandId = computed(() => `${props.name}-gel`);
const innerShadowId = computed(() => `${props.name}-inner`);
const bevelId = computed(() => `${props.name}-bevel`);
const shadowId = computed(() => `${props.name}-shadow`);
</script>

<style scoped>
.ios-icon {
    width: 100%;
    height: 100%;
    display: block;
    overflow: visible;
}
</style>
