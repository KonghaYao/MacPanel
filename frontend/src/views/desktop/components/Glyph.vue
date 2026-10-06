<template>
    <svg viewBox="0 0 24 24" class="glyph" aria-hidden="true">
        <g v-if="name === 'containers'" fill="none" stroke="currentColor" stroke-width="1.6">
            <rect x="3" y="9" width="13" height="10" rx="2" />
            <rect x="8" y="5" width="13" height="10" rx="2" />
        </g>
        <g v-else-if="name === 'images'" fill="none" stroke="currentColor" stroke-width="1.6">
            <rect x="3" y="4" width="18" height="16" rx="2" />
            <circle cx="9" cy="9" r="1.5" fill="currentColor" stroke="none" />
            <path d="M4 16.5 8.5 12l3.5 3 3-2.2L20 17" />
        </g>
        <g v-else-if="name === 'appStore'" fill="none" stroke="currentColor" stroke-width="1.6">
            <path d="M8 8V7a4 4 0 0 1 8 0v1" />
            <path d="M6.5 8h11l-.8 11.2a1.5 1.5 0 0 1-1.5 1.3H8.8a1.5 1.5 0 0 1-1.5-1.3L6.5 8Z" />
        </g>
        <g v-else-if="name === 'monitor'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
            <rect x="3" y="3.5" width="18" height="13" rx="2" />
            <path d="M7 13.2 10.2 9l2.8 2.2L17 7.5" />
            <path d="M9 20.5h6" />
        </g>
        <g
            v-else-if="name === 'terminal'"
            fill="none"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <path d="M5 7.5 9.2 12 5 16.5" />
            <path d="M12 16.5h7" />
        </g>
        <g v-else-if="name === 'files'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
            <path d="M3.5 8.2A2 2 0 0 1 5.5 6.2h3.4l1.8 2h8.8a2 2 0 0 1 2 2v7.2a2 2 0 0 1-2 2H5.5a2 2 0 0 1-2-2V8.2Z" />
        </g>
        <g v-else-if="name === 'database'" fill="none" stroke="currentColor" stroke-width="1.6">
            <ellipse cx="12" cy="6.5" rx="6.5" ry="2.6" />
            <path d="M5.5 6.5v11c0 1.5 2.9 2.6 6.5 2.6s6.5-1.1 6.5-2.6v-11" />
            <path d="M5.5 12c0 1.5 2.9 2.6 6.5 2.6s6.5-1.1 6.5-2.6" />
        </g>
        <g v-else-if="name === 'website'" fill="none" stroke="currentColor" stroke-width="1.6">
            <circle cx="12" cy="12" r="7.5" />
            <path d="M4.5 12h15" />
            <path d="M12 4.5c2.2 2.4 2.2 10.6 0 15" />
            <path d="M12 4.5c-2.2 2.4-2.2 10.6 0 15" />
        </g>
        <g v-else-if="name === 'firewall'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
            <path d="M12 3.2 18.5 6v5.2c0 4.1-2.7 6.8-6.5 8.4-3.8-1.6-6.5-4.3-6.5-8.4V6L12 3.2Z" />
        </g>
        <g v-else-if="name === 'cronjob'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
            <circle cx="12" cy="12" r="7.5" />
            <path d="M12 8v4.5l2.8 1.8" />
        </g>
        <g v-else-if="name === 'logs'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
            <path d="M7 3.5h7.2L19 8.2V20.5H7Z" />
            <path d="M14.2 3.5V8.2H19" />
            <path d="M10 12.2h6M10 15.8h6" stroke-linecap="round" />
        </g>
        <g v-else-if="name === 'toolbox'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round">
            <rect x="3.5" y="8.5" width="17" height="10.5" rx="2" />
            <path d="M8.5 8.5V6.8A1.8 1.8 0 0 1 10.3 5h3.4a1.8 1.8 0 0 1 1.8 1.8v1.7" />
            <path d="M3.5 13h17" />
        </g>
        <g v-else-if="name === 'gpu'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
            <rect x="7" y="7" width="10" height="10" rx="2" />
            <path
                d="M9.5 3.5v3.5M12 3.5v3.5M14.5 3.5v3.5M9.5 17v3.5M12 17v3.5M14.5 17v3.5M3.5 9.5H7M3.5 12H7M3.5 14.5H7M17 9.5h3.5M17 12h3.5M17 14.5h3.5"
            />
        </g>
        <g v-else-if="name === 'homebrew'" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
            <path d="M6.5 9h9v6.2A3.6 3.6 0 0 1 11.9 18.8H10.1A3.6 3.6 0 0 1 6.5 15.2V9Z" />
            <path d="M15.5 10.6h1.6a2 2 0 0 1 0 4H15.5" />
            <path d="M9 5.2c.2 1.1 1 1.1 1 2.2M12.2 4.4c.2 1.1 1 1.1 1 2.2" />
        </g>
        <g v-else fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
            <circle cx="12" cy="12" r="3" />
            <path
                d="M12 3.2v2.2M12 18.6V20.8M3.2 12h2.2M18.6 12h2.2M5.8 5.8l1.6 1.6M16.6 16.6l1.6 1.6M18.2 5.8l-1.6 1.6M7.4 16.6l-1.6 1.6"
            />
        </g>
    </svg>
</template>

<script setup lang="ts">
import type { DockGlyph } from '@/views/desktop/apps';

defineProps<{ name: DockGlyph }>();
</script>

<style scoped>
.glyph {
    width: 28px;
    height: 28px;
    color: #fff;
    display: block;
}
</style>
