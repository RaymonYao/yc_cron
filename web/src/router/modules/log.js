/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const logRouter = {
    path: '/log',
    name: 'log',
    component: Layout,
    redirect: '/log/logList',
    meta: {
        title: '日志管理',
        icon: 'el-icon-document'
    },
    children: [
        {
            path: "logList",
            name: 'logList',
            meta: {
                title: '日志列表'
            },
            component: () => import('../../views/log/LogList')
        }
    ]
}
export default logRouter
