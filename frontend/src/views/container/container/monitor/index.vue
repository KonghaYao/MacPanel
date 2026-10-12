<template>
    <DrawerPro
        v-model="monitorVisible"
        :header="$t('menu.monitor')"
        @close="handleClose"
        :resource="title"
        size="large"
    >
        <el-tabs v-model="activeTab" @tab-change="onTabChange">
            <el-tab-pane :label="$t('monitor.live')" name="live">
                <el-alert
                    v-if="!isRunning"
                    :title="$t('monitor.containerStopped')"
                    type="info"
                    :closable="false"
                    show-icon
                    class="mb-2.5"
                />
                <el-form label-position="top" @submit.prevent>
                    <el-form-item :label="$t('container.refreshTime')">
                        <el-select v-model="timeInterval" @change="changeTimer" :disabled="!isRunning">
                            <el-option label="3s" :value="3" />
                            <el-option label="5s" :value="5" />
                            <el-option label="10s" :value="10" />
                            <el-option label="30s" :value="30" />
                            <el-option label="60s" :value="60" />
                        </el-select>
                    </el-form-item>
                </el-form>
                <el-card>
                    <v-charts
                        height="200px"
                        id="cpuChart"
                        type="line"
                        renderer="svg"
                        :option="chartsOption['cpuChart']"
                        v-if="chartsOption['cpuChart']"
                    />
                </el-card>
                <el-card class="mt-2.5">
                    <v-charts
                        height="200px"
                        id="memoryChart"
                        type="line"
                        renderer="svg"
                        :option="chartsOption['memoryChart']"
                        v-if="chartsOption['memoryChart']"
                    />
                </el-card>
                <el-card class="mt-2.5">
                    <v-charts
                        height="200px"
                        id="ioChart"
                        type="line"
                        renderer="svg"
                        :option="chartsOption['ioChart']"
                        v-if="chartsOption['ioChart']"
                    />
                </el-card>
                <el-card class="mt-2.5">
                    <v-charts
                        height="200px"
                        id="networkChart"
                        type="line"
                        renderer="svg"
                        :option="chartsOption['networkChart']"
                        v-if="chartsOption['networkChart']"
                    />
                </el-card>
            </el-tab-pane>
            <el-tab-pane :label="$t('monitor.history')" name="history">
                <el-form label-position="top" @submit.prevent>
                    <el-form-item :label="$t('monitor.timeRange')">
                        <div class="history-toolbar">
                            <el-date-picker
                                v-model="timeRange"
                                type="datetimerange"
                                range-separator="-"
                                :start-placeholder="$t('commons.search.timeStart')"
                                :end-placeholder="$t('commons.search.timeEnd')"
                                :shortcuts="shortcuts"
                                @change="loadHistory"
                            />
                            <el-select v-model="aggregation" @change="loadHistory" class="history-aggregation">
                                <el-option :label="$t('monitor.avg')" value="avg" />
                                <el-option :label="$t('monitor.max')" value="max" />
                            </el-select>
                            <TableRefresh @search="loadHistory" />
                            <el-button v-permission @click="onClean" icon="Delete">
                                {{ $t('monitor.cleanMonitor') }}
                            </el-button>
                        </div>
                    </el-form-item>
                </el-form>
                <div v-loading="historyLoading">
                    <el-empty v-if="!historyLoading && historyEmpty" :description="$t('monitor.historyEmpty')" />
                    <template v-else>
                        <el-card>
                            <v-charts
                                height="200px"
                                id="historyCpuChart"
                                type="line"
                                renderer="svg"
                                :dataZoom="true"
                                :option="historyCharts['cpuChart']"
                                v-if="historyCharts['cpuChart']"
                            />
                        </el-card>
                        <el-card class="mt-2.5">
                            <v-charts
                                height="200px"
                                id="historyMemoryChart"
                                type="line"
                                renderer="svg"
                                :dataZoom="true"
                                :option="historyCharts['memoryChart']"
                                v-if="historyCharts['memoryChart']"
                            />
                        </el-card>
                        <el-card class="mt-2.5">
                            <v-charts
                                height="200px"
                                id="historyIoChart"
                                type="line"
                                renderer="svg"
                                :dataZoom="true"
                                :option="historyCharts['ioChart']"
                                v-if="historyCharts['ioChart']"
                            />
                        </el-card>
                        <el-card class="mt-2.5">
                            <v-charts
                                height="200px"
                                id="historyNetworkChart"
                                type="line"
                                renderer="svg"
                                :dataZoom="true"
                                :option="historyCharts['networkChart']"
                                v-if="historyCharts['networkChart']"
                            />
                        </el-card>
                    </template>
                </div>
            </el-tab-pane>
        </el-tabs>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, ref } from 'vue';
import { containerStats, loadContainerMonitor, cleanContainerMonitor } from '@/api/modules/container';
import { dateFormatForSecond } from '@/utils/date';
import { shortcuts } from '@/utils/shortcuts';
import { MsgSuccess } from '@/utils/message';
import { buildContainerHistoryCharts, defaultContainerHistoryRange } from '@/utils/container-monitor';
import VCharts from '@/components/v-charts/index.vue';
import TableRefresh from '@/components/table/TableRefresh.vue';
import i18n from '@/lang';

const title = ref();
const monitorVisible = ref(false);
const activeTab = ref('live');
const isRunning = ref(true);
const timeInterval = ref();
let timer: NodeJS.Timer | null = null;
let isInit = ref<boolean>(true);
interface DialogProps {
    containerID: string;
    container: string;
    state?: string;
}
const dialogData = ref<DialogProps>({
    containerID: '',
    container: '',
    state: 'running',
});

const acceptParams = async (params: DialogProps): Promise<void> => {
    monitorVisible.value = true;
    dialogData.value.containerID = params.containerID;
    dialogData.value.container = params.container;
    dialogData.value.state = params.state ?? 'running';
    isRunning.value = dialogData.value.state === 'running';
    title.value = params.container;
    cpuDatas.value = [];
    memDatas.value = [];
    cacheDatas.value = [];
    ioReadDatas.value = [];
    ioWriteDatas.value = [];
    netTxDatas.value = [];
    netRxDatas.value = [];
    timeDatas.value = [];
    timeInterval.value = 5;
    isInit.value = true;
    historyLoaded.value = false;
    historyCharts.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
    timeRange.value = defaultContainerHistoryRange();
    aggregation.value = 'avg';
    activeTab.value = isRunning.value ? 'live' : 'history';
    if (isRunning.value) {
        loadData();
        timer = setInterval(async () => {
            if (monitorVisible.value && activeTab.value === 'live') {
                isInit.value = false;
                loadData();
            }
        }, 1000 * timeInterval.value);
    }
    if (activeTab.value === 'history') {
        loadHistory();
    }
};

const cpuDatas = ref<Array<string>>([]);
const memDatas = ref<Array<string>>([]);
const cacheDatas = ref<Array<string>>([]);
const ioReadDatas = ref<Array<string>>([]);
const ioWriteDatas = ref<Array<string>>([]);
const netTxDatas = ref<Array<string>>([]);
const netRxDatas = ref<Array<string>>([]);
const timeDatas = ref<Array<string>>([]);
const chartsOption = ref({ cpuChart: null, memoryChart: null, ioChart: null, networkChart: null });

const timeRange = ref<[Date, Date]>(defaultContainerHistoryRange());
const aggregation = ref('avg');
const historyLoading = ref(false);
const historyLoaded = ref(false);
const historyEmpty = ref(false);
const historyCharts = ref({ cpuChart: null, memoryChart: null, ioChart: null, networkChart: null });

const onTabChange = (name: string | number) => {
    if (name === 'history' && !historyLoaded.value && !historyLoading.value) {
        loadHistory();
    }
};

const changeTimer = () => {
    clearInterval(Number(timer));
    timer = setInterval(async () => {
        if (monitorVisible.value && activeTab.value === 'live') {
            loadData();
        }
    }, 1000 * timeInterval.value);
};

const loadData = async () => {
    const res = await containerStats(dialogData.value.containerID);
    cpuDatas.value.push(res.data.cpuPercent.toFixed(2));
    if (cpuDatas.value.length > 20) {
        cpuDatas.value.splice(0, 1);
    }
    memDatas.value.push(res.data.memory.toFixed(2));
    if (memDatas.value.length > 20) {
        memDatas.value.splice(0, 1);
    }
    cacheDatas.value.push(res.data.cache.toFixed(2));
    if (cacheDatas.value.length > 20) {
        cacheDatas.value.splice(0, 1);
    }
    ioReadDatas.value.push(res.data.ioRead.toFixed(2));
    if (ioReadDatas.value.length > 20) {
        ioReadDatas.value.splice(0, 1);
    }
    ioWriteDatas.value.push(res.data.ioWrite.toFixed(2));
    if (ioWriteDatas.value.length > 20) {
        ioWriteDatas.value.splice(0, 1);
    }
    netTxDatas.value.push(res.data.networkTX.toFixed(2));
    if (netTxDatas.value.length > 20) {
        netTxDatas.value.splice(0, 1);
    }
    netRxDatas.value.push(res.data.networkRX.toFixed(2));
    if (netRxDatas.value.length > 20) {
        netRxDatas.value.splice(0, 1);
    }
    timeDatas.value.push(dateFormatForSecond(res.data.shotTime));
    if (timeDatas.value.length > 20) {
        timeDatas.value.splice(0, 1);
    }

    chartsOption.value['cpuChart'] = {
        title: 'CPU',
        xData: timeDatas.value,
        yData: [
            {
                name: 'CPU',
                data: cpuDatas.value,
            },
        ],
        formatStr: '%',
    };

    chartsOption.value['memoryChart'] = {
        title: i18n.global.t('monitor.memory'),
        xData: timeDatas.value,
        yData: [
            {
                name: i18n.global.t('monitor.memory'),
                data: memDatas.value,
            },
            {
                name: i18n.global.t('container.cache'),
                data: cacheDatas.value,
            },
        ],
        formatStr: 'MB',
    };

    chartsOption.value['ioChart'] = {
        title: i18n.global.t('monitor.disk') + ' IO',
        xData: timeDatas.value,
        yData: [
            {
                name: i18n.global.t('monitor.read'),
                data: ioReadDatas.value,
            },
            {
                name: i18n.global.t('monitor.write'),
                data: ioWriteDatas.value,
            },
        ],
        formatStr: 'MB',
    };

    chartsOption.value['networkChart'] = {
        title: i18n.global.t('monitor.network'),
        xData: timeDatas.value,
        yData: [
            {
                name: i18n.global.t('monitor.up'),
                data: netTxDatas.value,
            },
            {
                name: i18n.global.t('monitor.down'),
                data: netRxDatas.value,
            },
        ],
        formatStr: 'KB',
    };
};

const loadHistory = async () => {
    if (!timeRange.value || timeRange.value.length !== 2) {
        return;
    }
    historyLoading.value = true;
    historyEmpty.value = false;
    try {
        const res = await loadContainerMonitor({
            containerName: dialogData.value.container,
            startTime: timeRange.value[0],
            endTime: timeRange.value[1],
            aggregation: aggregation.value,
        });
        historyLoaded.value = true;
        if (!res.data.points || res.data.points.length === 0 || res.data.sampleCount === 0) {
            historyEmpty.value = true;
            historyCharts.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
            return;
        }
        historyCharts.value = buildContainerHistoryCharts(res.data);
    } catch (error) {
        historyEmpty.value = true;
        historyCharts.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
    } finally {
        historyLoading.value = false;
    }
};

const onClean = async () => {
    ElMessageBox.confirm(
        i18n.global.t('monitor.cleanHelper', [dialogData.value.container]),
        i18n.global.t('monitor.cleanMonitor'),
        {
            confirmButtonText: i18n.global.t('commons.button.confirm'),
            cancelButtonText: i18n.global.t('commons.button.cancel'),
            type: 'info',
        },
    ).then(async () => {
        historyLoading.value = true;
        await cleanContainerMonitor({ containerName: dialogData.value.container })
            .then(() => {
                MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
                loadHistory();
            })
            .catch(() => {
                historyLoading.value = false;
            });
    });
};

const handleClose = async () => {
    monitorVisible.value = false;
    clearInterval(Number(timer));
    timer = null;
    activeTab.value = 'live';
    chartsOption.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
    historyCharts.value = { cpuChart: null, memoryChart: null, ioChart: null, networkChart: null };
    historyLoaded.value = false;
    historyEmpty.value = false;
};

onBeforeUnmount(() => {
    handleClose;
});

defineExpose({
    acceptParams,
});
</script>

<style scoped lang="scss">
.history-toolbar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
}

.history-aggregation {
    width: 140px;
}
</style>
