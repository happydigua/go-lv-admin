import request from '@/utils/request';

// 获取所有设置
export const getSettings = () => {
    return request({
        url: '/settings',
        method: 'get'
    });
};

// 获取公开设置（无需登录）
export const getPublicSettings = () => {
    return request({
        url: '/settings/public',
        method: 'get'
    });
};

// 更新设置
export const updateSettings = (data: Record<string, string>) => {
    return request({
        url: '/settings',
        method: 'put',
        data
    });
};
