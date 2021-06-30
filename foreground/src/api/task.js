import service from '../utils/request'

// 任务列表
export const getTaskList = (data) => {
    return service({
        url: "/sys/getTaskList",
        method: 'post',
        data: data
    })
}