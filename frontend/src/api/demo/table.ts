import request from '@/utils/request';

// 获取列表
export const getDemoList = (params: any) => {
    return request({
        url: '/demo/list',
        method: 'get',
        params
    });
};

// 创建
export const createDemo = (data: any) => {
    return request({
        url: '/demo',
        method: 'post',
        data
    });
};

// 更新
export const updateDemo = (data: any) => {
    return request({
        url: `/demo/${data.ID}`,
        method: 'put',
        data
    });
};

// 删除
export const deleteDemo = (id: number) => {
    return request({
        url: `/demo/${id}`,
        method: 'delete'
    });
};
