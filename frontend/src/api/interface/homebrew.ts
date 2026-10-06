export namespace Homebrew {
    export interface Status {
        isExist: boolean;
        version: string;
        prefix: string;
        formulaCount: number;
        caskCount: number;
        brewPath: string;
        cellarPath: string;
    }

    export interface Package {
        name: string;
        version: string;
        type: string;
    }

    export interface SearchResult {
        name: string;
        type: string;
    }

    export interface MirrorConfig {
        bottleDomain: string;
        apiDomain: string;
        brewGitRemote: string;
        coreGitRemote: string;
        caskGitRemote: string;
    }

    export interface DoctorResult {
        output: string;
    }

    export interface ListReq {
        page: number;
        pageSize: number;
        info?: string;
        type?: string;
    }

    export interface SearchReq {
        query: string;
        type?: string;
    }

    export interface PackageReq {
        name: string;
        type?: string;
        taskID?: string;
    }

    export interface TaskReq {
        taskID?: string;
    }
}
