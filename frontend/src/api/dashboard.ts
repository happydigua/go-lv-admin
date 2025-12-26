import request from '@/utils/request';

export const getDashboardStats = () => {
    return request({
        url: '/dashboard/stats',
        method: 'get',
    });
};

export const getDashboardCharts = () => {
    return request({
        url: '/dashboard/charts',
        method: 'get',
    });
};
