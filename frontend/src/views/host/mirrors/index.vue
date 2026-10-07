<template>
    <div v-loading="loading">
        <LayoutContent :divider="true">
            <template #title>{{ $t('menu.mirrors') }}</template>
            <template #prompt>
                <el-alert type="info" :closable="false">
                    {{ $t('mirrors.intro') }}
                </el-alert>
            </template>
            <template #main>
                <el-tabs v-if="ecosystems.length" v-model="activeTab">
                    <el-tab-pane v-for="eco in ecosystems" :key="eco.id" :label="ecosystemTitle(eco.id)" :name="eco.id">
                        <div class="mb-3 flex flex-wrap items-center gap-2">
                            <el-tag v-if="eco.readError" type="danger">{{ $t('mirrors.readFailed') }}</el-tag>
                            <el-tag
                                v-else-if="eco.activePreset"
                                :type="eco.activePreset === 'official' ? 'info' : 'success'"
                            >
                                {{ activeLabel(eco) }}
                            </el-tag>
                            <el-tag v-else type="warning">{{ $t('mirrors.custom') }}</el-tag>
                        </div>

                        <div class="mb-3 flex flex-wrap items-center gap-1 text-sm text-gray-500">
                            <span>{{ $t('mirrors.configPath') }}</span>
                            <el-link
                                v-if="eco.configPath"
                                type="primary"
                                :underline="false"
                                class="break-all"
                                @click="openConfigPreview(eco.configPath)"
                            >
                                {{ eco.configPath }}
                            </el-link>
                            <CopyButton v-if="eco.configPath" :content="eco.configPath" />
                        </div>

                        <el-alert v-if="noteOf(eco.id)" class="mb-3" type="info" :closable="false">
                            {{ noteOf(eco.id) }}
                        </el-alert>
                        <el-alert v-if="eco.readError" class="mb-3" type="warning" :closable="false">
                            {{ eco.readError }}
                        </el-alert>

                        <div class="mb-4 flex flex-wrap items-center gap-2">
                            <div v-for="preset in eco.presets" :key="preset.id" class="flex items-center">
                                <el-button
                                    :type="eco.activePreset === preset.id ? 'primary' : 'default'"
                                    :disabled="applyingKey !== ''"
                                    :loading="applyingKey === keyOf(eco.id, preset.id)"
                                    @click="onPreset(eco.id, preset.id)"
                                >
                                    {{ presetLabel(preset) }}
                                </el-button>
                                <CopyButton :content="preset.snippet" />
                            </div>
                        </div>

                        <el-form v-if="drafts[eco.id]" label-position="top" class="max-w-3xl">
                            <el-form-item v-for="field in eco.fields" :key="field.key" :label="fieldLabel(field.key)">
                                <el-input
                                    v-if="field.kind === 'urls'"
                                    v-model="drafts[eco.id][field.key]"
                                    type="textarea"
                                    :rows="4"
                                />
                                <el-input v-else v-model="drafts[eco.id][field.key]" />
                            </el-form-item>
                            <el-form-item>
                                <el-button
                                    type="primary"
                                    plain
                                    :disabled="applyingKey !== ''"
                                    :loading="applyingKey === keyOf(eco.id, 'custom')"
                                    @click="onCustom(eco.id)"
                                >
                                    {{ $t('mirrors.saveCustom') }}
                                </el-button>
                            </el-form-item>
                        </el-form>

                        <div class="mt-2 rounded bg-[var(--el-fill-color-light)] p-3">
                            <div class="flex items-center justify-between gap-2">
                                <span class="text-sm">{{ $t('mirrors.snippet') }}</span>
                                <CopyButton :content="eco.snippet" :is-icon="false" />
                            </div>
                            <pre class="mt-2 whitespace-pre-wrap break-all text-xs leading-5">{{ eco.snippet }}</pre>
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </template>
        </LayoutContent>
        <TextPreview ref="textPreviewRef" />
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { applyMirror, listMirrors } from '@/api/modules/mirror';
import { Mirror } from '@/api/interface/mirror';
import { MsgSuccess } from '@/utils/message';
import TextPreview from '@/views/host/file-management/text-preview/index.vue';

const { t, te } = useI18n();
const loading = ref(false);
const applyingKey = ref('');
const activeTab = ref('');
const ecosystems = ref<Mirror.Ecosystem[]>([]);
const drafts = reactive<Record<string, Record<string, string>>>({});
const textPreviewRef = ref<InstanceType<typeof TextPreview>>();

const configFileName = (configPath: string) => {
    const parts = configPath.split(/[/\\]/);
    return parts[parts.length - 1] || configPath;
};

const openConfigPreview = (configPath: string) => {
    textPreviewRef.value?.acceptParams({ path: configPath, name: configFileName(configPath) });
};

const keyOf = (ecosystem: string, action: string) => `${ecosystem}:${action}`;

const textOrEmpty = (key: string) => (te(key) ? t(key) : '');

const ecosystemTitle = (id: string) => textOrEmpty(`mirrors.ecosystems.${id}`) || id;

const noteOf = (id: string) => textOrEmpty(`mirrors.notes.${id}`);

const fieldLabel = (key: string) => textOrEmpty(`mirrors.fields.${key}`) || key;

const presetLabel = (preset: Mirror.Preset) => textOrEmpty(`mirrors.presets.${preset.id}`) || preset.name;

const activeLabel = (eco: Mirror.Ecosystem) => {
    const preset = eco.presets.find((item) => item.id === eco.activePreset);
    return preset ? presetLabel(preset) : eco.activePreset;
};

const syncDrafts = (items: Mirror.Ecosystem[]) => {
    for (const eco of items) {
        const next: Record<string, string> = {};
        for (const field of eco.fields) {
            next[field.key] = eco.current?.[field.key] || '';
        }
        drafts[eco.id] = next;
    }
};

const ensureActiveTab = () => {
    if (!ecosystems.value.length) {
        activeTab.value = '';
        return;
    }
    if (!ecosystems.value.some((eco) => eco.id === activeTab.value)) {
        activeTab.value = ecosystems.value[0].id;
    }
};

const load = async () => {
    const res = await listMirrors();
    ecosystems.value = res.data || [];
    syncDrafts(ecosystems.value);
    ensureActiveTab();
};

const onPreset = async (ecosystem: string, presetId: string) => {
    applyingKey.value = keyOf(ecosystem, presetId);
    try {
        await applyMirror({ ecosystem, presetId });
        MsgSuccess(t('commons.msg.operationSuccess'));
        await load();
    } catch {
        return;
    } finally {
        applyingKey.value = '';
    }
};

const onCustom = async (ecosystem: string) => {
    applyingKey.value = keyOf(ecosystem, 'custom');
    try {
        await applyMirror({ ecosystem, values: { ...drafts[ecosystem] } });
        MsgSuccess(t('commons.msg.operationSuccess'));
        await load();
    } catch {
        return;
    } finally {
        applyingKey.value = '';
    }
};

onMounted(async () => {
    loading.value = true;
    try {
        await load();
    } catch {
        ecosystems.value = [];
    } finally {
        loading.value = false;
    }
});
</script>
