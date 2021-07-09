import service from '../utils/request'

// 节点列表
export const getWorkerList = (data) => {
    return service({
        url: "/worker/getWorkerList",
        method: 'post',
        data: data
    })
}

//保存节点
export const saveWorker = (data) => {
    return service({
        url: "/worker/saveWorker",
        method: 'post',
        data: data
    })
}

//删除节点
export const delWorker = (data) => {
    return service({
        url: "/worker/delWorker",
        method: 'post',
        params: data
    })
}