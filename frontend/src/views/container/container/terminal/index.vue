<template>
    <DrawerPro
        v-model="terminalVisible"
        :header="$t('menu.terminal')"
        @close="handleClose"
        :resource="title"
        :autoClose="!terminalOpen"
        size="large"
        :fullScreen="true"
    >
        <template #extra>
            <TerminalToolbar
                v-if="terminalOpen"
                inline
                :can-copy="!!terminalRef?.ui?.canCopy"
                show-reconnect
                @search="terminalRef?.openSearch()"
                @copy="terminalRef?.copySelection()"
                @paste="terminalRef?.pasteClipboard()"
                @clear="terminalRef?.clearScreen()"
                @reconnect="reconnect"
            >
                <el-button link @click="onClose()">{{ $t('commons.button.disConn') }}</el-button>
            </TerminalToolbar>
        </template>
        <template #content>
            <div class="terminal-pane">
            <el-form v-show="!terminalOpen" ref="formRef" :model="form" label-position="top">
                <el-form-item :label="$t('menu.container')" prop="containerID" v-if="form.containerIDList.length > 1">
                    <el-select placeholder="container" clearable v-model="form.containerID">
                        <el-option v-for="item in form.containerIDList" :key="item" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('commons.table.user')" prop="user">
                    <el-input placeholder="root" clearable v-model="form.user" />
                </el-form-item>
                <el-form-item
                    v-if="form.isCustom"
                    :label="$t('container.command')"
                    prop="command"
                    :rules="Rules.requiredInput"
                >
                    <el-checkbox class="p-w-100" border v-model="form.isCustom" @change="onChangeCommand">
                        {{ $t('container.custom') }}
                    </el-checkbox>
                    <el-input style="width: calc(100% - 100px)" clearable v-model="form.command" />
                </el-form-item>
                <el-form-item
                    v-if="!form.isCustom"
                    :label="$t('container.command')"
                    prop="command"
                    :rules="Rules.requiredSelect"
                >
                    <el-checkbox class="p-w-100" border v-model="form.isCustom" @change="onChangeCommand">
                        {{ $t('container.custom') }}
                    </el-checkbox>
                    <el-select style="width: calc(100% - 100px)" filterable clearable v-model="form.command">
                        <el-option value="/bin/ash" label="/bin/ash" />
                        <el-option value="/bin/bash" label="/bin/bash" />
                        <el-option value="/bin/sh" label="/bin/sh" />
                    </el-select>
                </el-form-item>

                <el-button @click="initTerm(formRef)">
                    {{ $t('commons.button.conn') }}
                </el-button>
            </el-form>
            <div v-if="terminalOpen" class="terminal-pane-body">
                <Terminal ref="terminalRef"></Terminal>
            </div>
            </div>
        </template>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { reactive, ref, nextTick } from 'vue';
import { ElForm, FormInstance } from 'element-plus';
import { Rules } from '@/global/form-rules';
import Terminal from '@/components/terminal/index.vue';
import TerminalToolbar from '@/components/terminal/toolbar.vue';
import { useGlobalStore } from '@/composables/useGlobalStore';
const { currentNode } = useGlobalStore();

const title = ref();
const terminalVisible = ref(false);
const terminalOpen = ref(false);
const form = reactive({
    isCustom: false,
    command: '',
    user: '',
    containerID: '',
    containerIDList: [],
    node: '',
});
const formRef = ref();
const terminalRef = ref<InstanceType<typeof Terminal> | null>(null);

interface DialogProps {
    containerID: string;
    title: string;
    node?: string;
}
const acceptParams = async (params: DialogProps): Promise<void> => {
    terminalVisible.value = true;
    form.containerIDList = params.containerID.split(',');
    form.containerID = form.containerIDList[0];
    title.value = params.title;
    form.isCustom = false;
    form.user = '';
    form.command = '/bin/sh';
    form.node = params.node || currentNode.value;
    terminalOpen.value = false;
};

const onChangeCommand = async () => {
    form.command = '';
};

const initTerm = (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid) => {
        if (!valid) return;
        terminalOpen.value = true;
        await nextTick();
        let args = `source=container&containerid=${form.containerID}&user=${form.user}&command=${form.command}`;
        if (form.node) {
            args += `&operateNode=${form.node}`;
        }
        terminalRef.value!.acceptParams({
            endpoint: '/api/v2/hosts/terminal/container',
            args: args,
            error: '',
            initCmd: '',
        });
    });
};

const onClose = () => {
    terminalRef.value?.onClose();
    terminalOpen.value = false;
};

const reconnect = async () => {
    onClose();
    await nextTick();
    initTerm(formRef.value);
};

function handleClose() {
    onClose();
    terminalVisible.value = false;
}

defineExpose({
    acceptParams,
});
</script>

<style lang="scss">
@use '@/components/terminal/pane.scss';
</style>
