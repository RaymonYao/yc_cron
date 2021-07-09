<template>
  <div>
    <el-form ref="searchForm" :model="searchForm" :inline="true">
      <el-form-item label="筛选条件" prop="condition">
        <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
          <el-option key="group_name" label="分组名称" value="group_name"/>
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
      <el-button @click="showGroupEdit({})"><i class="el-icon-plus"></i> 新增分组</el-button>
    </div>
    <el-table :data="tableData" border style="width: 100%">
      <el-table-column
          prop="group_id"
          label="ID"
          width="120"
          align="center">
      </el-table-column>
      <el-table-column
          prop="group_name"
          label="分组名称"
          width="200"
          align="center">
      </el-table-column>
      <el-table-column
          prop="description"
          label="分组描述"
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
          <el-button title="编辑" circle type="primary" @click="showGroupEdit(scope.row)" icon="el-icon-edit"></el-button>
          <el-button title="删除" type="danger" icon="el-icon-delete" circle @click="delGroup(scope.row)"
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
    <GroupEdit ref="groupEditorFormDrawer"></GroupEdit>
  </div>
</template>

<script>
import tableInfo from '@/plugins/mixins/tableInfo'
import {delGroup, getGroupList} from "@/api/group";
import GroupEdit from "./cpns/GroupEdit";
import dateTool from "@/plugins/mixins/dateTool";

export default {
  name: "GroupList",
  mixins: [tableInfo, dateTool],
  components: {GroupEdit},
  data() {
    return {
      searchForm: {}
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    getList: getGroupList,
    showGroupEdit(row) {
      this.$refs.groupEditorFormDrawer.setEditVal(row)
    },
    delGroup(row) {
      this.$confirm('确认删除该分组?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        await delGroup({group_id: row.group_id}).then((res) => {
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
