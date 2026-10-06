<template>
    <DrawerPro
        v-model="drawerVisible"
        :header="title"
        :resource="dialogData.rowData?.name || ''"
        @close="handleClose"
        size="large"
    >
        <el-form ref="formRef" label-position="top" :model="dialogData.rowData" :rules="rules" v-loading="loading">
            <el-form-item :label="$t('commons.table.name')" prop="name">
                <el-input v-model.trim="dialogData.rowData!.name" />
            </el-form-item>
            <el-form-item :label="$t('database.s3.source')" prop="source">
                <el-radio-group v-model="dialogData.rowData!.source" :disabled="dialogData.title === 'edit'">
                    <el-radio value="manual">{{ $t('database.s3.sourceManual') }}</el-radio>
                    <el-radio value="rustfs">{{ $t('database.s3.sourceRustFS') }}</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item
                v-if="dialogData.rowData!.source === 'rustfs'"
                :label="$t('database.s3.rustfs')"
                prop="appInstallId"
            >
                <el-select
                    v-model="dialogData.rowData!.appInstallId"
                    class="w-full"
                    :disabled="dialogData.title === 'edit'"
                >
                    <el-option
                        v-for="item in rustfsItems"
                        :key="item.appInstallId"
                        :label="item.name"
                        :value="item.appInstallId"
                        :disabled="!!item.error"
                    >
                        <span>{{ item.name }}</span>
                        <span class="text-gray-400 ml-2">{{ item.status }}</span>
                    </el-option>
                </el-select>
                <span v-if="selectedRustFS?.error" class="input-help">{{ selectedRustFS.error }}</span>
            </el-form-item>
            <template v-if="dialogData.rowData!.source === 'manual'">
                <el-form-item :label="$t('database.s3.endpoint')" prop="endpoint">
                    <el-input v-model.trim="dialogData.rowData!.endpoint" />
                    <span class="input-help">{{ $t('database.s3.endpointHelper') }}</span>
                </el-form-item>
                <el-form-item :label="$t('database.s3.accessKey')" prop="accessKey">
                    <el-input v-model.trim="dialogData.rowData!.accessKey" />
                </el-form-item>
                <el-form-item :label="$t('database.s3.secretKey')" prop="secretKey">
                    <el-input v-model="dialogData.rowData!.secretKey" type="password" show-password />
                </el-form-item>
                <el-form-item :label="$t('database.s3.region')" prop="region">
                    <el-input v-model.trim="dialogData.rowData!.region" />
                </el-form-item>
                <el-form-item :label="$t('database.s3.useSSL')">
                    <el-switch v-model="dialogData.rowData!.useSSL" />
                </el-form-item>
                <el-form-item :label="$t('database.s3.pathStyle')">
                    <el-switch v-model="dialogData.rowData!.pathStyle" />
                    <span class="input-help">{{ $t('database.s3.pathStyleHelper') }}</span>
                </el-form-item>
            </template>
            <el-form-item :label="$t('commons.table.description')" prop="description">
                <el-input type="textarea" :rows="3" v-model="dialogData.rowData!.description" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="drawerVisible = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button v-permission :disabled="loading" type="primary" @click="onSubmit(formRef)">
                    {{ $t('commons.button.confirm') }}
                </el-button>
            </span>
        </template>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { computed, reactive, ref } from 'vue';
import { Rules } from '@/global/form-rules';
import i18n from '@/lang';
import { ElForm } from 'element-plus';
import { MsgSuccess } from '@/utils/message';
import { Database } from '@/api/interface/database';
import { connectS3RustFS, createS3Connection, updateS3Connection } from '@/api/modules/database';

interface DialogProps {
    title: string;
    rowData?: Database.S3ConnectionCreate & { id?: number };
    rustfsItems?: Database.S3RustFSInstall[];
}

const loading = ref(false);
const title = ref('');
const drawerVisible = ref(false);
const rustfsItems = ref<Database.S3RustFSInstall[]>([]);
const dialogData = ref<DialogProps>({ title: '' });

const selectedRustFS = computed(() =>
    rustfsItems.value.find((item) => item.appInstallId === dialogData.value.rowData?.appInstallId),
);

const acceptParams = (params: DialogProps): void => {
    dialogData.value = {
        title: params.title,
        rowData: {
            name: '',
            source: 'manual',
            pathStyle: true,
            useSSL: false,
            ...params.rowData,
        },
    };
    rustfsItems.value = params.rustfsItems || [];
    title.value = i18n.global.t('commons.button.' + dialogData.value.title);
    drawerVisible.value = true;
};

const emit = defineEmits<{ (e: 'search', connId?: number): void }>();

const handleClose = () => {
    drawerVisible.value = false;
};

const rules = reactive({
    name: [Rules.requiredInput],
    source: [Rules.requiredSelect],
    endpoint: [Rules.requiredInput],
    accessKey: [Rules.requiredInput],
    secretKey: [Rules.requiredInput],
});

type FormInstance = InstanceType<typeof ElForm>;
const formRef = ref<FormInstance>();

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    const row = dialogData.value.rowData!;
    if (row.source === 'rustfs') {
        if (!row.appInstallId) {
            return;
        }
        loading.value = true;
        try {
            const res = await connectS3RustFS(row.appInstallId);
            await updateS3Connection({
                id: res.data.id,
                name: row.name || res.data.name,
                source: 'rustfs',
                description: row.description || '',
            });
            drawerVisible.value = false;
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            emit('search', res.data.id);
        } finally {
            loading.value = false;
        }
        return;
    }

    const fields = ['name', 'endpoint', 'accessKey'];
    if (dialogData.value.title === 'add') {
        fields.push('secretKey');
    }
    formEl.validateField(fields, async (valid) => {
        if (!valid) return;
        loading.value = true;
        try {
            if (dialogData.value.title === 'edit' && row.id) {
                await updateS3Connection({
                    id: row.id,
                    name: row.name,
                    source: 'manual',
                    endpoint: row.endpoint,
                    accessKey: row.accessKey,
                    secretKey: row.secretKey,
                    region: row.region,
                    bucket: row.bucket,
                    useSSL: row.useSSL,
                    pathStyle: row.pathStyle,
                    description: row.description,
                });
            } else {
                await createS3Connection(row);
            }
            drawerVisible.value = false;
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            emit('search');
        } finally {
            loading.value = false;
        }
    });
};

defineExpose({
    acceptParams,
});
</script>
