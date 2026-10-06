<template>
    <div v-loading="loading">
        <NoSuchService v-if="statusLoaded && !status.isExist" name="Homebrew" />

        <template v-else-if="status.isExist">
            <div class="app-status card-interval">
                <el-card>
                    <div class="flex w-full flex-col gap-4 md:flex-row">
                        <div class="flex flex-wrap gap-4 ml-3">
                            <el-tag effect="dark" type="success">Homebrew</el-tag>
                            <el-tag v-if="status.version">{{ status.version }}</el-tag>
                            <el-tag v-if="status.prefix">{{ status.prefix }}</el-tag>
                        </div>
                        <div class="flex flex-wrap gap-2 ml-3 text-sm text-gray-500">
                            <span>{{ $t('homebrew.formulaCount', [status.formulaCount]) }}</span>
                            <span>{{ $t('homebrew.caskCount', [status.caskCount]) }}</span>
                        </div>
                    </div>
                </el-card>
            </div>

            <LayoutContent :title="$t('menu.homebrew')">
                <template #main>
                    <el-tabs v-model="activeTab">
                        <el-tab-pane :label="$t('homebrew.tabs.installed')" name="installed">
                            <div class="flex flex-wrap gap-2 mb-4">
                                <el-select v-model="listType" class="w-40" @change="searchInstalled">
                                    <el-option :label="$t('homebrew.typeAll')" value="all" />
                                    <el-option :label="$t('homebrew.typeFormula')" value="formula" />
                                    <el-option :label="$t('homebrew.typeCask')" value="cask" />
                                </el-select>
                                <TableSearch @search="searchInstalled" v-model:searchName="searchName" />
                                <TableRefresh @search="searchInstalled" />
                            </div>
                            <ComplexTable
                                :pagination-config="paginationConfig"
                                @search="searchInstalled"
                                :data="installedData"
                            >
                                <el-table-column :label="$t('commons.table.name')" prop="name" min-width="160" />
                                <el-table-column :label="$t('commons.table.version')" prop="version" min-width="120" />
                                <el-table-column :label="$t('commons.table.type')" prop="type" min-width="100" />
                                <fu-table-operations :buttons="installedButtons" :label="$t('commons.table.operate')" fix />
                            </ComplexTable>
                        </el-tab-pane>

                        <el-tab-pane :label="$t('homebrew.tabs.search')" name="search">
                            <div class="flex flex-wrap gap-2 mb-4">
                                <el-input
                                    v-model="searchQuery"
                                    class="w-64"
                                    :placeholder="$t('homebrew.searchPlaceholder')"
                                    @keyup.enter="onSearchPackages"
                                />
                                <el-select v-model="searchType" class="w-40">
                                    <el-option :label="$t('homebrew.typeAll')" value="all" />
                                    <el-option :label="$t('homebrew.typeFormula')" value="formula" />
                                    <el-option :label="$t('homebrew.typeCask')" value="cask" />
                                </el-select>
                                <el-button type="primary" @click="onSearchPackages">
                                    {{ $t('commons.button.search') }}
                                </el-button>
                            </div>
                            <ComplexTable :data="searchData" :pagination-config="searchPagination">
                                <el-table-column :label="$t('commons.table.name')" prop="name" min-width="160" />
                                <el-table-column :label="$t('commons.table.type')" prop="type" min-width="100" />
                                <fu-table-operations :buttons="searchButtons" :label="$t('commons.table.operate')" fix />
                            </ComplexTable>
                        </el-tab-pane>

                        <el-tab-pane :label="$t('homebrew.tabs.mirror')" name="mirror">
                            <el-alert type="info" :closable="false" class="mb-4">
                                <template #title>{{ $t('homebrew.mirrorHelper') }}</template>
                            </el-alert>
                            <div class="flex flex-wrap gap-2 mb-4">
                                <el-button
                                    v-for="preset in mirrorPresets"
                                    :key="preset.key"
                                    @click="applyMirrorPreset(preset.key)"
                                >
                                    {{ preset.label }}
                                </el-button>
                                <el-button @click="clearMirror">{{ $t('homebrew.clearMirror') }}</el-button>
                            </div>
                            <el-form :model="mirrorForm" label-width="180px" class="max-w-3xl">
                                <el-form-item :label="$t('homebrew.bottleDomain')">
                                    <el-input v-model="mirrorForm.bottleDomain" :placeholder="$t('homebrew.bottleDomainTip')" />
                                </el-form-item>
                                <el-form-item :label="$t('homebrew.apiDomain')">
                                    <el-input v-model="mirrorForm.apiDomain" :placeholder="$t('homebrew.apiDomainTip')" />
                                </el-form-item>
                                <el-form-item :label="$t('homebrew.brewGitRemote')">
                                    <el-input v-model="mirrorForm.brewGitRemote" :placeholder="$t('homebrew.brewGitRemoteTip')" />
                                </el-form-item>
                                <el-form-item :label="$t('homebrew.coreGitRemote')">
                                    <el-input v-model="mirrorForm.coreGitRemote" :placeholder="$t('homebrew.coreGitRemoteTip')" />
                                </el-form-item>
                                <el-form-item :label="$t('homebrew.caskGitRemote')">
                                    <el-input v-model="mirrorForm.caskGitRemote" :placeholder="$t('homebrew.caskGitRemoteTip')" />
                                </el-form-item>
                                <el-form-item>
                                    <el-button v-permission v-node-admin type="primary" @click="saveMirror">
                                        {{ $t('commons.button.save') }}
                                    </el-button>
                                </el-form-item>
                            </el-form>
                        </el-tab-pane>

                        <el-tab-pane :label="$t('homebrew.tabs.maintenance')" name="maintenance">
                            <div class="flex flex-wrap gap-2 mb-4">
                                <el-button v-permission v-node-admin type="primary" @click="onUpdate">
                                    {{ $t('homebrew.brewUpdate') }}
                                </el-button>
                                <el-button v-permission v-node-admin type="primary" plain @click="onUpgradeAll">
                                    {{ $t('homebrew.upgradeAll') }}
                                </el-button>
                                <el-button v-permission v-node-admin @click="onDoctor">
                                    {{ $t('homebrew.brewDoctor') }}
                                </el-button>
                            </div>
                            <el-input
                                v-if="doctorOutput"
                                v-model="doctorOutput"
                                type="textarea"
                                :rows="18"
                                readonly
                            />
                        </el-tab-pane>
                    </el-tabs>
                </template>
            </LayoutContent>
        </template>

        <TaskLog ref="taskLogRef" width="70%" />
        <OpDialog ref="opRef" @search="searchInstalled" />
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import i18n from '@/lang';
import NoSuchService from '@/components/layout-content/no-such-service.vue';
import TaskLog from '@/components/log/task/index.vue';
import {
    doctorHomebrew,
    getHomebrewMirror,
    getHomebrewStatus,
    installHomebrew,
    listHomebrew,
    searchHomebrew,
    uninstallHomebrew,
    updateHomebrew,
    updateHomebrewMirror,
    upgradeHomebrew,
} from '@/api/modules/homebrew';
import { Homebrew } from '@/api/interface/homebrew';
import { newUUID } from '@/utils/id';
import { MsgSuccess } from '@/utils/message';

const loading = ref(false);
const statusLoaded = ref(false);
const status = ref<Homebrew.Status>({
    isExist: false,
    version: '',
    prefix: '',
    formulaCount: 0,
    caskCount: 0,
    brewPath: '',
    cellarPath: '',
});

const activeTab = ref('installed');
const listType = ref('all');
const searchName = ref('');
const installedData = ref<Homebrew.Package[]>([]);
const paginationConfig = reactive({
    cacheSize: 20,
    currentPage: 1,
    pageSize: 20,
    total: 0,
});

const searchQuery = ref('');
const searchType = ref('all');
const searchData = ref<Homebrew.SearchResult[]>([]);
const searchPagination = reactive({
    cacheSize: 20,
    currentPage: 1,
    pageSize: 20,
    total: 0,
});

const mirrorForm = reactive<Homebrew.MirrorConfig>({
    bottleDomain: '',
    apiDomain: '',
    brewGitRemote: '',
    coreGitRemote: '',
    caskGitRemote: '',
});

const doctorOutput = ref('');
const taskLogRef = ref();
const opRef = ref();

const mirrorPresets = [
    { key: 'tuna', label: i18n.global.t('homebrew.presetTuna') },
    { key: 'ustc', label: i18n.global.t('homebrew.presetUstc') },
    { key: 'aliyun', label: i18n.global.t('homebrew.presetAliyun') },
];

const mirrorPresetValues: Record<string, Homebrew.MirrorConfig> = {
    tuna: {
        bottleDomain: 'https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles',
        apiDomain: 'https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api',
        brewGitRemote: 'https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git',
        coreGitRemote: 'https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git',
        caskGitRemote: 'https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-cask.git',
    },
    ustc: {
        bottleDomain: 'https://mirrors.ustc.edu.cn/homebrew-bottles',
        apiDomain: 'https://mirrors.ustc.edu.cn/homebrew-bottles/api',
        brewGitRemote: 'https://mirrors.ustc.edu.cn/brew.git',
        coreGitRemote: 'https://mirrors.ustc.edu.cn/homebrew-core.git',
        caskGitRemote: 'https://mirrors.ustc.edu.cn/homebrew-cask.git',
    },
    aliyun: {
        bottleDomain: 'https://mirrors.aliyun.com/homebrew/homebrew-bottles',
        apiDomain: 'https://mirrors.aliyun.com/homebrew/homebrew-bottles/api',
        brewGitRemote: 'https://mirrors.aliyun.com/homebrew/brew.git',
        coreGitRemote: 'https://mirrors.aliyun.com/homebrew/homebrew-core.git',
        caskGitRemote: 'https://mirrors.aliyun.com/homebrew/homebrew-cask.git',
    },
};

const openTaskLog = (taskID: string) => {
    taskLogRef.value.openWithTaskID(taskID);
};

const loadStatus = async () => {
    loading.value = true;
    try {
        const res = await getHomebrewStatus();
        status.value = res.data;
        statusLoaded.value = true;
    } finally {
        loading.value = false;
    }
};

const searchInstalled = async () => {
    if (!status.value.isExist) {
        return;
    }
    loading.value = true;
    try {
        const res = await listHomebrew({
            page: paginationConfig.currentPage,
            pageSize: paginationConfig.pageSize,
            info: searchName.value,
            type: listType.value,
        });
        installedData.value = res.data.items || [];
        paginationConfig.total = res.data.total || 0;
    } finally {
        loading.value = false;
    }
};

const onSearchPackages = async () => {
    if (!searchQuery.value.trim()) {
        return;
    }
    loading.value = true;
    try {
        const res = await searchHomebrew({
            query: searchQuery.value.trim(),
            type: searchType.value,
        });
        searchData.value = res.data || [];
        searchPagination.total = searchData.value.length;
    } finally {
        loading.value = false;
    }
};

const runPackageTask = async (
    action: 'install' | 'uninstall' | 'upgrade',
    row: { name: string; type: string },
) => {
    const taskID = newUUID();
    const params = { name: row.name, type: row.type, taskID };
    if (action === 'install') {
        await installHomebrew(params);
    } else if (action === 'uninstall') {
        await uninstallHomebrew(params);
    } else {
        await upgradeHomebrew(params);
    }
    MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
    openTaskLog(taskID);
};

const installedButtons = [
    {
        label: i18n.global.t('commons.button.upgrade'),
        permission: true,
        nodeAdmin: true,
        click: (row: Homebrew.Package) => runPackageTask('upgrade', row),
    },
    {
        label: i18n.global.t('commons.button.uninstall'),
        permission: true,
        nodeAdmin: true,
        click: (row: Homebrew.Package) => {
            opRef.value.acceptParams({
                title: i18n.global.t('commons.button.uninstall'),
                names: [row.name],
                msg: i18n.global.t('commons.msg.operatorHelper', [
                    i18n.global.t('commons.button.uninstall'),
                    row.name,
                ]),
                api: null,
                params: row,
                submit: () => runPackageTask('uninstall', row),
            });
        },
    },
];

const searchButtons = [
    {
        label: i18n.global.t('commons.button.install'),
        permission: true,
        nodeAdmin: true,
        click: (row: Homebrew.SearchResult) => runPackageTask('install', row),
    },
];

const loadMirror = async () => {
    const res = await getHomebrewMirror();
    Object.assign(mirrorForm, res.data);
};

const applyMirrorPreset = (key: string) => {
    Object.assign(mirrorForm, mirrorPresetValues[key]);
};

const clearMirror = () => {
    Object.assign(mirrorForm, {
        bottleDomain: '',
        apiDomain: '',
        brewGitRemote: '',
        coreGitRemote: '',
        caskGitRemote: '',
    });
};

const saveMirror = async () => {
    await updateHomebrewMirror({ ...mirrorForm });
    MsgSuccess(i18n.global.t('commons.msg.saveSuccess'));
};

const onUpdate = async () => {
    const taskID = newUUID();
    await updateHomebrew({ taskID });
    MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
    openTaskLog(taskID);
};

const onUpgradeAll = async () => {
    const taskID = newUUID();
    await upgradeHomebrew({ name: '', type: 'formula', taskID });
    MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
    openTaskLog(taskID);
};

const onDoctor = async () => {
    loading.value = true;
    try {
        const res = await doctorHomebrew();
        doctorOutput.value = res.data.output || '';
    } finally {
        loading.value = false;
    }
};

onMounted(async () => {
    await loadStatus();
    if (status.value.isExist) {
        await Promise.all([searchInstalled(), loadMirror()]);
    }
});
</script>
