/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const workerRouter = {
    path: '/worker',
    name: 'worker',
    component: Layout,
    redirect: '/worker/workerList',
    meta: {
        title: '节点管理',
        icon: 'el-icon-share'
    },
    children: [
        {
            path: "workerList",
            name: 'workerList',
            meta: {
                title: '节点列表'
            },
            component: () => import('../../views/worker/WorkerList')
        }
    ]
}
export default workerRouter
