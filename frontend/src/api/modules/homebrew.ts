import http from '@/api';
import { ResPage } from '../interface';
import { Homebrew } from '../interface/homebrew';
import { TimeoutEnum } from '@/enums/http-enum';

export const getHomebrewStatus = () => {
    return http.get<Homebrew.Status>('/homebrew/status', {}, TimeoutEnum.T_60S);
};

export const searchHomebrew = (param: Homebrew.SearchReq) => {
    return http.post<Homebrew.SearchResult[]>('/homebrew/search', param, TimeoutEnum.T_60S);
};

export const listHomebrew = (param: Homebrew.ListReq) => {
    return http.get<ResPage<Homebrew.Package>>('/homebrew/list', param, TimeoutEnum.T_60S);
};

export const installHomebrew = (param: Homebrew.PackageReq) => {
    return http.post('/homebrew/install', param, TimeoutEnum.T_10M);
};

export const uninstallHomebrew = (param: Homebrew.PackageReq) => {
    return http.post('/homebrew/uninstall', param, TimeoutEnum.T_10M);
};

export const upgradeHomebrew = (param: Homebrew.PackageReq) => {
    return http.post('/homebrew/upgrade', param, TimeoutEnum.T_10M);
};

export const updateHomebrew = (param: Homebrew.TaskReq) => {
    return http.post('/homebrew/update', param, TimeoutEnum.T_10M);
};

export const getHomebrewMirror = () => {
    return http.get<Homebrew.MirrorConfig>('/homebrew/mirror');
};

export const updateHomebrewMirror = (param: Homebrew.MirrorConfig) => {
    return http.post('/homebrew/mirror/update', param, TimeoutEnum.T_60S);
};

export const doctorHomebrew = () => {
    return http.post<Homebrew.DoctorResult>('/homebrew/doctor', {}, TimeoutEnum.T_10M);
};
