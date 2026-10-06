import type { Component } from 'vue';
import type { DockGlyph } from '@/views/desktop/apps';
import ContainersIcon from '@/views/desktop/icons/ContainersIcon.vue';
import ImagesIcon from '@/views/desktop/icons/ImagesIcon.vue';
import AppStoreIcon from '@/views/desktop/icons/AppStoreIcon.vue';
import MonitorIcon from '@/views/desktop/icons/MonitorIcon.vue';
import TerminalIcon from '@/views/desktop/icons/TerminalIcon.vue';
import FilesIcon from '@/views/desktop/icons/FilesIcon.vue';
import DatabaseIcon from '@/views/desktop/icons/DatabaseIcon.vue';
import WebsiteIcon from '@/views/desktop/icons/WebsiteIcon.vue';
import FirewallIcon from '@/views/desktop/icons/FirewallIcon.vue';
import CronjobIcon from '@/views/desktop/icons/CronjobIcon.vue';
import LogsIcon from '@/views/desktop/icons/LogsIcon.vue';
import ToolboxIcon from '@/views/desktop/icons/ToolboxIcon.vue';
import GpuIcon from '@/views/desktop/icons/GpuIcon.vue';
import HomebrewIcon from '@/views/desktop/icons/HomebrewIcon.vue';
import SettingsIcon from '@/views/desktop/icons/SettingsIcon.vue';

export const DOCK_ICON_MAP: Record<DockGlyph, Component> = {
    containers: ContainersIcon,
    images: ImagesIcon,
    appStore: AppStoreIcon,
    monitor: MonitorIcon,
    terminal: TerminalIcon,
    files: FilesIcon,
    database: DatabaseIcon,
    website: WebsiteIcon,
    firewall: FirewallIcon,
    cronjob: CronjobIcon,
    logs: LogsIcon,
    toolbox: ToolboxIcon,
    gpu: GpuIcon,
    homebrew: HomebrewIcon,
    settings: SettingsIcon,
};
