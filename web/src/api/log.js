import service from '../utils/request'

// 日志列表
export const getLogList = (data) => {
    return service({
        url: "/log/getLogList",
        method: 'post',
        data: data
    })
}