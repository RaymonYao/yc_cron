<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="任务状态" prop="condition">
        <el-select style="width: 120px" v-model.number="searchForm.status">
          <el-option v-for="v in uStatus" :key="v.value" :label="v.label" :value="v.value"/>
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
      <el-button @click="showUserEdit({})"><i class="el-icon-plus"></i> 新增任务</el-button>
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
          prop="last_execute_at"
          label="上次执行时间"
          align="center">
      </el-table-column>
      <el-table-column
          prop="next_execute_at"
          label="下次执行时间"
          align="center">
      </el-table-column>
      <el-table-column align="center" label="操作">
        <template slot-scope="scope">
          <el-button circle type="primary" @click="showUserEdit(scope.row)" icon="el-icon-edit"></el-button>
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

    <el-drawer title="任务添加/编辑" :visible.sync="taskEditFormShow"
               direction="rtl"
               @opened="()=>{this.$refs.mInput.focus()}"
               :before-close="handleClose"
               size="30%">
      <el-form ref="taskForm" :model="taskForm" class="drawerForm" :rules="taskFormRule" label-width="70px">
        <el-form-item label="任务名称">
          <el-input v-model.trim="taskForm.task_name" ref="mInput"></el-input>
        </el-form-item>
        <el-form-item label="任务说明">
          <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 16}" v-model="taskForm.description"></el-input>
        </el-form-item>
        <el-form-item label="所属分组" v-model="taskForm.group_id">
          <el-select v-model="taskForm.group_id">
            <el-option v-for="item in groupOptions" :key="item.group_id" :label="item.group_name" :value="item.group_id"/>
          </el-select>
        </el-form-item>
        <el-form-item label="cron表达式">
          <el-input v-model.trim="taskForm.cron_spec"></el-input>
        </el-form-item>
        <el-form-item label="命令脚本">
          <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 16}" v-model="taskForm.command"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleClose">取 消</el-button>
          <el-button type="primary" @click="saveTask">确 定</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {getTaskList, saveTask} from "../../api/task";
import {getGroupList} from "../../api/group";

export default {
  name: "UserList",
  mixins: [tableInfo],
  data() {
    return {
      taskEditFormShow: false,
      taskForm: {
        task_id: 0,
        group_id:0,
        task_name: '',
        description: '',
        cron_spec: '',
        command: ''
      },
      groupOptions: '',
      taskFormRule: {},
      uStatus: [
        {
          value: 1,
          label: '正常'
        },
        {
          value: 2,
          label: '禁用'
        }
      ],
      searchForm: {
        status: 1   //默认查状态正常的用户
      }
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    getList: getTaskList,
    showUserEdit(row) {
      console.log(row)
      for (let k in this.taskForm) {
        this.$set(this.taskForm, k, row[k] ? row[k] : '')
      }
      this.getGroups()
      this.taskEditFormShow = true
    },
    saveTask() {
      this.$refs.uForm.validate(async (valid) => {
        if (valid) {
          await saveTask(this.userForm).then((res) => {
            this.$message({
              type: 'success',
              message: res.msg,
              showClose: true
            });
            this.getTableData()
            this.taskEditFormShow = false
          }).catch(() => {
          })
        }
      })
    },
    handleClose() {
      this.$refs.taskForm.clearValidate()
      this.taskEditFormShow = false
    },
    async getGroups() {
      await getGroupList({}).then((res) => {
        this.groupOptions = res.data.list
      }).catch(() => {
      })
    },
  }
}
</script>
