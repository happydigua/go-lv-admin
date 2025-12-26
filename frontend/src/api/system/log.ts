import request from '@/utils/request';

// 获取操作日志列表
export const getOperationLogList = (params: {
    page: number;
    pageSize: number;
    username?: string;
    module?: string;
    action?: string;
}) => {
    return request({
        url: '/system/log/list',
        method: 'get',
        params,
    });
};

// 批量删除操作日志
export const deleteOperationLogs = (ids: number[]) => {
    return request({
        url: '/system/log',
        method: 'delete',
        data: { ids },
    });
};

// 清空操作日志
export const clearOperationLogs = () => {
    return request({
        url: '/system/log/clear',
        method: 'delete',
    });
};
