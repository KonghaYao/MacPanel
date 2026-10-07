export namespace Mirror {
    export interface Field {
        key: string;
        kind: string;
    }

    export interface Preset {
        id: string;
        name: string;
        values: Record<string, string>;
        snippet: string;
    }

    export interface Ecosystem {
        id: string;
        configPath: string;
        activePreset: string;
        fields: Field[];
        current: Record<string, string>;
        snippet: string;
        presets: Preset[];
        readError?: string;
    }

    export interface ApplyReq {
        ecosystem: string;
        presetId?: string;
        values?: Record<string, string>;
    }
}
