<template>
    <div>
        <FireRouter />
        <LayoutContent :title="$t('ssh.hostManage', 2)">
            <template #leftToolBar>
                <el-button v-permission v-node-admin type="primary" @click="onOpenDialog">
                    {{ $t('commons.button.create') }}
                </el-button>
                <el-button v-permission v-node-admin plain :disabled="selects.length === 0" @click="onDelete(null)">
                    {{ $t('commons.button.delete') }}
                </el-button>
            </template>
            <template #rightToolBar>
                <TableSearch @search="search()" v-model:searchName="searchName" />
                <TableRefresh @search="search()" />
            </template>
            <template #main>
                <ComplexTable
                    :pagination-config="paginationConfig"
                    v-model:selects="selects"
                    :data="data"
                    v-loading="loading"
                    @search="search"
                    :heightDiff="260"
                >
                    <el-table-column type="selection" fix />
                    <el-table-column :label="$t('home.alias')" prop="alias" min-width="120" />
                    <el-table-column :label="$t('ssh.hostAddress')" prop="hostName" min-width="160" />
                    <el-table-column :label="$t('commons.table.user')" prop="user" min-width="100" />
                    <el-table-column :label="$t('ssh.authStatus')" prop="authStatus" min-width="120">
                        <template #default="{ row }">
                            <el-tag :type="row.authStatus === 'authorized' ? 'success' : 'danger'">
                                {{ authStatusLabel(row.authStatus) }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('commons.table.date')" min-width="160">
                        <template #default="{ row }">
                            {{ dateFormatWithoutSeconds(row.createdAt) }}
                        </template>
                    </el-table-column>
                    <fu-table-operations width="160px" :buttons="buttons" :label="$t('commons.table.operate')" fix />
                </ComplexTable>
            </template>
        </LayoutContent>

        <OpDialog ref="opRef" @search="search" @submit="onSubmitDelete()" />
        <Operate ref="dialogRef" @search="search" />
    </div>
</template>

<script setup lang="ts">
import FireRouter from '@/views/host/ssh/index.vue';
import Operate from '@/views/host/ssh/hosts/operate/index.vue';
import { deleteSSHHosts, searchSSHHosts, testSSHHost } from '@/api/modules/host';
import { Host } from '@/api/interface/host';
import i18n from '@/lang';
import { MsgError, MsgSuccess } from '@/utils/message';
import { dateFormatWithoutSeconds } from '@/utils/date';
import { reactive, ref } from 'vue';

const loading = ref(false);
const data = ref<Host.SSHHostInfo[]>([]);
const selects = ref<Host.SSHHostInfo[]>([]);
const searchName = ref('');
const dialogRef = ref();
const opRef = ref();
const operateIDs = ref<number[]>([]);

const paginationConfig = reactive({
    cacheSizeKey: 'ssh-host-page-size',
    currentPage: 1,
    pageSize: Number(localStorage.getItem('ssh-host-page-size')) || 20,
    total: 0,
    small: true,
});

const authStatusLabel = (status: string) => {
    if (status === 'authorized') {
        return i18n.global.t('ssh.authAuthorized');
    }
    return i18n.global.t('ssh.authFailed');
};

const buttons = [
    {
        label: i18n.global.t('ssh.testConnection'),
        permission: true,
        click: (row: Host.SSHHostInfo) => {
            onTest(row);
        },
    },
    {
        label: i18n.global.t('commons.button.delete'),
        permission: true,
        nodeAdmin: true,
        click: (row: Host.SSHHostInfo) => {
            onDelete(row);
        },
    },
];

const search = async () => {
    loading.value = true;
    try {
        const res = await searchSSHHosts({
            page: paginationConfig.currentPage,
            pageSize: paginationConfig.pageSize,
            info: searchName.value,
        });
        data.value = res.data.items || [];
        paginationConfig.total = res.data.total || 0;
    } catch (error) {
        MsgError(error);
    } finally {
        loading.value = false;
    }
};

const onOpenDialog = () => {
    dialogRef.value.acceptParams();
};

const onDelete = (row: Host.SSHHostInfo | null) => {
    operateIDs.value = row ? [row.id] : selects.value.map((item) => item.id);
    opRef.value.acceptParams({
        title: i18n.global.t('commons.button.delete'),
        names: row ? [row.alias] : selects.value.map((item) => item.alias),
        count: operateIDs.value.length,
    });
};

const onSubmitDelete = async () => {
    try {
        await deleteSSHHosts(operateIDs.value);
        MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        search();
    } catch (error) {
        MsgError(error);
    }
};

const onTest = async (row: Host.SSHHostInfo) => {
    loading.value = true;
    try {
        const res = await testSSHHost({ id: row.id });
        if (res.data.status) {
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        } else {
            MsgError(i18n.global.t('ssh.testFailed'));
        }
    } catch (error) {
        MsgError(error);
    } finally {
        loading.value = false;
    }
};

search();
</script>
