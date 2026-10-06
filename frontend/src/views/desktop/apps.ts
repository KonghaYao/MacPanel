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
        tint: 'linear-gradient(180deg, #64b5ff 0%, #1d4ed8 100%)',
    },
    {
        id: 'images',
        titleKey: 'desktop.apps.images',
        route: '/containers/image',
        glyph: 'images',
        tint: 'linear-gradient(180deg, #d8b4fe 0%, #7c3aed 100%)',
    },
    {
        id: 'appStore',
        titleKey: 'desktop.apps.appStore',
        route: '/apps/all',
        glyph: 'appStore',
        tint: 'linear-gradient(180deg, #7dd3fc 0%, #2563eb 100%)',
    },
    {
        id: 'monitor',
        titleKey: 'desktop.apps.monitor',
        route: '/hosts/monitor/monitor',
        glyph: 'monitor',
        tint: 'linear-gradient(180deg, #6b7280 0%, #111827 100%)',
    },
    {
        id: 'terminal',
        titleKey: 'desktop.apps.terminal',
        route: '/terminal',
        glyph: 'terminal',
        tint: 'linear-gradient(180deg, #3f3f46 0%, #09090b 100%)',
    },
    {
        id: 'files',
        titleKey: 'desktop.apps.files',
        route: '/hosts/files',
        glyph: 'files',
        tint: 'linear-gradient(180deg, #93c5fd 0%, #1d4ed8 100%)',
    },
    {
        id: 'database',
        titleKey: 'desktop.apps.database',
        route: '/databases/mysql',
        glyph: 'database',
        tint: 'linear-gradient(180deg, #6ee7b7 0%, #047857 100%)',
    },
    {
        id: 'website',
        titleKey: 'desktop.apps.website',
        route: '/websites',
        glyph: 'website',
        tint: 'linear-gradient(180deg, #67e8f9 0%, #0369a1 100%)',
    },
    {
        id: 'firewall',
        titleKey: 'desktop.apps.firewall',
        route: '/hosts/firewall/rules',
        glyph: 'firewall',
        tint: 'linear-gradient(180deg, #fda4af 0%, #be123c 100%)',
    },
    {
        id: 'cronjob',
        titleKey: 'desktop.apps.cronjob',
        route: '/cronjobs/cronjob',
        glyph: 'cronjob',
        tint: 'linear-gradient(180deg, #fcd34d 0%, #b45309 100%)',
    },
    {
        id: 'logs',
        titleKey: 'desktop.apps.logs',
        route: '/logs/operation',
        glyph: 'logs',
        tint: 'linear-gradient(180deg, #d6d3d1 0%, #57534e 100%)',
    },
    {
        id: 'toolbox',
        titleKey: 'desktop.apps.toolbox',
        route: '/toolbox/device',
        glyph: 'toolbox',
        tint: 'linear-gradient(180deg, #fdba74 0%, #c2410c 100%)',
    },
    {
        id: 'gpu',
        titleKey: 'desktop.apps.gpu',
        route: '/ai/gpu/current',
        glyph: 'gpu',
        tint: 'linear-gradient(180deg, #c4b5fd 0%, #6d28d9 100%)',
    },
    {
        id: 'homebrew',
        titleKey: 'desktop.apps.homebrew',
        route: '/homebrew/index',
        glyph: 'homebrew',
        tint: 'linear-gradient(180deg, #fbbf24 0%, #b45309 100%)',
        platformFeature: 'homebrew',
    },
    {
        id: 'settings',
        titleKey: 'desktop.apps.settings',
        route: '/settings',
        glyph: 'settings',
        tint: 'linear-gradient(180deg, #d1d5db 0%, #4b5563 100%)',
        dividerBefore: true,
    },
];

export const DESKTOP_ICON_IDS = ['containers', 'terminal', 'files', 'settings'];

export const findDesktopApp = (id: string) => DESKTOP_APPS.find((item) => item.id === id);
