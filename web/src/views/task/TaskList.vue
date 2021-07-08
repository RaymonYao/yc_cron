<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="任务状态" prop="condition">
        <el-select style="width: 120px" v-model.number="searchForm.status">
          <el-option v-for="v in tStatus" :key="v.value" :label="v.label" :value="v.value"/>
        </el-select>
      </el-form-item>
      <el-form-item label="筛选条件" prop="condition">
        <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
          <el-option key="task_name" label="任务名称" value="task_name"/>
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
      <el-button @click="showTaskEdit({})"><i class="el-icon-plus"></i> 新增任务</el-button>
    </div>
    <el-table :data="tableData" border style="width: 100%">
      <el-table-column
          prop="task_id"
          label="ID"
          width="120"
          align="center">
      </el-table-column>
      <el-table-column
          prop="task_name"
          label="任务名称"
          width="180"
          align="center">
      </el-table-column>
      <el-table-column
          prop="cron_spec"
          label="时间表达式"
          align="center">
      </el-table-column>
      <el-table-column
          prop="prev_execute_time"
          label="上次执行时间"
          align="center"
          :formatter="dateFormat">
      </el-table-column>
      <el-table-column
          prop="next_execute_time"
          label="下次执行时间"
          align="center"
          :formatter="dateFormat">
      </el-table-column>
      <el-table-column align="center" label="操作">
        <template slot-scope="scope">
          <el-button title="编辑" type="primary" circle @click="showTaskEdit(scope.row)" icon="el-icon-edit"></el-button>
          <el-button title="暂停" type="success" icon="el-icon-switch-button" circle @click="pauseTask(scope.row)"
                     slot="reference" v-show="scope.row.status"></el-button>
          <el-button title="开启" type="warning" icon="el-icon-caret-right" circle @click="startTask(scope.row)"
                     slot="reference" v-show="!scope.row.status"></el-button>
          <el-button title="立即执行" type="primary" icon="el-icon-refresh-right" circle @click="runTask(scope.row)"
                     slot="reference"></el-button>
          <el-button title="强杀" type="danger" icon="el-icon-scissors" circle @click="killTask(scope.row)"
                     slot="reference"></el-button>
          <el-button title="删除" type="danger" icon="el-icon-delete" circle @click="delTask(scope.row)"
                     slot="reference" v-show="!scope.row.status"></el-button>
          <el-button title="查看日志" type="info" icon="el-icon-document" circle @click="delTask(scope.row)"
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
    <TaskEdit ref="taskEditorFormDrawer"></TaskEdit>
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {delTask, getTaskList, killTask, pauseTask, runTask, startTask} from "@/api/task";
import TaskEdit from "./cpns/TaskEdit";
import dateTool from "@/plugins/mixins/dateTool";

export default {
  name: "TaskList",
  mixins: [tableInfo, dateTool],
  components: {TaskEdit},
  data() {
    return {
      searchForm: {
        status: -1   //默认请选择
      },
      tStatus: [
        {
          value: -1,
          label: '请选择'
        },
        {
          value: 0,
          label: '未开启'
        },
        {
          value: 1,
          label: '已开启'
        }
      ],
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    getList: getTaskList,
    showTaskEdit(row) {
      this.$refs.taskEditorFormDrawer.setEditVal(row)
    },
    delTask(row) {
      this.$confirm('确认删除该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await delTask({task_id: row.task_id}).then((res) => {
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
    },
    pauseTask(row) {
      this.$confirm('确认暂停该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await pauseTask({task_id: row.task_id}).then((res) => {
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
    },
    startTask(row) {
      this.$confirm('确认开启该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await startTask({task_id: row.task_id}).then((res) => {
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
    },
    runTask(row) {
      this.$confirm('确认执行该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await runTask({task_id: row.task_id}).then((res) => {
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
    },
    killTask(row) {
      this.$confirm('确认强杀该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await killTask({task_id: row.task_id}).then((res) => {
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
