<template>
    <DrawerPro
        v-model="open"
        :title="$t('menu.terminal')"
        @close="handleClose"
        :resource="title"
        :autoClose="false"
        :fullScreen="true"
    >
        <template #extra>
            <TerminalToolbar
                inline
                :can-copy="!!terminalRef?.ui?.canCopy"
                show-reconnect
                @search="terminalRef?.openSearch()"
                @copy="terminalRef?.copySelection()"
                @paste="terminalRef?.pasteClipboard()"
                @clear="terminalRef?.clearScreen()"
                @reconnect="reconnect"
            />
        </template>
        <div class="terminal-pane">
            <div class="terminal-pane-body">
                <Terminal ref="terminalRef"></Terminal>
            </div>
        </div>

        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="handleClose">
                    {{ $t('commons.button.disConn') }}
                </el-button>
            </span>
        </template>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { nextTick, ref } from 'vue';
import Terminal from '@/components/terminal/index.vue';
import TerminalToolbar from '@/components/terminal/toolbar.vue';
import { closeOllamaModel } from '@/api/modules/ai';

const title = ref();
const open = ref(false);
const itemName = ref();
const terminalRef = ref();

interface DialogProps {
    name: string;
}
const acceptParams = async (params: DialogProps): Promise<void> => {
    itemName.value = params.name;
    open.value = true;
    initTerm();
};

const initTerm = () => {
    nextTick(() => {
        terminalRef.value.acceptParams({
            endpoint: '/api/v2/hosts/terminal/container',
            args: `source=ollama&name=${itemName.value}`,
            error: '',
            initCmd: '',
        });
    });
};

const reconnect = async () => {
    terminalRef.value?.onClose();
    await nextTick();
    initTerm();
};

const onClose = async () => {
    await closeOllamaModel(itemName.value)
        .then(() => {
            terminalRef.value?.onClose();
        })
        .catch(() => {
            terminalRef.value?.onClose();
        });
};

function handleClose() {
    onClose();
    open.value = false;
}

defineExpose({
    acceptParams,
});
</script>

<style lang="scss">
@use '@/components/terminal/pane.scss';
</style>
