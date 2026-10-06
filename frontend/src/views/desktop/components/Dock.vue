<template>
    <div class="dock-wrap">
        <div ref="rowRef" class="dock" @mousemove="onMove" @mouseleave="onLeave">
            <template v-for="(app, index) in apps" :key="app.id">
                <span v-if="app.dividerBefore" class="sep" />
                <button
                    type="button"
                    class="dock-item"
                    :class="{ bounce: bouncingId === app.id }"
                    :style="itemStyle(index)"
                    :aria-label="t(app.titleKey)"
                    @click="open(app.id)"
                    @animationend="onBounceEnd(app.id)"
                >
                    <span class="tile">
                        <Glyph :name="app.glyph" />
                    </span>
                    <span class="tip">{{ t(app.titleKey) }}</span>
                    <i class="dot" :class="{ on: runningIds.includes(app.id) }" />
                </button>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { DesktopAppDef } from '@/views/desktop/apps';
import Glyph from '@/views/desktop/components/Glyph.vue';

defineProps<{
    apps: DesktopAppDef[];
    runningIds: string[];
}>();

const emit = defineEmits<{ open: [id: string] }>();

const { t } = useI18n();
const rowRef = ref<HTMLElement>();
const scales = ref<number[]>([]);
const bouncingId = ref('');

const itemStyle = (index: number) => {
    const scale = scales.value[index] || 1;
    return {
        '--scale': String(scale),
        transform: `scale(${scale})`,
        zIndex: String(Math.round(scale * 10)),
    };
};

const onMove = (event: MouseEvent) => {
    const row = rowRef.value;
    if (!row) {
        return;
    }
    const rowRect = row.getBoundingClientRect();
    const next: number[] = [];
    row.querySelectorAll<HTMLElement>('.dock-item').forEach((item) => {
        const center = rowRect.left + item.offsetLeft + item.offsetWidth / 2;
        const distance = Math.abs(event.clientX - center);
        const influence = 140;
        const linear = Math.max(0, 1 - distance / influence);
        const eased = linear * linear * (3 - 2 * linear);
        next.push(1 + eased * 0.75);
    });
    scales.value = next;
};

const onLeave = () => {
    scales.value = [];
};

const open = (id: string) => {
    bouncingId.value = '';
    requestAnimationFrame(() => {
        bouncingId.value = id;
    });
    emit('open', id);
};

const onBounceEnd = (id: string) => {
    if (bouncingId.value === id) {
        bouncingId.value = '';
    }
};
</script>

<style scoped lang="scss">
.dock-wrap {
    position: relative;
    z-index: 5;
    display: flex;
    justify-content: center;
    padding: 0 16px 10px;
    pointer-events: none;
    overflow: visible;
}

.dock {
    position: relative;
    pointer-events: auto;
    display: flex;
    align-items: flex-end;
    gap: 4px;
    max-width: calc(100vw - 24px);
    padding: 6px 12px 4px;
    border-radius: 20px;
    background: rgba(255, 255, 255, 0.22);
    border: 0.5px solid rgba(255, 255, 255, 0.38);
    box-shadow:
        0 0 0 0.5px rgba(0, 0, 0, 0.06),
        0 12px 36px rgba(0, 0, 0, 0.22),
        inset 0 0.5px 0 rgba(255, 255, 255, 0.5);
    -webkit-backdrop-filter: blur(36px) saturate(1.15);
    backdrop-filter: blur(36px) saturate(1.15);
}

.dock-item {
    position: relative;
    width: 60px;
    height: 72px;
    padding: 0 0 8px;
    border: 0;
    background: transparent;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
    transform-origin: center bottom;
    transition: transform 0.12s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    cursor: pointer;
}

.dock-item.bounce {
    animation: mac-dock-bounce 0.55s ease;
}

.tile {
    width: 60px;
    height: 60px;
    overflow: visible;
    display: block;
    position: relative;
}

.tip {
    position: absolute;
    left: 50%;
    bottom: calc(100% + 6px);
    transform: translateX(-50%) scale(calc(1 / var(--scale, 1)));
    transform-origin: center bottom;
    padding: 3px 8px;
    border-radius: 6px;
    background: rgba(246, 246, 246, 0.92);
    color: #1d1d1f;
    font-size: 12px;
    line-height: 1.3;
    white-space: nowrap;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
    opacity: 0;
    pointer-events: none;
}

.dock-item:hover .tip,
.dock-item:focus-visible .tip {
    opacity: 1;
}

.dot {
    width: 4px;
    height: 4px;
    margin-top: 4px;
    border-radius: 50%;
    background: rgba(40, 40, 40, 0.72);
    visibility: hidden;
}

.dot.on {
    visibility: visible;
}

.sep {
    width: 1px;
    height: 42px;
    margin: 0 4px 14px;
    background: rgba(0, 0, 0, 0.22);
    align-self: flex-end;
}

:global(.mac-desktop.is-dark) .dock {
    background: rgba(30, 30, 32, 0.42);
    border-color: rgba(255, 255, 255, 0.12);
    box-shadow:
        0 0 0 0.5px rgba(0, 0, 0, 0.24),
        0 12px 36px rgba(0, 0, 0, 0.38),
        inset 0 0.5px 0 rgba(255, 255, 255, 0.08);
}

:global(.mac-desktop.is-dark) .tip {
    background: rgba(44, 44, 46, 0.92);
    color: #f5f5f7;
}

:global(.mac-desktop.is-dark) .dot {
    background: rgba(255, 255, 255, 0.88);
}

:global(.mac-desktop.is-dark) .sep {
    background: rgba(255, 255, 255, 0.28);
}

@keyframes mac-dock-bounce {
    0%,
    100% {
        transform: translateY(0) scale(var(--scale, 1));
    }
    35% {
        transform: translateY(-18px) scale(var(--scale, 1));
    }
    58% {
        transform: translateY(0) scale(var(--scale, 1));
    }
    78% {
        transform: translateY(-8px) scale(var(--scale, 1));
    }
}

@media (prefers-reduced-motion: reduce) {
    .dock-item.bounce {
        animation: none;
    }
}
</style>
