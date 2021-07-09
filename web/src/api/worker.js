import service from '../utils/request'

// 日志列表
export const getWorkerList = (data) => {
    return service({
        url: "/worker/getWorkerList",
        method: 'post',
        data: data
    })
}