<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="节点状态" prop="condition">
        <el-select style="width: 120px" v-model.number="searchForm.status">
          <el-option v-for="v in lStatus" :key="v.value" :label="v.label" :value="v.value"/>
        </el-select>
      </el-form-item>
      <el-form-item label="筛选条件" prop="condition">
        <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
          <el-option key="task_name" label="节点IP" value="task_name"/>
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
      <el-table-column align="center" label="操作">

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
import {getWorkerList} from "@/api/worker";
import dateTool from "@/plugins/mixins/dateTool";

export default {
  name: "WorkerList",
  mixins: [tableInfo, dateTool],
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
    this.getTableData()
  },
  methods: {
    getList: getWorkerList
  }
}
</script>
