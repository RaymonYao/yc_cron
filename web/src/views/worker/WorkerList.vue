<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="筛选条件" prop="condition">
        <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
          <el-option key="worker_ip" label="节点IP" value="worker_ip"/>
        </el-select>
      </el-form-item>
      <el-form-item prop="searchValue">
        <el-input v-model="searchForm.sValue" placeholder="搜索关键词"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="getTableData(1)" type="primary">查询</el-button>
      </el-form-item>
    </el-form>
    <div class="button-box">
      <el-button @click="showWorkerEdit({})"><i class="el-icon-plus"></i> 新增节点</el-button>
    </div>
    <el-table :data="tableData" border style="width: 100%">
      <el-table-column
          prop="worker_id"
          label="ID"
          width="120"
          align="center">
      </el-table-column>
      <el-table-column
          prop="worker_ip"
          label="节点IP"
          align="center"
          width="150">
      </el-table-column>
      <el-table-column
          prop="status"
          label="节点状态"
          align="center"
          width="100">
        <template slot-scope="scope">
          <el-button v-if="scope.row.status === 0" type="success" icon="el-icon-check" circle></el-button>
          <el-button v-if="scope.row.status === -1" type="danger" icon="el-icon-close" circle></el-button>
        </template>
      </el-table-column>
      <el-table-column
          prop="description"
          label="节点描述"
          width="500"
          align="center">
      </el-table-column>
      <el-table-column
          prop="create_time"
          label="创建时间"
          align="center"
          :formatter="dateFormat">
      </el-table-column>
      <el-table-column align="center" label="操作">
        <template slot-scope="scope">
          <el-button title="编辑" v-show="scope.row.status === -1" circle type="primary"
                     @click="showWorkerEdit(scope.row)" icon="el-icon-edit"></el-button>
          <el-button title="删除" v-show="scope.row.status === -1" type="danger" icon="el-icon-delete" circle
                     @click="delWorker(scope.row)"
                     slot="reference"></el-button>
          <el-button title="重启节点" v-show="scope.row.status === -1" type="danger" icon="el-icon-caret-right" circle
                     @click="delWorker(scope.row)"
                     slot="reference"></el-button>
          <el-button title="剔除节点" v-show="scope.row.status === 0" type="success" icon="el-icon-close-notification"
                     circle @click="delWorker(scope.row)"
                     slot="reference"></el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="page-content">
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
      </el-pagination>
    </div>
    <WorkerEdit ref="workerEditorFormDrawer"></WorkerEdit>
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {delWorker, getWorkerList} from "@/api/worker";
import dateTool from "@/plugins/mixins/dateTool";
import WorkerEdit from "./cpns/WorkerEdit";

export default {
  name: "WorkerList",
  mixins: [tableInfo, dateTool],
  components: {WorkerEdit},
  data() {
    return {
      searchForm: {
        status: -2   //默认请选择
      },
      wStatus: [
        {
          value: -2,
          label: '请选择'
        },
        {
          value: 0,
          label: '正常'
        },
        {
          value: -1,
          label: '异常'
        }
      ]
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    getList: getWorkerList,
    showWorkerEdit(row) {
      this.$refs.workerEditorFormDrawer.setEditVal(row)
    },
    delWorker(row) {
      this.$confirm('确认删除该节点?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await delWorker({worker_id: row.worker_id}).then((res) => {
          this.$message({
            type: 'success',
            message: res.msg
          })
          this.getTableData()
        }).catch(() => {
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消'
        })
      })
    }
  }
}
</script>
