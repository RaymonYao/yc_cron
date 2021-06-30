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
          <el-option key="user_name" label="任务名" value="user_name"/>
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
          prop="user_id"
          label="ID"
          width="120">
      </el-table-column>
      <el-table-column
          prop="user_name"
          label="登录名"
          width="180">
      </el-table-column>
      <el-table-column
          prop="nick_name"
          label="用户昵称">
      </el-table-column>
      <el-table-column
          label="用户状态">
      </el-table-column>
      <el-table-column
          prop="create_at"
          label="创建时间">
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

    <el-drawer title="任务添加/编辑" :visible.sync="userEditFormShow"
               direction="rtl"
               @opened="()=>{this.$refs.mInput.focus()}"
               :before-close="handleClose"
               size="30%">
      <el-form ref="uForm" :model="taskForm" class="drawerForm" :rules="taskFormRule" label-width="70px">
        <el-form-item label="任务名称">
          <el-input v-model.trim="taskForm.task_name" ref="mInput"></el-input>
        </el-form-item>
        <el-form-item label="任务说明">
          <el-input type="textarea" v-model="taskForm.description"></el-input>
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="taskForm.group_id" placeholder="未分组">
            <el-option label="采购" value="shanghai"></el-option>
            <el-option label="销售" value="beijing"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="cron表达式">
          <el-input v-model.trim="taskForm.cron_spec"></el-input>
        </el-form-item>
        <el-form-item label="命令脚本">
          <el-input type="textarea" v-model="taskForm.command"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleClose">取 消</el-button>
          <el-button type="primary" @click="userSave">确 定</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {getTaskList} from "../../api/task";

export default {
  name: "UserList",
  mixins: [tableInfo],
  data() {
    return {
      userEditFormShow: false,
      taskForm: {
        user_id: 0,
        group_id: 0,
        task_name: '',
        description: '',
        cron_spec: '',
        command: ''
      },
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
    // this.getTableData()
  },
  methods: {
    getList: getTaskList,
    userSave: function () {

    },
    handleClose: function () {

    }
  }
}
</script>
