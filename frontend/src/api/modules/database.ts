import http from '@/api';
import { deepCopy } from '@/utils/misc';
import { encodeBase64, encodeBase64Fields } from '@/utils/base64';
import { ResPage, DescriptionUpdate } from '../interface';
import { Database } from '../interface/database';
import { TimeoutEnum } from '@/enums/http-enum';

// common
export const loadDBBaseInfo = (type: string, database: string) => {
    return http.post<Database.BaseInfo>(`/databases/common/info`, { type: type, name: database });
};
export const loadDBFile = (type: string, database: string) => {
    return http.post<string>(`/databases/common/load/file`, { type: type, name: database });
};
export const updateDBFile = (params: Database.DBConfUpdate) => {
    return http.post(`/databases/common/update/conf`, params);
};

// pg
export const addPostgresqlDB = (params: Database.PostgresqlDBCreate) => {
    let request = deepCopy(params) as Database.PostgresqlDBCreate;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/pg`, request, TimeoutEnum.T_40S);
};
export const bindPostgresqlUser = (params: Database.PgBind) => {
    return http.post(`/databases/pg/bind`, params, TimeoutEnum.T_40S);
};
export const changePrivileges = (params: Database.PgChangePrivileges) => {
    return http.post(`/databases/pg/privileges`, params, TimeoutEnum.T_40S);
};
export const searchPostgresqlDBs = (params: Database.SearchDBWithPage, node?: string) => {
    const query = node ? `?operateNode=${node}` : '';
    return http.post<ResPage<Database.PostgresqlDBInfo>>(`/databases/pg/search${query}`, params);
};
export const updatePostgresqlDescription = (params: DescriptionUpdate) => {
    return http.post(`/databases/pg/description`, params);
};
export const loadPgFromRemote = (database: string) => {
    return http.post(`/databases/pg/${database}/load`);
};
export const deleteCheckPostgresqlDB = (params: Database.PostgresqlDBDeleteCheck) => {
    return http.post<Database.DBResource[]>(`/databases/pg/del/check`, params, TimeoutEnum.T_40S);
};
export const updatePostgresqlPassword = (params: Database.ChangeInfo) => {
    let request = deepCopy(params) as Database.ChangeInfo;
    encodeBase64Fields(request, ['value']);
    return http.post(`/databases/pg/password`, request, TimeoutEnum.T_40S);
};
export const deletePostgresqlDB = (params: Database.PostgresqlDBDelete) => {
    return http.post(`/databases/pg/del`, params, TimeoutEnum.T_40S);
};

// mysql
export const searchMysqlDBs = (params: Database.SearchDBWithPage, node?: string) => {
    const query = node ? `?operateNode=${node}` : '';
    return http.post<ResPage<Database.MysqlDBInfo>>(`/databases/search${query}`, params);
};
export const addMysqlDB = (params: Database.MysqlDBCreate) => {
    let request = deepCopy(params) as Database.MysqlDBCreate;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases`, request);
};
export const searchMysqlUsers = (params: Database.MysqlUserSearch) => {
    return http.post<Database.MysqlUser[]>(`/databases/users/search`, params);
};
export const createMysqlUser = (params: Database.MysqlUserCreate) => {
    let request = deepCopy(params) as Database.MysqlUserCreate;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/users`, request);
};
export const deleteMysqlUser = (params: Database.MysqlUserDelete) => {
    return http.post(`/databases/users/del`, params);
};
export const updateMysqlUser = (params: Database.MysqlUserUpdate) => {
    return http.post(`/databases/users/update`, params);
};
export const updateMysqlUserPassword = (params: Database.MysqlUserPassword) => {
    let request = deepCopy(params) as Database.MysqlUserPassword;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/users/password`, request);
};
export const saveMysqlUserPassword = (params: Database.MysqlUserPassword) => {
    let request = deepCopy(params) as Database.MysqlUserPassword;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/users/password/save`, request);
};
export const searchMysqlGrants = (params: Database.MysqlUserSearch) => {
    return http.post<Database.MysqlGrant[]>(`/databases/grants/search`, params);
};
export const searchMysqlGrantSummary = (params: Database.MysqlGrantSummarySearch) => {
    return http.post<Record<string, Database.MysqlUser[]>>(`/databases/grants/summary`, params);
};
export const grantMysqlUser = (params: Database.MysqlGrantCreate) => {
    return http.post(`/databases/grants`, params);
};
export const revokeMysqlGrant = (params: Database.MysqlGrantDelete) => {
    return http.post(`/databases/grants/del`, params);
};
export const loadDBFromRemote = (params: Database.MysqlLoadDB) => {
    return http.post(`/databases/load`, params);
};
export const updateMysqlAccess = (params: Database.ChangeInfo) => {
    return http.post(`/databases/change/access`, params);
};
export const updateMysqlPassword = (params: Database.ChangeInfo) => {
    let request = deepCopy(params) as Database.ChangeInfo;
    encodeBase64Fields(request, ['value']);
    return http.post(`/databases/change/password`, request);
};
export const updateMysqlDescription = (params: DescriptionUpdate) => {
    return http.post(`/databases/description/update`, params);
};
export const updateMysqlVariables = (params: Database.VariablesUpdate) => {
    return http.post(`/databases/variables/update`, params);
};
export const deleteCheckMysqlDB = (params: Database.MysqlDBDeleteCheck) => {
    return http.post<Array<string>>(`/databases/del/check`, params);
};
export const deleteMysqlDB = (params: Database.MysqlDBDelete) => {
    return http.post(`/databases/del`, params);
};

// mongodb
export const searchMongodbDBs = (params: Database.SearchDBWithPage, node?: string) => {
    const query = node ? `?operateNode=${node}` : '';
    return http.post<ResPage<Database.MongodbDBInfo>>(`/databases/mongodb/search${query}`, params);
};
export const addMongodbDB = (params: Database.MongodbDBCreate) => {
    let request = deepCopy(params) as Database.MongodbDBCreate;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/mongodb`, request, TimeoutEnum.T_40S);
};
export const loadMongodbFromRemote = (params: Database.MongodbLoadDB) => {
    return http.post(`/databases/mongodb/load`, params, TimeoutEnum.T_40S);
};
export const bindMongodbUser = (params: Database.MongodbBind) => {
    let request = deepCopy(params) as Database.MongodbBind;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/mongodb/bind`, request, TimeoutEnum.T_40S);
};
export const updateMongodbPassword = (params: Database.MongodbPassword) => {
    let request = deepCopy(params) as Database.MongodbPassword;
    encodeBase64Fields(request, ['password']);
    return http.post(`/databases/mongodb/password`, request, TimeoutEnum.T_40S);
};
export const updateMongodbRootPassword = (params: Database.ChangeInfo) => {
    let request = deepCopy(params) as Database.ChangeInfo;
    encodeBase64Fields(request, ['value']);
    return http.post(`/databases/mongodb/root/password`, request, TimeoutEnum.T_40S);
};
export const updateMongodbDescription = (params: DescriptionUpdate) => {
    return http.post(`/databases/mongodb/description`, params, TimeoutEnum.T_40S);
};
export const deleteCheckMongodbDB = (params: Database.MongodbDBDeleteCheck) => {
    return http.post<Database.DBResource[]>(`/databases/mongodb/del/check`, params, TimeoutEnum.T_40S);
};
export const deleteMongodbDB = (params: Database.MongodbDBDelete) => {
    return http.post(`/databases/mongodb/del`, params, TimeoutEnum.T_40S);
};
export const loadMongodbPrivileges = (params: Database.MongodbPrivilegesLoad) => {
    return http.post<string>(`/databases/mongodb/privileges`, params, TimeoutEnum.T_40S);
};
export const changeMongodbPrivileges = (params: Database.MongodbPrivileges) => {
    return http.post(`/databases/mongodb/privileges/change`, params, TimeoutEnum.T_40S);
};

export const loadMysqlVariables = (type: string, database: string) => {
    return http.post<Database.MysqlVariables>(`/databases/variables`, { type: type, name: database });
};
export const loadMysqlStatus = (type: string, database: string) => {
    return http.post<Database.MysqlStatus>(`/databases/status`, { type: type, name: database });
};
export const loadRemoteAccess = (type: string, database: string) => {
    return http.post<boolean>(`/databases/remote`, { type: type, name: database });
};
export const loadFormatCollations = (database: string) => {
    return http.post<Array<Database.FormatCollationOption>>(`/databases/format/options`, { name: database });
};

// redis
export const loadRedisStatus = (type: string, database: string) => {
    return http.post<Database.RedisStatus>(`/databases/redis/status`, { type: type, name: database });
};
export const loadRedisConf = (type: string, database: string) => {
    return http.post<Database.RedisConf>(`/databases/redis/conf`, { type: type, name: database });
};
export const redisPersistenceConf = (type: string, database: string) => {
    return http.post<Database.RedisPersistenceConf>(`/databases/redis/persistence/conf`, {
        type: type,
        name: database,
    });
};
export const checkRedisCli = () => {
    return http.get<boolean>(`/databases/redis/check`);
};
export const installRedisCli = () => {
    return http.post(`/databases/redis/install/cli`, {}, TimeoutEnum.T_5M);
};
export const changeRedisPassword = (database: string, password: string) => {
    if (password) {
        password = encodeBase64(password);
    }
    return http.post(`/databases/redis/password`, { database: database, value: password });
};
export const updateRedisPersistenceConf = (params: Database.RedisConfPersistenceUpdate) => {
    return http.post(`/databases/redis/persistence/update`, params);
};
export const updateRedisConf = (params: Database.RedisConfUpdate) => {
    return http.post(`/databases/redis/conf/update`, params);
};
export const updateRedisConfByFile = (params: Database.DBConfUpdate) => {
    return http.post(`/databases/redis/conffile/update`, params);
};

// database
export const getDatabase = (name: string) => {
    return http.get<Database.DatabaseInfo>(`/databases/db/${name}`);
};
export const searchDatabases = (params: Database.SearchDatabasePage) => {
    return http.post<ResPage<Database.DatabaseInfo>>(`/databases/db/search`, params);
};
export const listDatabases = (type: string, node?: string) => {
    const query = node ? `?operateNode=${node}` : '';
    return http.get<Array<Database.DatabaseOption>>(`/databases/db/list/${type}${query}`);
};
export const listDbItems = (type: string) => {
    return http.get<Array<Database.DbItem>>(`/databases/db/item/${type}`);
};
export const checkDatabase = (params: Database.DatabaseCreate) => {
    let request = deepCopy(params) as Database.DatabaseCreate;
    if (request.ssl) {
        encodeBase64Fields(request, ['clientKey', 'clientCert', 'rootCert']);
    }

    return http.post<boolean>(`/databases/db/check`, request, TimeoutEnum.T_60S);
};
export const addDatabase = (params: Database.DatabaseCreate) => {
    let request = deepCopy(params) as Database.DatabaseCreate;
    if (request.ssl) {
        encodeBase64Fields(request, ['clientKey', 'clientCert', 'rootCert']);
    }

    return http.post(`/databases/db`, request, TimeoutEnum.T_60S);
};
export const editDatabase = (params: Database.DatabaseUpdate) => {
    let request = deepCopy(params) as Database.DatabaseCreate;
    if (request.ssl) {
        encodeBase64Fields(request, ['clientKey', 'clientCert', 'rootCert']);
    }

    return http.post(`/databases/db/update`, request, TimeoutEnum.T_60S);
};
export const deleteCheckDatabase = (id: number) => {
    return http.post<Database.DBResource[]>(`/databases/db/del/check`, { id: id });
};
export const deleteDatabase = (params: Database.DatabaseDelete) => {
    return http.post(`/databases/db/del`, params);
};

export const listS3RustFS = () => {
    return http.get<Database.S3RustFSList>(`/databases/s3/rustfs`);
};
export const connectS3RustFS = (appInstallId: number) => {
    return http.post<Database.S3Connection>(`/databases/s3/rustfs/connect`, { appInstallId });
};
export const listS3Connections = () => {
    return http.get<Database.S3Connection[]>(`/databases/s3/conn`);
};
export const createS3Connection = (params: Database.S3ConnectionCreate) => {
    let request = deepCopy(params) as Database.S3ConnectionCreate;
    encodeBase64Fields(request, ['secretKey']);
    return http.post(`/databases/s3/conn`, request);
};
export const updateS3Connection = (params: Database.S3ConnectionCreate & { id: number }) => {
    let request = deepCopy(params) as Database.S3ConnectionCreate & { id: number };
    if (request.secretKey) {
        encodeBase64Fields(request, ['secretKey']);
    }
    return http.post(`/databases/s3/conn/update`, request);
};
export const deleteS3Connection = (id: number) => {
    return http.post(`/databases/s3/conn/del`, { id });
};
export const testS3Connection = (id: number) => {
    return http.post(`/databases/s3/conn/test`, { id }, TimeoutEnum.T_40S);
};
export const listS3Buckets = (id: number) => {
    return http.post<Database.S3BucketInfo[]>(`/databases/s3/buckets`, { id }, TimeoutEnum.T_40S);
};
export const createS3Bucket = (id: number, bucket: string, region?: string) => {
    return http.post(`/databases/s3/buckets/create`, { id, bucket, region });
};
export const deleteS3Bucket = (id: number, bucket: string) => {
    return http.post(`/databases/s3/buckets/del`, { id, bucket });
};
export const listS3Objects = (id: number, bucket: string, prefix: string) => {
    return http.post<Database.S3ObjectList>(`/databases/s3/objects`, { id, bucket, prefix }, TimeoutEnum.T_40S);
};
export const createS3Folder = (id: number, bucket: string, prefix: string, name: string) => {
    return http.post(`/databases/s3/folder`, { id, bucket, prefix, name });
};
export const deleteS3Objects = (id: number, bucket: string, keys: string[]) => {
    return http.post(`/databases/s3/objects/del`, { id, bucket, keys }, TimeoutEnum.T_5M);
};
export const previewS3Object = (id: number, bucket: string, key: string) => {
    return http.post<Database.S3Preview>(`/databases/s3/objects/preview`, { id, bucket, key }, TimeoutEnum.T_40S);
};
export const downloadS3Object = (id: number, bucket: string, key: string) => {
    return http.download<BlobPart>(
        `/databases/s3/objects/download`,
        { id, bucket, key },
        { responseType: 'blob', timeout: TimeoutEnum.T_10M },
    );
};
export const uploadS3Object = (params: FormData) => {
    return http.upload(`/databases/s3/objects/upload`, params, {
        timeout: TimeoutEnum.T_10M,
    });
};
