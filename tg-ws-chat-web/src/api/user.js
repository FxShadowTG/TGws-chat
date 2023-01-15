import request from "@/api";

// 查询主机列表
export function getHostCount(query) {
    return request({
        url: '/gethostcount',
        method: 'get',
        params: query
    })
}
