import i18n from '@/lang';
import { Container } from '@/api/interface/container';
import { dateFormatWithoutYear } from '@/utils/date';

export interface ContainerChartOption {
    title: string;
    xData: Array<string>;
    yData: Array<{ name: string; data: Array<number> }>;
    formatStr: string;
}

export interface ContainerHistoryCharts {
    cpuChart: ContainerChartOption;
    memoryChart: ContainerChartOption;
    ioChart: ContainerChartOption;
    networkChart: ContainerChartOption;
}

const historyWindowSeconds = 3600;

export function defaultContainerHistoryRange(): [Date, Date] {
    const end = new Date();
    const start = new Date(end.getTime() - historyWindowSeconds * 1000);
    return [start, end];
}

function toFixedNumber(value: number): number {
    if (value === undefined || value === null || Number.isNaN(value)) {
        return 0;
    }
    return Number(Number(value).toFixed(2));
}

export function buildContainerHistoryCharts(data: Container.ContainerMonitorData): ContainerHistoryCharts {
    const points = data?.points ?? [];
    const xData = points.map((item) => dateFormatWithoutYear(item.createdAt));
    const series = (key: keyof Container.ContainerMonitorPoint) =>
        points.map((item) => toFixedNumber(item[key] as number));

    return {
        cpuChart: {
            title: 'CPU',
            xData,
            yData: [{ name: 'CPU', data: series('cpuPercent') }],
            formatStr: '%',
        },
        memoryChart: {
            title: i18n.global.t('monitor.memory'),
            xData,
            yData: [
                { name: i18n.global.t('monitor.memory'), data: series('memory') },
                { name: i18n.global.t('container.cache'), data: series('cache') },
            ],
            formatStr: 'MB',
        },
        ioChart: {
            title: i18n.global.t('monitor.disk') + ' IO',
            xData,
            yData: [
                { name: i18n.global.t('monitor.read'), data: series('ioRead') },
                { name: i18n.global.t('monitor.write'), data: series('ioWrite') },
            ],
            formatStr: 'KB/s',
        },
        networkChart: {
            title: i18n.global.t('monitor.network'),
            xData,
            yData: [
                { name: i18n.global.t('monitor.up'), data: series('networkTX') },
                { name: i18n.global.t('monitor.down'), data: series('networkRX') },
            ],
            formatStr: 'KB/s',
        },
    };
}
