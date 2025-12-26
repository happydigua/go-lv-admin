import request from '@/utils/request';

export const login = (data: any) => {
    return request({
        url: '/base/login',
        method: 'post',
        data,
    });
};
