<template>
    <div v-loading="loading">
        <LayoutContent :title="$t('database.s3.title')">
            <template #app v-if="currentConn?.source === 'rustfs' && currentConn.appName">
                <AppStatus
                    :key="currentConn.appInstallId"
                    app-key="rustfs"
                    :app-name="currentConn.appName"
                    :hide-setting="true"
                    v-model:loading="loading"
                    @after="onLoadBuckets"
                />
            </template>
            <template #leftToolBar>
                <el-button v-permission type="primary" :disabled="!currentConn || !currentBucket" @click="onUploadClick">
                    {{ $t('commons.button.upload') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn || !currentBucket" @click="onCreateFolder">
                    {{ $t('database.s3.createFolder') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn" @click="onCreateBucket">
                    {{ $t('database.s3.createBucket') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn || !currentBucket" @click="onDeleteBucket">
                    {{ $t('database.s3.deleteBucket') }}
                </el-button>
                <el-button
                    v-permission
                    :disabled="selectedRows.length === 0"
                    @click="onDeleteSelected"
                >
                    {{ $t('commons.button.delete') }}
                </el-button>
                <el-upload
                    ref="uploadRef"
                    class="hidden-upload"
                    :auto-upload="false"
                    :show-file-list="false"
                    :on-change="onUploadChange"
                />
            </template>
            <template #rightToolBar>
                <el-select v-if="connections.length" v-model="currentConnId" class="p-w-200" @change="onChangeConn">
                    <template #prefix>{{ $t('database.s3.connection') }}</template>
                    <el-option v-for="item in connections" :key="item.id" :label="item.name" :value="item.id">
                        <span>{{ item.name }}</span>
                        <el-tag class="ml-2" size="small">
                            {{
                                item.source === 'rustfs'
                                    ? $t('database.s3.sourceRustFS')
                                    : $t('database.s3.sourceManual')
                            }}
                        </el-tag>
                    </el-option>
                </el-select>
                <el-select
                    v-if="currentConn"
                    v-model="currentBucket"
                    class="p-w-200"
                    :placeholder="$t('database.s3.bucket')"
                    @change="onChangeBucket"
                >
                    <template #prefix>{{ $t('database.s3.bucket') }}</template>
                    <el-option v-for="item in buckets" :key="item.name" :label="item.name" :value="item.name" />
                </el-select>
                <el-button v-permission type="primary" plain @click="onConnectRustFS">
                    {{ $t('database.s3.connectRustFS') }}
                </el-button>
                <el-button v-permission type="primary" plain @click="onAddConn">
                    {{ $t('database.s3.addConnection') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn" @click="onEditConn">
                    {{ $t('commons.button.edit') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn" @click="onTestConn">
                    {{ $t('database.s3.testConnection') }}
                </el-button>
                <el-button v-permission :disabled="!currentConn" @click="onDeleteConn">
                    {{ $t('commons.button.delete') }}
                </el-button>
            </template>
            <template #main>
                <div v-if="isLoaded && !connections.length" class="app-warn">
                    <div class="flex flex-col gap-2 items-center justify-center w-full sm:flex-row">
                        <span v-if="!rustfs.installed">{{ $t('app.checkInstalledWarn', ['RustFS']) }}</span>
                        <span v-else>{{ $t('database.s3.selectConnection') }}</span>
                        <span
                            v-if="!rustfs.installed"
                            class="flex items-center justify-center gap-0.5 cursor-pointer"
                            @click="goInstallRustFS"
                        >
                            <el-icon><Position /></el-icon>
                            {{ $t('database.s3.goInstallRustFS') }}
                        </span>
                    </div>
                    <div>
                        <img src="@/assets/images/no_app.svg" />
                    </div>
                </div>
                <div v-else-if="currentConn">
                    <el-breadcrumb separator="/" class="mb-4">
                        <el-breadcrumb-item>
                            <el-link type="primary" @click="enterPrefix('')">{{ currentBucket || '/' }}</el-link>
                        </el-breadcrumb-item>
                        <el-breadcrumb-item v-for="item in prefixSegments" :key="item.prefix">
                            <el-link type="primary" @click="enterPrefix(item.prefix)">{{ item.name }}</el-link>
                        </el-breadcrumb-item>
                    </el-breadcrumb>
                    <ComplexTable :data="objects" v-model:selects="selectedRows">
                        <el-table-column type="selection" width="48" />
                        <el-table-column :label="$t('commons.table.name')" min-width="240" show-overflow-tooltip>
                            <template #default="{ row }">
                                <el-link v-if="row.prefix" type="primary" @click="enterPrefix(row.key)">
                                    {{ row.name }}/
                                </el-link>
                                <span v-else>{{ row.name }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="$t('file.size')" width="140">
                            <template #default="{ row }">
                                <span v-if="!row.prefix">{{ computeSize2(row.size || 0) }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column :label="$t('commons.table.date')" width="200">
                            <template #default="{ row }">
                                <span v-if="!row.prefix && row.lastModified">
                                    {{ dateFormatSimpleWithSecond(row.lastModified) }}
                                </span>
                            </template>
                        </el-table-column>
                        <fu-table-operations :buttons="buttons" :label="$t('commons.table.operate')" fix />
                    </ComplexTable>
                </div>
            </template>
        </LayoutContent>
        <ConnDrawer ref="connRef" @search="onConnSaved" />
        <DrawerPro v-model="previewVisible" :header="$t('database.s3.preview')" :resource="preview?.key || ''" size="large">
            <div v-if="preview">
                <el-alert
                    v-if="preview.tooLarge"
                    :title="$t('database.s3.previewTooLarge')"
                    type="warning"
                    :closable="false"
                    class="mb-4"
                />
                <el-alert
                    v-else-if="preview.kind === 'none'"
                    :title="$t('database.s3.previewUnsupported')"
                    type="info"
                    :closable="false"
                    class="mb-4"
                />
                <img
                    v-else-if="preview.kind === 'image' && preview.content"
                    class="max-w-full"
                    :src="'data:' + (preview.contentType || 'image/*') + ';base64,' + preview.content"
                />
                <pre v-else-if="preview.kind === 'text'" class="preview-text">{{ preview.content }}</pre>
            </div>
        </DrawerPro>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { ElMessageBox, UploadFile } from 'element-plus';
import AppStatus from '@/components/app-status/index.vue';
import ConnDrawer from '@/views/database/s3/conn/index.vue';
import { Database } from '@/api/interface/database';
import {
    connectS3RustFS,
    createS3Bucket,
    createS3Folder,
    deleteS3Bucket,
    deleteS3Connection,
    deleteS3Objects,
    downloadS3Object,
    listS3Buckets,
    listS3Connections,
    listS3Objects,
    listS3RustFS,
    previewS3Object,
    testS3Connection,
    uploadS3Object,
} from '@/api/modules/database';
import i18n from '@/lang';
import { MsgSuccess } from '@/utils/message';
import { routerToNameWithQuery } from '@/utils/router';
import { computeSize2 } from '@/utils/size';
import { dateFormatSimpleWithSecond } from '@/utils/date';

const loading = ref(false);
const isLoaded = ref(false);
const connections = ref<Database.S3Connection[]>([]);
const currentConnId = ref<number>();
const buckets = ref<Database.S3BucketInfo[]>([]);
const currentBucket = ref('');
const prefix = ref('');
const objects = ref<Database.S3ObjectItem[]>([]);
const selectedRows = ref<Database.S3ObjectItem[]>([]);
const rustfs = ref<Database.S3RustFSList>({ installed: false, items: [] });
const connRef = ref();
const uploadRef = ref();
const previewVisible = ref(false);
const preview = ref<Database.S3Preview>();

const currentConn = computed(() => connections.value.find((item) => item.id === currentConnId.value));

const prefixSegments = computed(() => {
    const parts = prefix.value.split('/').filter(Boolean);
    const items: { name: string; prefix: string }[] = [];
    let current = '';
    for (const part of parts) {
        current += part + '/';
        items.push({ name: part, prefix: current });
    }
    return items;
});

const buttons = [
    {
        label: i18n.global.t('database.s3.preview'),
        disabled: (row: Database.S3ObjectItem) => row.prefix,
        click: (row: Database.S3ObjectItem) => onPreview(row),
    },
    {
        label: i18n.global.t('commons.button.download'),
        disabled: (row: Database.S3ObjectItem) => row.prefix,
        click: (row: Database.S3ObjectItem) => onDownload(row),
    },
    {
        label: i18n.global.t('commons.button.delete'),
        click: (row: Database.S3ObjectItem) => onDelete([row]),
    },
];

const goInstallRustFS = () => {
    routerToNameWithQuery('AppAll', { install: 'rustfs' });
};

const loadRustFS = async () => {
    const res = await listS3RustFS();
    rustfs.value = res.data;
};

const loadConnections = async (selectId?: number) => {
    const res = await listS3Connections();
    connections.value = res.data || [];
    if (selectId && connections.value.some((item) => item.id === selectId)) {
        currentConnId.value = selectId;
    } else if (currentConnId.value && connections.value.some((item) => item.id === currentConnId.value)) {
        // keep
    } else {
        currentConnId.value = connections.value[0]?.id;
    }
};

const onLoadBuckets = async () => {
    if (!currentConnId.value) {
        buckets.value = [];
        currentBucket.value = '';
        objects.value = [];
        return;
    }
    const res = await listS3Buckets(currentConnId.value);
    buckets.value = res.data || [];
    if (currentBucket.value && buckets.value.some((item) => item.name === currentBucket.value)) {
        await loadObjects();
        return;
    }
    const preferred = currentConn.value?.bucket;
    currentBucket.value = buckets.value.find((item) => item.name === preferred)?.name || buckets.value[0]?.name || '';
    if (currentBucket.value) {
        prefix.value = '';
        await loadObjects();
    } else {
        objects.value = [];
    }
};

const loadObjects = async () => {
    if (!currentConnId.value || !currentBucket.value) {
        objects.value = [];
        return;
    }
    const res = await listS3Objects(currentConnId.value, currentBucket.value, prefix.value);
    objects.value = res.data.objects || [];
    selectedRows.value = [];
};

const onChangeConn = async () => {
    prefix.value = '';
    currentBucket.value = '';
    await onLoadBuckets();
};

const onChangeBucket = async () => {
    prefix.value = '';
    await loadObjects();
};

const enterPrefix = async (next: string) => {
    prefix.value = next;
    await loadObjects();
};

const onConnSaved = async (connId?: number) => {
    await loadConnections(connId);
    await onLoadBuckets();
};

const onAddConn = () => {
    connRef.value.acceptParams({
        title: 'add',
        rowData: { name: '', source: 'manual', pathStyle: true, useSSL: false },
        rustfsItems: rustfs.value.items,
    });
};

const onEditConn = () => {
    if (!currentConn.value) return;
    connRef.value.acceptParams({
        title: 'edit',
        rowData: { ...currentConn.value, secretKey: '' },
        rustfsItems: rustfs.value.items,
    });
};

const onConnectRustFS = async () => {
    await loadRustFS();
    if (!rustfs.value.installed) {
        goInstallRustFS();
        return;
    }
    const usable = rustfs.value.items.filter((item) => !item.error);
    if (usable.length === 1) {
        loading.value = true;
        try {
            const res = await connectS3RustFS(usable[0].appInstallId);
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            await onConnSaved(res.data.id);
        } finally {
            loading.value = false;
        }
        return;
    }
    connRef.value.acceptParams({
        title: 'add',
        rowData: {
            name: usable[0]?.name || '',
            source: 'rustfs',
            appInstallId: usable[0]?.appInstallId,
        },
        rustfsItems: rustfs.value.items,
    });
};

const onTestConn = async () => {
    if (!currentConnId.value) return;
    loading.value = true;
    try {
        await testS3Connection(currentConnId.value);
        MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        await onLoadBuckets();
    } finally {
        loading.value = false;
    }
};

const onDeleteConn = async () => {
    if (!currentConn.value) return;
    try {
        await ElMessageBox.confirm(i18n.global.t('commons.msg.delete'), i18n.global.t('commons.button.delete'), {
            type: 'warning',
        });
    } catch {
        return;
    }
    await deleteS3Connection(currentConn.value.id);
    MsgSuccess(i18n.global.t('commons.msg.deleteSuccess'));
    currentConnId.value = undefined;
    await loadConnections();
    await onLoadBuckets();
};

const onCreateBucket = async () => {
    if (!currentConnId.value) return;
    let value = '';
    try {
        const res = await ElMessageBox.prompt(
            i18n.global.t('database.s3.bucketName'),
            i18n.global.t('database.s3.createBucket'),
        );
        value = res.value;
    } catch {
        return;
    }
    const name = value?.trim();
    if (!name) return;
    await createS3Bucket(currentConnId.value, name, currentConn.value?.region);
    currentBucket.value = name;
    MsgSuccess(i18n.global.t('commons.msg.createSuccess'));
    await onLoadBuckets();
};

const onDeleteBucket = async () => {
    if (!currentConnId.value || !currentBucket.value) return;
    try {
        await ElMessageBox.confirm(
            i18n.global.t('database.s3.deleteBucketHelper', [currentBucket.value]),
            i18n.global.t('database.s3.deleteBucket'),
            { type: 'warning' },
        );
    } catch {
        return;
    }
    await deleteS3Bucket(currentConnId.value, currentBucket.value);
    currentBucket.value = '';
    MsgSuccess(i18n.global.t('commons.msg.deleteSuccess'));
    await onLoadBuckets();
};

const onCreateFolder = async () => {
    if (!currentConnId.value || !currentBucket.value) return;
    let value = '';
    try {
        const res = await ElMessageBox.prompt(
            i18n.global.t('database.s3.folderName'),
            i18n.global.t('database.s3.createFolder'),
        );
        value = res.value;
    } catch {
        return;
    }
    const name = value?.trim();
    if (!name) return;
    await createS3Folder(currentConnId.value, currentBucket.value, prefix.value, name);
    MsgSuccess(i18n.global.t('commons.msg.createSuccess'));
    await loadObjects();
};

const onUploadClick = () => {
    const input = uploadRef.value?.$el?.querySelector('input') as HTMLInputElement | undefined;
    input?.click();
};

const onUploadChange = async (file: UploadFile) => {
    if (!file.raw || !currentConnId.value || !currentBucket.value) return;
    const form = new FormData();
    form.append('id', String(currentConnId.value));
    form.append('bucket', currentBucket.value);
    form.append('prefix', prefix.value);
    form.append('file', file.raw);
    loading.value = true;
    try {
        await uploadS3Object(form);
        MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        await loadObjects();
    } finally {
        loading.value = false;
        uploadRef.value?.clearFiles?.();
    }
};

const onPreview = async (row: Database.S3ObjectItem) => {
    if (!currentConnId.value || !currentBucket.value) return;
    loading.value = true;
    try {
        const res = await previewS3Object(currentConnId.value, currentBucket.value, row.key);
        preview.value = res.data;
        previewVisible.value = true;
    } finally {
        loading.value = false;
    }
};

const onDownload = async (row: Database.S3ObjectItem) => {
    if (!currentConnId.value || !currentBucket.value) return;
    const blob = await downloadS3Object(currentConnId.value, currentBucket.value, row.key);
    const downloadUrl = window.URL.createObjectURL(new Blob([blob]));
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = downloadUrl;
    a.download = row.name || 'object';
    a.dispatchEvent(new MouseEvent('click'));
    setTimeout(() => window.URL.revokeObjectURL(downloadUrl), 1000);
};

const onDeleteSelected = async () => {
    await onDelete(selectedRows.value);
};

const onDelete = async (rows: Database.S3ObjectItem[]) => {
    if (!currentConnId.value || !currentBucket.value || rows.length === 0) return;
    const folder = rows.some((row) => row.prefix);
    const message = folder
        ? i18n.global.t('database.s3.deleteFolderHelper', [rows.map((row) => row.key).join(', ')])
        : i18n.global.t('database.s3.deleteObjectHelper', [rows.map((row) => row.key).join(', ')]);
    try {
        await ElMessageBox.confirm(message, i18n.global.t('commons.button.delete'), { type: 'warning' });
    } catch {
        return;
    }
    await deleteS3Objects(
        currentConnId.value,
        currentBucket.value,
        rows.map((row) => row.key),
    );
    MsgSuccess(i18n.global.t('commons.msg.deleteSuccess'));
    await loadObjects();
};

onMounted(async () => {
    loading.value = true;
    try {
        await loadRustFS();
        await loadConnections();
        await onLoadBuckets();
    } finally {
        loading.value = false;
        isLoaded.value = true;
    }
});
</script>

<style scoped>
.hidden-upload {
    display: none;
}
.preview-text {
    white-space: pre-wrap;
    word-break: break-word;
    font-size: 13px;
    line-height: 1.5;
}
</style>
