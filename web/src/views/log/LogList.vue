<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="执行状态" prop="condition">
        <el-select style="width: 120px" v-model.number="searchForm.status">
          <el-option v-for="v in lStatus" :key="v.value" :label="v.label" :value="v.value"/>
        </el-select>
      </el-form-item>
      <el-form-item label="筛选条件" prop="condition">
        <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
          <el-option key="task_name" label="任务名称" value="task_name"/>
          <el-option key="task_id" label="任务ID" value="task_id"/>
        </el-select>
      </el-form-item>
      <el-form-item prop="searchValue">
        <el-input v-model="searchForm.sValue" placeholder="搜索关键词"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="getTableData(1)" type="primary">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border style="width: 100%">
      <el-table-column
          prop="log_id"
          label="ID"
          width="120"
          align="center">
      </el-table-column>
      <el-table-column
          prop="task_id"
          label="任务ID"
          width="120"
          align="center">
      </el-table-column>
      <el-table-column
          prop="Task.task_name"
          label="任务名称"
          width="300"
          align="center">
      </el-table-column>
      <el-table-column
          prop="start_time"
          label="开始时间"
          align="center"
          width="200"
          :formatter="dateFormat">
      </el-table-column>
      <el-table-column
          prop="end_time"
          label="结束时间"
          align="center"
          width="200"
          :formatter="dateFormat">
      </el-table-column>
      <el-table-column
          prop="process_time"
          label="消耗时间(s)"
          align="center"
          width="100">
        <template slot-scope="scope">
          {{ scope.row.end_time - scope.row.start_time }}s
        </template>
      </el-table-column>
      <el-table-column
          prop="status"
          label="执行状态"
          align="center"
          width="150">
        <template slot-scope="scope">
          <el-button v-if="scope.row.status === 0" type="success" icon="el-icon-check" circle></el-button>
          <el-button v-if="scope.row.status === -1" type="danger" icon="el-icon-close" circle></el-button>
        </template>
      </el-table-column>
      <el-table-column
          prop="run_worker"
          label="执行节点"
          align="center"
          width="250">
      </el-table-column>
      <el-table-column align="center" label="操作">
        <template slot-scope="scope">
          <el-popover placement="bottom"
                      trigger="hover"
                      width="1000">
            <LogInfo :out_put="scope.row.out_put" :error="scope.row.error"></LogInfo>
            <el-button slot="reference">查看输出</el-button>
          </el-popover>
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
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {getLogList} from "@/api/log";
import dateTool from "@/plugins/mixins/dateTool";
import LogInfo from "./cpns/LogInfo";

export default {
  name: "LogList",
  mixins: [tableInfo, dateTool],
  components: {LogInfo},
  data() {
    return {
      searchForm: {
        status: -2   //默认请选择
      },
      lStatus: [
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
          label: '错误'
        }
      ]
    }
  },
  created() {
    if (this.$route.query.sCondition !== undefined && this.$route.query.sValue !== undefined) {
      this.searchForm.sCondition = this.$route.query.sCondition
      this.searchForm.sValue = this.$route.query.sValue
    }
    this.getTableData()
  },
  methods: {
    getList: getLogList
  }
}
</script>
