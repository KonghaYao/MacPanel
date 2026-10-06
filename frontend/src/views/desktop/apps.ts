import { ICON_PALETTES } from '@/views/desktop/icons/palettes';

export type DockGlyph =
    | 'containers'
    | 'images'
    | 'appStore'
    | 'monitor'
    | 'terminal'
    | 'files'
    | 'database'
    | 'website'
    | 'firewall'
    | 'cronjob'
    | 'logs'
    | 'toolbox'
    | 'gpu'
    | 'homebrew'
    | 'settings';

export interface DesktopAppDef {
    id: string;
    titleKey: string;
    route: string;
    glyph: DockGlyph;
    tint: string;
    platformFeature?: string;
    dividerBefore?: boolean;
}

export interface DesktopWindowState {
    id: string;
    appId: string;
    titleKey: string;
    kind: 'iframe' | 'about';
    route: string;
    x: number;
    y: number;
    width: number;
    height: number;
    z: number;
    minimized: boolean;
    maximized: boolean;
}

export const DESKTOP_APPS: DesktopAppDef[] = [
    {
        id: 'containers',
        titleKey: 'desktop.apps.containers',
        route: '/containers/container',
        glyph: 'containers',
        tint: ICON_PALETTES.containers.tint,
    },
    {
        id: 'images',
        titleKey: 'desktop.apps.images',
        route: '/containers/image',
        glyph: 'images',
        tint: ICON_PALETTES.images.tint,
    },
    {
        id: 'appStore',
        titleKey: 'desktop.apps.appStore',
        route: '/apps/all',
        glyph: 'appStore',
        tint: ICON_PALETTES.appStore.tint,
    },
    {
        id: 'monitor',
        titleKey: 'desktop.apps.monitor',
        route: '/hosts/monitor/monitor',
        glyph: 'monitor',
        tint: ICON_PALETTES.monitor.tint,
    },
    {
        id: 'terminal',
        titleKey: 'desktop.apps.terminal',
        route: '/terminal',
        glyph: 'terminal',
        tint: ICON_PALETTES.terminal.tint,
    },
    {
        id: 'files',
        titleKey: 'desktop.apps.files',
        route: '/hosts/files',
        glyph: 'files',
        tint: ICON_PALETTES.files.tint,
    },
    {
        id: 'database',
        titleKey: 'desktop.apps.database',
        route: '/databases/mysql',
        glyph: 'database',
        tint: ICON_PALETTES.database.tint,
    },
    {
        id: 'website',
        titleKey: 'desktop.apps.website',
        route: '/websites',
        glyph: 'website',
        tint: ICON_PALETTES.website.tint,
    },
    {
        id: 'firewall',
        titleKey: 'desktop.apps.firewall',
        route: '/hosts/firewall/rules',
        glyph: 'firewall',
        tint: ICON_PALETTES.firewall.tint,
    },
    {
        id: 'cronjob',
        titleKey: 'desktop.apps.cronjob',
        route: '/cronjobs/cronjob',
        glyph: 'cronjob',
        tint: ICON_PALETTES.cronjob.tint,
    },
    {
        id: 'logs',
        titleKey: 'desktop.apps.logs',
        route: '/logs/operation',
        glyph: 'logs',
        tint: ICON_PALETTES.logs.tint,
    },
    {
        id: 'toolbox',
        titleKey: 'desktop.apps.toolbox',
        route: '/toolbox/device',
        glyph: 'toolbox',
        tint: ICON_PALETTES.toolbox.tint,
    },
    {
        id: 'gpu',
        titleKey: 'desktop.apps.gpu',
        route: '/ai/gpu/current',
        glyph: 'gpu',
        tint: ICON_PALETTES.gpu.tint,
    },
    {
        id: 'homebrew',
        titleKey: 'desktop.apps.homebrew',
        route: '/homebrew/index',
        glyph: 'homebrew',
        tint: ICON_PALETTES.homebrew.tint,
        platformFeature: 'homebrew',
    },
    {
        id: 'settings',
        titleKey: 'desktop.apps.settings',
        route: '/settings',
        glyph: 'settings',
        tint: ICON_PALETTES.settings.tint,
        dividerBefore: true,
    },
];

export const DESKTOP_ICON_IDS = ['containers', 'terminal', 'files', 'settings'];

export const findDesktopApp = (id: string) => DESKTOP_APPS.find((item) => item.id === id);
