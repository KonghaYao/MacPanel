import type { DockGlyph } from '@/views/desktop/apps';

const symbolAsset = (name: string) =>
    new URL(`../../../assets/desktop/icons/symbols/${name}.png`, import.meta.url).href;

export const SYMBOL_URLS: Record<DockGlyph, string> = {
    containers: symbolAsset('containers'),
    images: symbolAsset('images'),
    appStore: symbolAsset('appStore'),
    monitor: symbolAsset('monitor'),
    terminal: symbolAsset('terminal'),
    files: symbolAsset('files'),
    database: symbolAsset('database'),
    website: symbolAsset('website'),
    firewall: symbolAsset('firewall'),
    cronjob: symbolAsset('cronjob'),
    logs: symbolAsset('logs'),
    toolbox: symbolAsset('toolbox'),
    gpu: symbolAsset('gpu'),
    homebrew: symbolAsset('homebrew'),
    settings: symbolAsset('settings'),
};
