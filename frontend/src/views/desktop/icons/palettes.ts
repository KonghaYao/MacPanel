import type { DockGlyph } from '@/views/desktop/apps';

export interface IconPalette {
    /** top → mid → bottom gradient stops */
    colors: [string, string, string];
    /** desktop tint string for apps.ts */
    tint: string;
}

export const ICON_PALETTES: Record<DockGlyph, IconPalette> = {
    containers: {
        colors: ['#4db8ff', '#2496ed', '#1260b0'],
        tint: 'linear-gradient(180deg, #4db8ff 0%, #2496ed 52%, #1260b0 100%)',
    },
    images: {
        colors: ['#ff7eb3', '#d946ef', '#6366f1'],
        tint: 'linear-gradient(180deg, #ff7eb3 0%, #d946ef 52%, #6366f1 100%)',
    },
    appStore: {
        colors: ['#66b8ff', '#0a84ff', '#0066d6'],
        tint: 'linear-gradient(180deg, #66b8ff 0%, #0a84ff 52%, #0066d6 100%)',
    },
    monitor: {
        colors: ['#636366', '#3a3a3c', '#1c1c1e'],
        tint: 'linear-gradient(180deg, #636366 0%, #3a3a3c 52%, #1c1c1e 100%)',
    },
    terminal: {
        colors: ['#48484a', '#2c2c2e', '#0d0d0d'],
        tint: 'linear-gradient(180deg, #48484a 0%, #2c2c2e 52%, #0d0d0d 100%)',
    },
    files: {
        colors: ['#7ec4ff', '#3d9cf0', '#147ce5'],
        tint: 'linear-gradient(180deg, #7ec4ff 0%, #3d9cf0 52%, #147ce5 100%)',
    },
    database: {
        colors: ['#5cd6b8', '#34c759', '#248a3d'],
        tint: 'linear-gradient(180deg, #5cd6b8 0%, #34c759 52%, #248a3d 100%)',
    },
    website: {
        colors: ['#ffffff', '#f0f0f5', '#d8d8de'],
        tint: 'linear-gradient(180deg, #ffffff 0%, #f0f0f5 52%, #d8d8de 100%)',
    },
    firewall: {
        colors: ['#ff8a80', '#ff453a', '#c41e1e'],
        tint: 'linear-gradient(180deg, #ff8a80 0%, #ff453a 52%, #c41e1e 100%)',
    },
    cronjob: {
        colors: ['#ffb340', '#ff9500', '#cc7700'],
        tint: 'linear-gradient(180deg, #ffb340 0%, #ff9500 52%, #cc7700 100%)',
    },
    logs: {
        colors: ['#8e8e93', '#636366', '#48484a'],
        tint: 'linear-gradient(180deg, #8e8e93 0%, #636366 52%, #48484a 100%)',
    },
    toolbox: {
        colors: ['#ffd060', '#ff9f0a', '#e68600'],
        tint: 'linear-gradient(180deg, #ffd060 0%, #ff9f0a 52%, #e68600 100%)',
    },
    gpu: {
        colors: ['#b388ff', '#8b5cf6', '#6d28d9'],
        tint: 'linear-gradient(180deg, #b388ff 0%, #8b5cf6 52%, #6d28d9 100%)',
    },
    homebrew: {
        colors: ['#f0c878', '#d4a056', '#a67c2a'],
        tint: 'linear-gradient(180deg, #f0c878 0%, #d4a056 52%, #a67c2a 100%)',
    },
    settings: {
        colors: ['#c7c7cc', '#98989d', '#636366'],
        tint: 'linear-gradient(180deg, #c7c7cc 0%, #98989d 52%, #636366 100%)',
    },
};
