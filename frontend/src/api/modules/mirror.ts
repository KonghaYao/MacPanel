import http from '@/api';
import { Mirror } from '../interface/mirror';

export const listMirrors = () => {
    return http.get<Mirror.Ecosystem[]>('/mirrors');
};

export const applyMirror = (param: Mirror.ApplyReq) => {
    return http.post('/mirrors/apply', param);
};
