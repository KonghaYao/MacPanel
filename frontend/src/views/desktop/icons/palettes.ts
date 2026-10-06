import type { DockGlyph } from '@/views/desktop/apps';

export interface IconPalette {
    /** top → mid → bottom gradient stops */
    colors: [string, string, string];
    /** desktop tint string for apps.ts */
    tint: string;
}

export const ICON_PALETTES: Record<DockGlyph, IconPalette> = {
    containers: {
        colors: ['#3d9ae8', '#2496ed', '#1565b8'],
        tint: 'linear-gradient(180deg, #3d9ae8 0%, #2496ed 45%, #1565b8 100%)',
    },
    images: {
        colors: ['#ff6b9d', '#c850c0', '#4158d0'],
        tint: 'linear-gradient(180deg, #ff6b9d 0%, #c850c0 45%, #4158d0 100%)',
    },
    appStore: {
        colors: ['#5eb3ff', '#1a9cff', '#0070e0'],
        tint: 'linear-gradient(180deg, #5eb3ff 0%, #1a9cff 45%, #0070e0 100%)',
    },
    monitor: {
        colors: ['#48484a', '#2c2c2e', '#1c1c1e'],
        tint: 'linear-gradient(180deg, #48484a 0%, #2c2c2e 45%, #1c1c1e 100%)',
    },
    terminal: {
        colors: ['#3a3a3c', '#1c1c1e', '#000000'],
        tint: 'linear-gradient(180deg, #3a3a3c 0%, #1c1c1e 45%, #000000 100%)',
    },
    files: {
        colors: ['#6eb4ff', '#3d9cf0', '#1a7ad9'],
        tint: 'linear-gradient(180deg, #6eb4ff 0%, #3d9cf0 45%, #1a7ad9 100%)',
    },
    database: {
        colors: ['#4db6ac', '#26a69a', '#00796b'],
        tint: 'linear-gradient(180deg, #4db6ac 0%, #26a69a 45%, #00796b 100%)',
    },
    website: {
        colors: ['#ffffff', '#e8e8ed', '#c7c7cc'],
        tint: 'linear-gradient(180deg, #ffffff 0%, #e8e8ed 45%, #c7c7cc 100%)',
    },
    firewall: {
        colors: ['#ff6961', '#ff453a', '#d70015'],
        tint: 'linear-gradient(180deg, #ff6961 0%, #ff453a 45%, #d70015 100%)',
    },
    cronjob: {
        colors: ['#ffffff', '#f2f2f7', '#d1d1d6'],
        tint: 'linear-gradient(180deg, #ffffff 0%, #f2f2f7 45%, #d1d1d6 100%)',
    },
    logs: {
        colors: ['#636366', '#48484a', '#2c2c2e'],
        tint: 'linear-gradient(180deg, #636366 0%, #48484a 45%, #2c2c2e 100%)',
    },
    toolbox: {
        colors: ['#ffb340', '#ff9500', '#e68600'],
        tint: 'linear-gradient(180deg, #ffb340 0%, #ff9500 45%, #e68600 100%)',
    },
    gpu: {
        colors: ['#9d7aff', '#7c4dff', '#5e35b1'],
        tint: 'linear-gradient(180deg, #9d7aff 0%, #7c4dff 45%, #5e35b1 100%)',
    },
    homebrew: {
        colors: ['#e8b86d', '#c8956c', '#8b6914'],
        tint: 'linear-gradient(180deg, #e8b86d 0%, #c8956c 45%, #8b6914 100%)',
    },
    settings: {
        colors: ['#aeaeb2', '#8e8e93', '#636366'],
        tint: 'linear-gradient(180deg, #aeaeb2 0%, #8e8e93 45%, #636366 100%)',
    },
};
