<template>
    <DialogPro
        v-model="drawerVisible"
        :header="$t('ssh.addHost')"
        @close="handleClose"
        size="small"
        :autoClose="false"
    >
        <el-alert class="mb-4" :closable="false" type="info">
            {{ $t('ssh.hostAddHelper') }}
        </el-alert>
        <el-form ref="formRef" label-position="top" :rules="rules" :model="form" v-loading="loading">
            <el-form-item :label="$t('home.alias')" prop="alias">
                <el-input v-model="form.alias" :placeholder="$t('ssh.aliasHelper')" />
            </el-form-item>
            <el-form-item :label="$t('ssh.hostAddress')" prop="hostName">
                <el-input v-model="form.hostName" :placeholder="$t('ssh.hostAddressHelper')" />
            </el-form-item>
            <el-form-item :label="$t('commons.table.user')" prop="user">
                <el-input v-model="form.user" />
            </el-form-item>
            <el-form-item :label="$t('commons.login.password')" prop="password">
                <el-input v-model="form.password" type="password" show-password />
                <span class="input-help">{{ $t('ssh.hostPasswordHelper') }}</span>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="handleClose">{{ $t('commons.button.cancel') }}</el-button>
            <el-button type="primary" :loading="loading" @click="onSubmit">
                {{ $t('commons.button.confirm') }}
            </el-button>
        </template>
    </DialogPro>
</template>

<script setup lang="ts">
import { createSSHHost } from '@/api/modules/host';
import { Host } from '@/api/interface/host';
import i18n from '@/lang';
import { MsgError, MsgSuccess } from '@/utils/message';
import { reactive, ref } from 'vue';

const emit = defineEmits(['search']);

const drawerVisible = ref(false);
const loading = ref(false);
const formRef = ref();

const form = reactive<Host.SSHHostOperate>({
    alias: '',
    hostName: '',
    user: '',
    port: 22,
    password: '',
});

const aliasPattern = /^[A-Za-z0-9][A-Za-z0-9_-]*$/;
const hostPattern =
    /^(\d{1,3}\.){3}\d{1,3}$|^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)*$/;

const rules = {
    alias: [
        { required: true, message: i18n.global.t('commons.rule.requiredInput'), trigger: 'blur' },
        {
            validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
                if (!aliasPattern.test(value || '')) {
                    callback(new Error(i18n.global.t('ssh.aliasHelper')));
                    return;
                }
                callback();
            },
            trigger: 'blur',
        },
    ],
    hostName: [
        { required: true, message: i18n.global.t('commons.rule.requiredInput'), trigger: 'blur' },
        {
            validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
                if (!hostPattern.test(value || '')) {
                    callback(new Error(i18n.global.t('ssh.hostAddressHelper')));
                    return;
                }
                callback();
            },
            trigger: 'blur',
        },
    ],
    user: [{ required: true, message: i18n.global.t('commons.rule.requiredInput'), trigger: 'blur' }],
    password: [{ required: true, message: i18n.global.t('commons.rule.requiredInput'), trigger: 'blur' }],
};

const resetForm = () => {
    form.alias = '';
    form.hostName = '';
    form.user = '';
    form.port = 22;
    form.password = '';
};

const acceptParams = () => {
    resetForm();
    drawerVisible.value = true;
};

const handleClose = () => {
    drawerVisible.value = false;
    formRef.value?.resetFields();
};

const onSubmit = async () => {
    await formRef.value?.validate(async (valid: boolean) => {
        if (!valid) {
            return;
        }
        loading.value = true;
        try {
            await createSSHHost({ ...form });
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            emit('search');
            handleClose();
        } catch (error) {
            MsgError(error);
        } finally {
            loading.value = false;
        }
    });
};

defineExpose({ acceptParams });
</script>
