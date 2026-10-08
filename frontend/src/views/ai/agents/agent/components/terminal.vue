<template>
    <DrawerPro
        v-model="terminalVisible"
        :header="$t('menu.terminal')"
        @close="handleClose"
        :resource="title"
        fullScreen
        :size="isFullScreen ? 'full' : 'large'"
        :autoClose="false"
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
                @reconnect="reConnect"
            />
        </template>
        <template #content>
            <div class="terminal-pane">
                <el-form v-show="!terminalOpen" label-position="top" @submit.prevent>
                    <el-form-item :label="$t('commons.table.user')" prop="user">
                        <el-select v-model="form.user" :disabled="form.users.length <= 1" @change="reConnect">
                            <el-option v-for="item in form.users" :key="item" :label="item" :value="item" />
                        </el-select>
                    </el-form-item>
                </el-form>
                <div v-if="terminalOpen" class="terminal-pane-body">
                    <Terminal ref="terminalRef"></Terminal>
                </div>
            </div>
        </template>
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
import { nextTick, reactive, ref } from 'vue';
import Terminal from '@/components/terminal/index.vue';
import TerminalToolbar from '@/components/terminal/toolbar.vue';
import { useGlobalStore } from '@/composables/useGlobalStore';

const { currentNode, isFullScreen } = useGlobalStore();

const title = ref('');
const terminalVisible = ref(false);
const terminalOpen = ref(false);
const terminalRef = ref<InstanceType<typeof Terminal> | null>(null);
const form = reactive({
    containerID: '',
    users: [] as string[],
    user: '',
    shell: '',
    initCmd: '',
    node: '',
});

interface DialogProps {
    containerID: string;
    title: string;
    users: string[];
    shell: string;
    initCmd?: string;
    node?: string;
}

const acceptParams = async (params: DialogProps): Promise<void> => {
    terminalVisible.value = true;
    form.containerID = params.containerID;
    form.users = [...params.users];
    form.user = form.users[0];
    form.shell = params.shell;
    form.initCmd = params.initCmd || '';
    form.node = params.node || currentNode.value;
    title.value = params.title;
    await reConnect();
};

const initTerm = async () => {
    terminalOpen.value = true;
    await nextTick();
    let args = `source=container&containerid=${form.containerID}&user=${form.user}&command=${form.shell}`;
    if (form.node) {
        args += `&operateNode=${form.node}`;
    }
    terminalRef.value?.acceptParams({
        endpoint: '/api/v2/hosts/terminal/container',
        args,
        error: '',
        initCmd: form.initCmd,
    });
};

const reConnect = async () => {
    terminalRef.value?.onClose();
    terminalOpen.value = false;
    await nextTick();
    await initTerm();
};

const onClose = () => {
    terminalRef.value?.onClose();
    terminalOpen.value = false;
};

const handleClose = () => {
    onClose();
    terminalVisible.value = false;
};

defineExpose({
    acceptParams,
});
</script>

<style lang="scss">
@use '@/components/terminal/pane.scss';
</style>
