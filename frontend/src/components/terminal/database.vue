<template>
    <DrawerPro
        v-model="open"
        :header="$t('menu.terminal')"
        @close="handleClose"
        :resource="resourceName"
        :autoClose="!open"
        size="large"
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
        <template #content>
            <div class="terminal-pane">
                <div class="terminal-pane-body">
                    <Terminal ref="terminalRef"></Terminal>
                </div>
            </div>
        </template>
    </DrawerPro>
</template>

<script lang="ts" setup>
import { ref, nextTick } from 'vue';
import Terminal from '@/components/terminal/index.vue';
import TerminalToolbar from '@/components/terminal/toolbar.vue';

const open = ref(false);
const terminalRef = ref<InstanceType<typeof Terminal> | null>(null);
const database = ref();
const databaseType = ref();
const resourceName = ref();
const command = ref('/bin/sh');
const user = ref('');
const containerID = ref('');
const initCmd = ref('');
const waitForPrompt = ref('');

interface DialogProps {
    databaseType?: string;
    database?: string;
    resourceName?: string;
    command?: string;
    user?: string;
    containerID?: string;
    initCmd?: string;
    waitForPrompt?: string;
}
const acceptParams = async (params: DialogProps): Promise<void> => {
    database.value = params.database || '';
    databaseType.value = params.databaseType || '';
    resourceName.value = params.resourceName || params.database || params.containerID || '';
    command.value = params.command || '/bin/sh';
    user.value = params.user || '';
    containerID.value = params.containerID || '';
    initCmd.value = params.initCmd || '';
    waitForPrompt.value = params.waitForPrompt || '';
    open.value = false;
    await initTerm();
};

const initTerm = async () => {
    open.value = true;
    await nextTick();
    const args = containerID.value
        ? `source=container&containerid=${containerID.value}&user=${user.value}&command=${command.value}`
        : `source=database&databaseType=${databaseType.value}&database=${database.value}`;
    terminalRef.value!.acceptParams({
        endpoint: '/api/v2/hosts/terminal/container',
        args,
        error: '',
        initCmd: initCmd.value,
        waitForPrompt: waitForPrompt.value,
    });
};

const reconnect = async () => {
    terminalRef.value?.onClose();
    await nextTick();
    await initTerm();
};

function handleClose() {
    terminalRef.value?.onClose();
    open.value = false;
}

defineExpose({
    acceptParams,
});
</script>

<style lang="scss">
@use '@/components/terminal/pane.scss';
</style>
