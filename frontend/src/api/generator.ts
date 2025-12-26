import request from '@/utils/request';

// 获取数据库表列表
export const getTables = () => {
    return request({ url: '/generator/tables', method: 'get' });
};

// 获取表的列信息
export const getTableColumns = (tableName: string) => {
    return request({ url: '/generator/columns', method: 'get', params: { tableName } });
};

// 预览生成的代码
export const previewCode = (data: any) => {
    return request({ url: '/generator/preview', method: 'post', data });
};

// 生成代码
export const generateCode = (data: any) => {
    return request({ url: '/generator/generate', method: 'post', data });
};
