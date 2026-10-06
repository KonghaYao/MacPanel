import type { DockGlyph } from '@/views/desktop/apps';

export interface IconPalette {
    /** top → mid → bottom gradient stops */
    colors: [string, string, string];
    /** desktop tint string for apps.ts */
    tint: string;
}

export const ICON_PALETTES: Record<DockGlyph, IconPalette> = {
    containers: {
        colors: ['#7ec8ff', '#3d9cf0', '#1565c0'],
        tint: 'linear-gradient(180deg, #7ec8ff 0%, #3d9cf0 48%, #1565c0 100%)',
    },
    images: {
        colors: ['#ffffff', '#f2f2f2', '#d8d8d8'],
        tint: 'linear-gradient(180deg, #ffffff 0%, #f2f2f2 48%, #d8d8d8 100%)',
    },
    appStore: {
        colors: ['#7ec0ff', '#4a9ff5', '#0066cc'],
        tint: 'linear-gradient(180deg, #7ec0ff 0%, #4a9ff5 48%, #0066cc 100%)',
    },
    monitor: {
        colors: ['#6ee07a', '#34c759', '#1a9e3a'],
        tint: 'linear-gradient(180deg, #6ee07a 0%, #34c759 48%, #1a9e3a 100%)',
    },
    terminal: {
        colors: ['#5a5a5c', '#3a3a3c', '#1a1a1a'],
        tint: 'linear-gradient(180deg, #5a5a5c 0%, #3a3a3c 48%, #1a1a1a 100%)',
    },
    files: {
        colors: ['#7ec4ff', '#4a9ff5', '#2272c9'],
        tint: 'linear-gradient(180deg, #7ec4ff 0%, #4a9ff5 48%, #2272c9 100%)',
    },
    database: {
        colors: ['#7ee08a', '#34c759', '#1e8a38'],
        tint: 'linear-gradient(180deg, #7ee08a 0%, #34c759 48%, #1e8a38 100%)',
    },
    website: {
        colors: ['#ffffff', '#eef4ff', '#c8d8f0'],
        tint: 'linear-gradient(180deg, #ffffff 0%, #eef4ff 48%, #c8d8f0 100%)',
    },
    firewall: {
        colors: ['#ff8a80', '#ff3333', '#cc0000'],
        tint: 'linear-gradient(180deg, #ff8a80 0%, #ff3333 48%, #cc0000 100%)',
    },
    cronjob: {
        colors: ['#ffc870', '#ff9500', '#cc7700'],
        tint: 'linear-gradient(180deg, #ffc870 0%, #ff9500 48%, #cc7700 100%)',
    },
    logs: {
        colors: ['#fff8b8', '#ffeb3b', '#e6c800'],
        tint: 'linear-gradient(180deg, #fff8b8 0%, #ffeb3b 48%, #e6c800 100%)',
    },
    toolbox: {
        colors: ['#ffb84d', '#ff8800', '#cc6600'],
        tint: 'linear-gradient(180deg, #ffb84d 0%, #ff8800 48%, #cc6600 100%)',
    },
    gpu: {
        colors: ['#c4a0ff', '#8b5cf6', '#5b21b6'],
        tint: 'linear-gradient(180deg, #c4a0ff 0%, #8b5cf6 48%, #5b21b6 100%)',
    },
    homebrew: {
        colors: ['#f5d080', '#d4a050', '#a07828'],
        tint: 'linear-gradient(180deg, #f5d080 0%, #d4a050 48%, #a07828 100%)',
    },
    settings: {
        colors: ['#c8c8cc', '#98989d', '#636366'],
        tint: 'linear-gradient(180deg, #c8c8cc 0%, #98989d 48%, #636366 100%)',
    },
};
