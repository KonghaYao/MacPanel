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
    },
    {
        id: 'images',
        titleKey: 'desktop.apps.images',
        route: '/containers/image',
        glyph: 'images',
    },
    {
        id: 'appStore',
        titleKey: 'desktop.apps.appStore',
        route: '/apps/all',
        glyph: 'appStore',
    },
    {
        id: 'monitor',
        titleKey: 'desktop.apps.monitor',
        route: '/hosts/monitor/monitor',
        glyph: 'monitor',
    },
    {
        id: 'terminal',
        titleKey: 'desktop.apps.terminal',
        route: '/terminal',
        glyph: 'terminal',
    },
    {
        id: 'files',
        titleKey: 'desktop.apps.files',
        route: '/hosts/files',
        glyph: 'files',
    },
    {
        id: 'database',
        titleKey: 'desktop.apps.database',
        route: '/databases/mysql',
        glyph: 'database',
    },
    {
        id: 'website',
        titleKey: 'desktop.apps.website',
        route: '/websites',
        glyph: 'website',
    },
    {
        id: 'firewall',
        titleKey: 'desktop.apps.firewall',
        route: '/hosts/firewall/rules',
        glyph: 'firewall',
    },
    {
        id: 'cronjob',
        titleKey: 'desktop.apps.cronjob',
        route: '/cronjobs/cronjob',
        glyph: 'cronjob',
    },
    {
        id: 'logs',
        titleKey: 'desktop.apps.logs',
        route: '/logs/operation',
        glyph: 'logs',
    },
    {
        id: 'toolbox',
        titleKey: 'desktop.apps.toolbox',
        route: '/toolbox/device',
        glyph: 'toolbox',
    },
    {
        id: 'gpu',
        titleKey: 'desktop.apps.gpu',
        route: '/ai/gpu/current',
        glyph: 'gpu',
    },
    {
        id: 'homebrew',
        titleKey: 'desktop.apps.homebrew',
        route: '/homebrew/index',
        glyph: 'homebrew',
        platformFeature: 'homebrew',
    },
    {
        id: 'settings',
        titleKey: 'desktop.apps.settings',
        route: '/settings',
        glyph: 'settings',
        dividerBefore: true,
    },
];

export const DESKTOP_ICON_IDS = ['containers', 'terminal', 'files', 'settings'];

export const findDesktopApp = (id: string) => DESKTOP_APPS.find((item) => item.id === id);
