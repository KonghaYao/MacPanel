<template>
    <svg viewBox="0 0 64 64" class="dock-icon" aria-hidden="true">
        <defs>
            <linearGradient :id="`${name}-bg`" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" :stop-color="gradTop" />
                <stop offset="100%" :stop-color="gradBottom" />
            </linearGradient>
            <linearGradient :id="`${name}-shine`" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#fff" stop-opacity="0.42" />
                <stop offset="45%" stop-color="#fff" stop-opacity="0.08" />
                <stop offset="100%" stop-color="#fff" stop-opacity="0" />
            </linearGradient>
        </defs>
        <rect width="64" height="64" rx="14" :fill="`url(#${name}-bg)`" />
        <rect width="64" height="64" rx="14" :fill="`url(#${name}-shine)`" />

        <g v-if="name === 'containers'" fill="#fff">
            <rect x="14" y="24" width="24" height="18" rx="3" opacity="0.92" />
            <rect x="26" y="18" width="24" height="18" rx="3" opacity="0.78" />
            <rect x="14" y="30" width="24" height="18" rx="3" opacity="0.55" />
        </g>

        <g v-else-if="name === 'images'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linejoin="round">
            <rect x="12" y="16" width="40" height="32" rx="4" fill="rgba(255,255,255,0.18)" />
            <circle cx="22" cy="26" r="3.5" fill="#fff" stroke="none" />
            <path d="M16 40l10-9 8 7 6-5 8 7" stroke-linecap="round" />
        </g>

        <g v-else-if="name === 'appStore'" fill="none" stroke="#fff" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 24v-2.5a6 6 0 0 1 12 0V24" />
            <path d="M17 24h30l-2.2 22a2.2 2.2 0 0 1-2.2 1.9H21.4a2.2 2.2 0 0 1-2.2-1.9L17 24Z" fill="rgba(255,255,255,0.2)" />
            <path d="M26 34h12" />
        </g>

        <g v-else-if="name === 'monitor'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="13" y="15" width="38" height="26" rx="3.5" fill="rgba(255,255,255,0.12)" />
            <path d="M20 33l7-8 6 5 8-10" />
            <path d="M24 46h16" />
            <path d="M32 41v5" />
        </g>

        <g v-else-if="name === 'terminal'" fill="none" stroke="#fff" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round">
            <rect x="12" y="14" width="40" height="36" rx="4" fill="rgba(255,255,255,0.08)" />
            <path d="M20 26l6 6-6 6" />
            <path d="M30 38h14" />
        </g>

        <g v-else-if="name === 'files'" fill="#fff">
            <path d="M14 18h14l4 5h18a3 3 0 0 1 3 3v18a3 3 0 0 1-3 3H17a3 3 0 0 1-3-3V21a3 3 0 0 1 3-3Z" opacity="0.95" />
            <path d="M28 18v5h18" fill="none" stroke="rgba(255,255,255,0.55)" stroke-width="1.5" />
        </g>

        <g v-else-if="name === 'database'" fill="none" stroke="#fff" stroke-width="2.2">
            <ellipse cx="32" cy="20" rx="14" ry="5.5" fill="rgba(255,255,255,0.18)" />
            <path d="M18 20v22c0 3 6.3 5.5 14 5.5s14-2.5 14-5.5V20" />
            <path d="M18 32c0 3 6.3 5.5 14 5.5s14-2.5 14-5.5" />
        </g>

        <g v-else-if="name === 'website'" fill="none" stroke="#fff" stroke-width="2.1">
            <circle cx="32" cy="32" r="16" fill="rgba(255,255,255,0.12)" />
            <path d="M16 32h32" />
            <path d="M32 16c4.5 5 4.5 19 0 24" />
            <path d="M32 16c-4.5 5-4.5 19 0 24" />
        </g>

        <g v-else-if="name === 'firewall'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linejoin="round">
            <path d="M32 13 46 19v12c0 9.5-6.2 15.8-14 18.5-7.8-2.7-14-9-14-18.5V19L32 13Z" fill="rgba(255,255,255,0.18)" />
            <path d="M24 31l5.5 5.5L40 26" stroke-linecap="round" />
        </g>

        <g v-else-if="name === 'cronjob'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linecap="round">
            <circle cx="32" cy="32" r="16" fill="rgba(255,255,255,0.12)" />
            <path d="M32 20v12l8 5" />
            <circle cx="32" cy="32" r="2" fill="#fff" stroke="none" />
        </g>

        <g v-else-if="name === 'logs'" fill="none" stroke="#fff" stroke-width="2.1" stroke-linejoin="round">
            <path d="M20 14h16l10 10v24a2.5 2.5 0 0 1-2.5 2.5H20a2.5 2.5 0 0 1-2.5-2.5V16.5A2.5 2.5 0 0 1 20 14Z" fill="rgba(255,255,255,0.16)" />
            <path d="M36 14v10h10" />
            <path d="M24 34h20M24 40h16" stroke-linecap="round" />
        </g>

        <g v-else-if="name === 'toolbox'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="14" y="24" width="36" height="22" rx="3.5" fill="rgba(255,255,255,0.14)" />
            <path d="M22 24v-4.5a3 3 0 0 1 3-3h14a3 3 0 0 1 3 3V24" />
            <path d="M14 34h36" />
            <path d="M28 18v-3M36 18v-3" />
        </g>

        <g v-else-if="name === 'gpu'" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round">
            <rect x="18" y="20" width="28" height="24" rx="4" fill="rgba(255,255,255,0.16)" />
            <path d="M24 14v6M32 14v6M40 14v6M24 44v6M32 44v6M40 44v6M14 26h4M14 32h4M14 38h4M46 26h4M46 32h4M46 38h4" />
        </g>

        <g v-else-if="name === 'homebrew'" fill="none" stroke="#fff" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M22 18h20l-2 28H24L22 18Z" fill="rgba(255,255,255,0.16)" />
            <path d="M20 26h24" />
            <path d="M38 30h6a2.5 2.5 0 0 1 0 5h-6" />
            <path d="M26 14c1 2 2 2 2 4M34 12c1 2 2 2 2 4" />
        </g>

        <g v-else fill="none" stroke="#fff" stroke-width="2.2" stroke-linecap="round">
            <circle cx="32" cy="32" r="8" fill="rgba(255,255,255,0.16)" />
            <path d="M32 18v4M32 42v4M18 32h4M42 32h4M22.5 22.5l2.8 2.8M38.7 38.7l2.8 2.8M41.5 22.5l-2.8 2.8M25.3 38.7l-2.8 2.8" />
        </g>
    </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DockGlyph } from '@/views/desktop/apps';

const props = defineProps<{ name: DockGlyph }>();

const palette: Record<DockGlyph, [string, string]> = {
    containers: ['#7ec8ff', '#1d4ed8'],
    images: ['#e9d5ff', '#7c3aed'],
    appStore: ['#93c5fd', '#2563eb'],
    monitor: ['#9ca3af', '#111827'],
    terminal: ['#52525b', '#09090b'],
    files: ['#93c5fd', '#1d4ed8'],
    database: ['#6ee7b7', '#047857'],
    website: ['#67e8f9', '#0369a1'],
    firewall: ['#fda4af', '#be123c'],
    cronjob: ['#fcd34d', '#b45309'],
    logs: ['#d6d3d1', '#57534e'],
    toolbox: ['#fdba74', '#c2410c'],
    gpu: ['#c4b5fd', '#6d28d9'],
    homebrew: ['#fbbf24', '#b45309'],
    settings: ['#e5e7eb', '#4b5563'],
};

const gradTop = computed(() => palette[props.name][0]);
const gradBottom = computed(() => palette[props.name][1]);
</script>

<style scoped>
.dock-icon {
    width: 100%;
    height: 100%;
    display: block;
    filter: drop-shadow(0 1px 1px rgba(0, 0, 0, 0.18));
}
</style>
