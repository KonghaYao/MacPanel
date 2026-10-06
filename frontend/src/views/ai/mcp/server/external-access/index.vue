<template>
    <DialogPro v-model="open" :title="$t('aiTools.mcp.externalAccessTitle')" size="w-60" @close="handleClose">
        <el-alert type="info" :closable="false" show-icon class="mb-4">
            {{ $t('aiTools.mcp.externalAccessIntro') }}
        </el-alert>
        <ol class="m-0 list-decimal space-y-3 pl-5 leading-6">
            <li>{{ $t('aiTools.mcp.externalAccessStep1') }}</li>
            <li>{{ $t('aiTools.mcp.externalAccessStep2') }}</li>
            <li>{{ $t('aiTools.mcp.externalAccessStep3') }}</li>
            <li>
                <div>{{ $t('aiTools.mcp.externalAccessStep4') }}</div>
                <ul class="mt-2 list-disc space-y-1 pl-5">
                    <li>
                        <code>{{ $t('aiTools.mcp.externalAccessStep4Cursor') }}</code>
                    </li>
                    <li>
                        <code>{{ $t('aiTools.mcp.externalAccessStep4Claude') }}</code>
                    </li>
                </ul>
            </li>
            <li>{{ $t('aiTools.mcp.externalAccessStep5') }}</li>
            <li>{{ $t('aiTools.mcp.externalAccessStep6', [mcpEndpoint]) }}</li>
        </ol>
        <div class="mt-4">
            <div class="mb-2 flex items-center gap-2">
                <span class="font-medium">{{ $t('aiTools.mcp.externalAccessExample') }}</span>
                <CopyButton :content="exampleConfig" />
            </div>
            <pre class="overflow-auto rounded bg-[var(--el-fill-color-light)] p-3 text-xs leading-5">{{
                exampleConfig
            }}</pre>
        </div>
        <template #footer>
            <el-button @click="open = false">{{ $t('commons.button.cancel') }}</el-button>
            <el-button type="primary" icon="Position" @click="toApiKeys">
                {{ $t('aiTools.mcp.externalAccessGoApiKeys') }}
            </el-button>
        </template>
    </DialogPro>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import DialogPro from '@/components/dialog-pro/index.vue';
import CopyButton from '@/components/copy-button/index.vue';
import { buildMcpExampleConfig, getMcpEndpointUrl } from '@/utils/mcp-config';

const open = ref(false);
const router = useRouter();
const mcpEndpoint = computed(() => getMcpEndpointUrl());
const exampleConfig = computed(() => buildMcpExampleConfig());

const acceptParams = (): void => {
    open.value = true;
};

const handleClose = (): void => {
    open.value = false;
};

const toApiKeys = (): void => {
    router.push({ path: '/settings/apikeys', query: { uncached: 'true' } });
    open.value = false;
};

defineExpose({
    acceptParams,
});
</script>
